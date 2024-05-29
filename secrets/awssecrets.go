/*
 * Copyright 2024 Hypermode, Inc.
 */

package secrets

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"time"

	"hmruntime/aws"
	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/utils"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/hypermodeAI/manifest"
)

type awsSecretsProvider struct {
	prefix string
	client *secretsmanager.Client
	cache  map[string]*types.SecretValueEntry
}

func (sp *awsSecretsProvider) initialize(ctx context.Context) {
	// Initialize the Secrets Manager service client.
	// This is safe to hold onto for the lifetime of the application.
	// See https://github.com/aws/aws-sdk-go-v2/discussions/2566
	cfg := aws.GetAwsConfig()
	sp.client = secretsmanager.NewFromConfig(cfg)

	// Set the prefix based on the namespace for the backend.
	ns, err := config.GetNamespace()
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to get the namespace.")
	}
	sp.prefix = ns + "/"

	// Populate the cache with all secrets in the namespace.
	err = sp.populateSecretsCache(ctx)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to populate the secrets cache.")
	}

	// Monitor for updates to the secrets.
	go sp.monitorForUpdates(ctx)
}

func (sp *awsSecretsProvider) getHostSecrets(ctx context.Context, host manifest.HostInfo) (map[string]string, error) {
	secrets, err := sp.getSecrets(ctx, host.Name+"/")
	if err != nil {
		return nil, err
	}

	// Migrate old auth header secret to the new location
	// TODO: Remove this when we no longer need to support the old manifest format
	oldAuthHeaderSecret, ok := sp.cache[host.Name]
	if ok {
		if manifestdata.Manifest.Version == 1 {
			secrets[manifest.V1AuthHeaderVariableName] = *oldAuthHeaderSecret.SecretString
			logger.Warn(ctx).Msg("Used deprecated auth header secret.  Please update the manifest to use a template such as {{SECRET_NAME}} and migrate the old secret in Secrets Manager.")
		} else {
			logger.Warn(ctx).Msg("The manifest is current, but the deprecated auth header secret was found.  Please remove the old secret in Secrets Manager.")
		}
	}

	return secrets, nil
}

func (sp *awsSecretsProvider) getSecrets(ctx context.Context, prefix string) (map[string]string, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	results := make(map[string]string)
	for key, secret := range sp.cache {
		if strings.HasPrefix(key, prefix) {
			name := strings.TrimPrefix(key, prefix)
			results[name] = *secret.SecretString
		}
	}

	return results, nil
}

func (sp *awsSecretsProvider) getSecretValue(ctx context.Context, name string) (string, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	val, ok := sp.cache[name]
	if !ok {
		return "", fmt.Errorf("secret %s not found", name)
	}

	return *val.SecretString, nil
}

func (sp *awsSecretsProvider) populateSecretsCache(ctx context.Context) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	p := secretsmanager.NewBatchGetSecretValuePaginator(sp.client, &secretsmanager.BatchGetSecretValueInput{
		Filters: []types.Filter{{
			Key:    types.FilterNameStringTypeName,
			Values: []string{sp.prefix},
		}},
	})

	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return err
		}

		if len(page.SecretValues) == 0 {
			continue
		}

		if len(sp.cache) == 0 {
			sp.cache = make(map[string]*types.SecretValueEntry, len(page.SecretValues))
		}

		for _, secret := range page.SecretValues {
			key := strings.TrimPrefix(*secret.Name, sp.prefix)
			sp.cache[key] = &secret
		}
	}

	if len(sp.cache) == 0 {
		sp.cache = make(map[string]*types.SecretValueEntry)
	}

	return nil
}

func (sp *awsSecretsProvider) monitorForUpdates(ctx context.Context) {
	ticker := time.NewTicker(config.RefreshInterval)
	defer ticker.Stop()

	for {
		transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
		defer transaction.Finish()

		secrets, err := sp.listSecrets(ctx)
		if err != nil {
			logger.Err(ctx, err).Msg("Failed to list secrets.")
		} else {
			// Update the cache with the new secret values.
			// We don't need to batch these requests, because it's unlikely that we'll have more than a few modified secrets at a time.
			remainder := maps.Clone(sp.cache)
			for _, secret := range secrets {

				key := strings.TrimPrefix(*secret.Name, sp.prefix)
				if cachedSecret, ok := sp.cache[key]; ok {
					delete(remainder, key)
					if sp.getCurrentVersionId(secret) == *cachedSecret.VersionId {
						continue
					}
				}

				secretValue, err := sp.client.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
					SecretId: secret.ARN,
				})
				if err != nil {
					logger.Err(ctx, err).Msgf("Failed to get secret value for %s.", *secret.Name)
					continue
				}

				sp.cache[key] = &types.SecretValueEntry{
					ARN:           secretValue.ARN,
					CreatedDate:   secretValue.CreatedDate,
					Name:          secretValue.Name,
					SecretBinary:  secretValue.SecretBinary,
					SecretString:  secretValue.SecretString,
					VersionId:     secretValue.VersionId,
					VersionStages: secretValue.VersionStages,
				}
			}

			// Remove secrets that were deleted
			for key := range remainder {
				delete(sp.cache, key)
			}
		}

		// Wait for next cycle
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (sp *awsSecretsProvider) getCurrentVersionId(s types.SecretListEntry) string {
	for k, v := range s.SecretVersionsToStages {
		if v[0] == "AWSCURRENT" {
			return k
		}
	}
	return ""
}

func (sp *awsSecretsProvider) listSecrets(ctx context.Context) ([]types.SecretListEntry, error) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	p := secretsmanager.NewListSecretsPaginator(sp.client, &secretsmanager.ListSecretsInput{
		Filters: []types.Filter{{
			Key:    types.FilterNameStringTypeName,
			Values: []string{sp.prefix},
		}},
	})

	var results []types.SecretListEntry
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		if len(results) == 0 {
			results = make([]types.SecretListEntry, 0, len(page.SecretList))
		}

		results = append(results, page.SecretList...)
	}

	return results, nil
}

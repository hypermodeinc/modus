/*
 * Copyright 2024 Hypermode, Inc.
 */

package secrets

import (
	"context"
	"fmt"
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/aws"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
	"github.com/hypermodeAI/manifest"
)

type awsSecretsProvider struct {
	prefix string
	client *secretsmanager.Client
	cache  map[string]*types.SecretValueEntry
	mu     sync.RWMutex
}

func (sp *awsSecretsProvider) initialize(ctx context.Context) {
	// Initialize the Secrets Manager service client.
	// This is safe to hold onto for the lifetime of the application.
	// See https://github.com/aws/aws-sdk-go-v2/discussions/2566
	cfg := aws.GetAwsConfig()
	sp.client = secretsmanager.NewFromConfig(cfg)

	// Set the prefix based on the namespace for the backend.
	sp.prefix = config.GetNamespace() + "/"

	// Populate the cache with all secrets in the namespace.
	err := sp.populateSecretsCache(ctx)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to populate the secrets cache.")
	}

	// Monitor for updates to the secrets.
	go sp.monitorForUpdates(ctx)
}

func (sp *awsSecretsProvider) getHostSecrets(host manifest.HostInfo) (map[string]string, error) {
	hostName := host.HostName()
	secrets, err := sp.getSecrets(hostName + "/")
	if err != nil {
		return nil, err
	}

	return secrets, nil
}

func (sp *awsSecretsProvider) getSecrets(prefix string) (map[string]string, error) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	results := make(map[string]string)
	for key, secret := range sp.cache {
		if strings.HasPrefix(key, prefix) {
			name := strings.TrimPrefix(key, prefix)
			results[name] = *secret.SecretString
		}
	}

	return results, nil
}

func (sp *awsSecretsProvider) getSecretValue(name string) (string, error) {
	sp.mu.RLock()
	defer sp.mu.RUnlock()

	val, ok := sp.cache[name]
	if !ok {
		return "", fmt.Errorf("secret %s not found", name)
	}

	return *val.SecretString, nil
}

func (sp *awsSecretsProvider) populateSecretsCache(ctx context.Context) error {
	sp.mu.Lock()
	defer sp.mu.Unlock()

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
			logger.Info(ctx).Str("key", key).Msg("Secret loaded.")
		}
	}

	if len(sp.cache) == 0 {
		sp.cache = make(map[string]*types.SecretValueEntry)
		logger.Info(ctx).Msg("No secrets loaded.")
	}

	return nil
}

func (sp *awsSecretsProvider) monitorForUpdates(ctx context.Context) {
	ticker := time.NewTicker(config.RefreshInterval)
	defer ticker.Stop()

	for {
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

				sp.mu.Lock()
				sp.cache[key] = &types.SecretValueEntry{
					ARN:           secretValue.ARN,
					CreatedDate:   secretValue.CreatedDate,
					Name:          secretValue.Name,
					SecretBinary:  secretValue.SecretBinary,
					SecretString:  secretValue.SecretString,
					VersionId:     secretValue.VersionId,
					VersionStages: secretValue.VersionStages,
				}
				sp.mu.Unlock()

				logger.Info(ctx).Str("key", key).Msg("Secret updated.")
			}

			// Remove secrets that were deleted
			if len(remainder) > 0 {
				sp.mu.Lock()
				for key := range remainder {
					delete(sp.cache, key)
					logger.Info(ctx).Str("key", key).Msg("Secret removed.")
				}
				sp.mu.Unlock()
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

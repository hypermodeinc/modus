/*
 * Copyright 2024 Hypermode, Inc.
 */

package secrets

import (
	"context"
	"fmt"
	"path"
	"strings"

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
	smClient *secretsmanager.Client
}

func (sp *awsSecretsProvider) initialize(ctx context.Context) {
	// Initialize the Secrets Manager service client.
	// This is safe to hold onto for the lifetime of the application.
	// See https://github.com/aws/aws-sdk-go-v2/discussions/2566
	cfg := aws.GetAwsConfig()
	sp.smClient = secretsmanager.NewFromConfig(cfg)
}

func (sp *awsSecretsProvider) getHostSecrets(ctx context.Context, host manifest.HostInfo) (map[string]string, error) {
	secrets, err := sp.getSecrets(ctx, host.Name)
	if err != nil {
		return nil, err
	}

	// Migrate old auth header secret to the new location
	// TODO: Remove this when we no longer need to support the old manifest format
	oldAuthHeaderSecret, ok := secrets[""]
	if ok {
		if manifestdata.Manifest.Version == 1 {
			secrets[manifest.V1AuthHeaderVariableName] = oldAuthHeaderSecret
			delete(secrets, "")
			logger.Warn(ctx).Msg("Used deprecated auth header secret.  Please update the manifest to use a template such as {{SECRET_NAME}} and migrate the old secret in Secrets Manager.")
		} else {
			logger.Warn(ctx).Msg("The manifest is current, but the deprecated auth header secret was found.  Please remove the old secret in Secrets Manager.")
		}
	}

	return secrets, nil
}

func (sp *awsSecretsProvider) getSecrets(ctx context.Context, prefix string) (map[string]string, error) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// All secrets are prefixed with the namespace for the backend.
	ns, err := config.GetNamespace()
	if err != nil {
		return nil, err
	}
	prefix = path.Join(ns, prefix)

	// TODO: use the secrets cache the secrets, but we can't just look them up in the cache by prefix
	// We'll need to update the cache automatically when secrets are modified.

	// TODO: handle pagination

	out, err := sp.smClient.BatchGetSecretValue(ctx, &secretsmanager.BatchGetSecretValueInput{
		Filters: []types.Filter{{
			Key:    types.FilterNameStringTypeName,
			Values: []string{prefix},
		}},
	})

	if err != nil {
		return nil, err
	}

	secrets := make(map[string]string)
	for _, secret := range out.SecretValues {
		name := strings.Trim(strings.TrimPrefix(*secret.Name, prefix), "/")
		secrets[name] = *secret.SecretString
	}

	// TODO: Cache the secret for future use
	// secretsCache[secretId] = *secretValue.SecretString

	return secrets, nil
}

func (sp *awsSecretsProvider) getSecretValue(ctx context.Context, name string) (string, error) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// // Return the secret from the cache if it exists
	// if secret, ok := secretsCache[secretId]; ok {
	// 	return secret, nil
	// }

	// All secrets are prefixed with the namespace for the backend.
	ns, err := config.GetNamespace()
	if err != nil {
		return "", err
	}
	name = path.Join(ns, name)

	secretValue, err := sp.smClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &name,
	})

	if err != nil {
		return "", fmt.Errorf("error getting secret: %w", err)
	}
	if secretValue.SecretString == nil {
		return "", fmt.Errorf("secret string was empty")
	}

	// // Cache the secret for future use
	// secretsCache[secretId] = *secretValue.SecretString

	return *secretValue.SecretString, nil
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"
	"strings"

	"hmruntime/utils"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager/types"
)

var smClient *secretsmanager.Client

// TODO: Prefetch secrets on startup and refresh them periodically

var secretsCache = make(map[string]string)

func initSecretsManager() {
	// Initialize the Secrets Manager service client.
	// This is safe to hold onto for the lifetime of the application.
	// See https://github.com/aws/aws-sdk-go-v2/discussions/2566
	smClient = secretsmanager.NewFromConfig(awsConfig)
}

func GetSecrets(ctx context.Context, prefix string) (map[string]string, error) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// TODO: use the secrets cache the secrets, but we can't just look them up in the cache by prefix
	// We'll need to update the cache automatically when secrets are modified.

	out, err := smClient.BatchGetSecretValue(ctx, &secretsmanager.BatchGetSecretValueInput{
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

func GetSecret(ctx context.Context, secretId string) (string, error) {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	// Return the secret from the cache if it exists
	if secret, ok := secretsCache[secretId]; ok {
		return secret, nil
	}

	secretValue, err := smClient.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
		SecretId: &secretId,
	})

	if err != nil {
		return "", fmt.Errorf("error getting secret: %w", err)
	}
	if secretValue.SecretString == nil {
		return "", fmt.Errorf("secret string was empty")
	}

	// Cache the secret for future use
	secretsCache[secretId] = *secretValue.SecretString

	return *secretValue.SecretString, nil
}

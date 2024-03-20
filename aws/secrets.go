/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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

func GetSecretString(ctx context.Context, secretId string) (string, error) {

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

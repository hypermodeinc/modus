/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// TODO: Prefetch secrets on startup and refresh them periodically

var secretsCache = make(map[string]string)

func GetSecretString(ctx context.Context, secretId string) (string, error) {

	// Return the secret from the cache if it exists
	if secret, ok := secretsCache[secretId]; ok {
		return secret, nil
	}

	svc := secretsmanager.NewFromConfig(awsConfig)
	secretValue, err := svc.GetSecretValue(ctx, &secretsmanager.GetSecretValueInput{
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

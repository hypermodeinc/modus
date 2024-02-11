/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretString(ctx context.Context, secretId string) (string, error) {

	if !awsEnabled {
		return "", fmt.Errorf("unable to retrieve secret because AWS functionality is disabled")
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

	return *secretValue.SecretString, nil
}

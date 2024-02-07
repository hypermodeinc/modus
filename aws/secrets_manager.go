/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecretString(secretId string) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		return "", fmt.Errorf("error getting Secrets Manager session: %w", err)
	}

	svc := secretsmanager.New(sess)
	secretValue, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
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

/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

func GetSecretManagerSession() (*secretsmanager.SecretsManager, error) {
	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	svc := secretsmanager.New(sess)

	return svc, nil
}

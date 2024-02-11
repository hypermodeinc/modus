/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var awsConfig aws.Config
var awsEnabled bool

func Initialize(ctx context.Context) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("error loading AWS configuration: %w", err)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("error getting AWS caller identity: %w", err)
	}

	awsConfig = cfg
	awsEnabled = true

	fmt.Printf("AWS configuration loaded:\n")
	fmt.Printf("  Region:  %s\n", awsConfig.Region)
	fmt.Printf("  Account: %s\n", *identity.Account)
	fmt.Printf("  User ID: %s\n", *identity.UserId)
	fmt.Printf("  ARN: %s\n", *identity.Arn)

	return nil
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"
	"log"

	hmConfig "hmruntime/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var awsConfig aws.Config
var awsEnabled bool
var useS3PluginStorage bool

func UseAwsForPluginStorage() bool {
	return useS3PluginStorage
}

func Initialize(ctx context.Context) error {

	useS3PluginStorage = hmConfig.S3Bucket != ""
	defer func() {
		if !useS3PluginStorage {
			fmt.Println("S3 bucket name is not set.  Using local storage for plugins.")
		} else if !awsEnabled {
			log.Fatalln("S3 bucket name is set, but AWS configuration failed to load.  Exiting.")
		} else {
			fmt.Printf("Using S3 bucket %s for plugin storage.\n", hmConfig.S3Bucket)
		}
	}()

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

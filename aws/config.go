/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"

	hmConfig "hmruntime/config"
	"hmruntime/logger"

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
			logger.Info(ctx).Msg("S3 bucket name is not set.  Using local storage for plugins.")
		} else if !awsEnabled {
			logger.Fatal(ctx).Msg("S3 bucket name is set, but AWS configuration failed to load.  Exiting.")
		} else {
			logger.Info(ctx).
				Str("bucket", hmConfig.S3Bucket).
				Msg("Using S3 for plugin storage.")
		}
	}()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Warn(ctx).Err(err).
			Msg("Error loading AWS configuration.")
		return nil
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		logger.Warn(ctx).Err(err).
			Msg("Error getting AWS caller identity.")
		return nil
	}

	awsConfig = cfg
	awsEnabled = true

	logger.Info(ctx).
		Str("region", awsConfig.Region).
		Str("account", *identity.Account).
		Str("userid", *identity.UserId).
		Msg("AWS configuration loaded.")

	return nil
}

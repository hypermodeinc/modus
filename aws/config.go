/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"

	"hmruntime/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var awsConfig aws.Config

func GetAwsConfig() aws.Config {
	return awsConfig
}

func Initialize(ctx context.Context) error {
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

	logger.Info(ctx).
		Str("region", awsConfig.Region).
		Str("account", *identity.Account).
		Str("userid", *identity.UserId).
		Msg("AWS configuration loaded.")

	return nil
}

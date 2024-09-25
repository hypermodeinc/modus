/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	hmConfig "hypruntime/config"
	"hypruntime/logger"
	"hypruntime/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

var awsConfig aws.Config

func GetAwsConfig() aws.Config {
	return awsConfig
}

func Initialize(ctx context.Context) {
	if !(hmConfig.UseAwsStorage || hmConfig.UseAwsSecrets) {
		return
	}

	err := initialize(ctx)
	if err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to initialize AWS.  Exiting.")
	}
}

func initialize(ctx context.Context) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryer(func() aws.Retryer {
		return retry.NewStandard(func(o *retry.StandardOptions) {
			// double the default values to allow for longer retries
			// see https://github.com/aws/aws-sdk-go-v2/discussions/2561
			o.MaxAttempts = retry.DefaultMaxAttempts * 2 // increase from 3 to 6 attempts
			o.MaxBackoff = retry.DefaultMaxBackoff * 2   // increase from 20s to 40s
		})
	}))

	if err != nil {
		return fmt.Errorf("error loading AWS configuration: %w", err)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("error getting AWS caller identity: %w", err)
	}

	awsConfig = cfg

	logger.Info(ctx).
		Str("region", awsConfig.Region).
		Str("account", *identity.Account).
		Str("userid", *identity.UserId).
		Msg("AWS configuration loaded.")

	return nil
}

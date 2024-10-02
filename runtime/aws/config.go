/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"

	hmConfig "github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

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

	logger.Info(ctx).
		Str("region", awsConfig.Region).
		Str("account", *identity.Account).
		Str("userid", *identity.UserId).
		Msg("AWS configuration loaded.")

	return nil
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package storage

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type awsStorageProvider struct {
	s3Client *s3.Client
	s3Bucket string
	s3Path   string
}

func (stg *awsStorageProvider) initialize(ctx context.Context) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	appConfig := app.Config()
	stg.s3Bucket = appConfig.S3Bucket()
	stg.s3Path = appConfig.S3Path()

	if stg.s3Bucket == "" {
		const msg = "An S3 bucket is required when using AWS storage.  Exiting."
		sentryutils.CaptureError(ctx, nil, msg)
		logger.Fatal(ctx).Msg(msg)
	}

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		const msg = "Failed to load AWS configuration.  Exiting."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
	}

	client := sts.NewFromConfig(cfg)
	identity, err := client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		const msg = "Failed to get AWS caller identity.  Exiting."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Fatal(ctx, err).Msg(msg)
	}

	stg.s3Client = s3.NewFromConfig(cfg)

	logger.Info(ctx).
		Str("region", cfg.Region).
		Str("account", *identity.Account).
		Str("userid", *identity.UserId).
		Msg("AWS configuration loaded.")
}

func (stg *awsStorageProvider) listFiles(ctx context.Context, patterns ...string) ([]FileInfo, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: &stg.s3Bucket,
		Prefix: &stg.s3Path,
	}

	result, err := stg.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in S3 bucket: %w", err)
	}

	var files = make([]FileInfo, 0, *result.KeyCount)
	for _, obj := range result.Contents {

		_, filename := path.Split(*obj.Key)

		matched := false
		for _, pattern := range patterns {
			if match, err := path.Match(pattern, filename); err == nil && match {
				matched = true
				break
			}
		}
		if !matched {
			continue
		}

		files = append(files, FileInfo{
			Name:         filename,
			Hash:         *obj.ETag,
			LastModified: *obj.LastModified,
		})
	}

	return files, nil
}

func (stg *awsStorageProvider) getFileContents(ctx context.Context, name string) ([]byte, error) {
	key := path.Join(stg.s3Path, name)
	input := &s3.GetObjectInput{
		Bucket: &stg.s3Bucket,
		Key:    &key,
	}

	obj, err := stg.s3Client.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get file %s from S3: %w", name, err)
	}

	defer obj.Body.Close()
	content, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read contents of file %s from S3: %w", name, err)
	}

	return content, nil
}

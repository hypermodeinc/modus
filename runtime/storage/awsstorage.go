/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package storage

import (
	"context"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/hypermodeinc/modus/runtime/aws"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type awsStorageProvider struct {
	s3Client *s3.Client
}

func (stg *awsStorageProvider) initialize(ctx context.Context) {
	if config.S3Bucket == "" {
		logger.Fatal(ctx).Msg("An S3 bucket is required when using AWS storage.  Exiting.")
	}

	// Initialize the S3 service client.
	// This is safe to hold onto for the lifetime of the application.
	// See https://github.com/aws/aws-sdk-go-v2/discussions/2566
	cfg := aws.GetAwsConfig()
	stg.s3Client = s3.NewFromConfig(cfg)
}

func (stg *awsStorageProvider) listFiles(ctx context.Context, extension string) ([]FileInfo, error) {

	input := &s3.ListObjectsV2Input{
		Bucket: &config.S3Bucket,
		Prefix: &config.S3Path,
	}

	result, err := stg.s3Client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in S3 bucket: %w", err)
	}

	var files = make([]FileInfo, 0, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, extension) {
			continue
		}

		_, filename := path.Split(*obj.Key)
		files = append(files, FileInfo{
			Name:         filename,
			Hash:         *obj.ETag,
			LastModified: *obj.LastModified,
		})
	}

	return files, nil
}

func (stg *awsStorageProvider) getFileContents(ctx context.Context, name string) ([]byte, error) {
	key := path.Join(config.S3Path, name)
	input := &s3.GetObjectInput{
		Bucket: &config.S3Bucket,
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

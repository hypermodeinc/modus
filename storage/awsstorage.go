/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"fmt"
	"hmruntime/aws"
	"hmruntime/config"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
)

type awsStorage struct {
}

func (s *awsStorage) initialize() {
	if config.S3Bucket == "" {
		log.Fatal().Msg("An S3 bucket is required when using AWS storage.  Exiting.")
	}
}

func (s *awsStorage) listFiles(ctx context.Context, extension string) ([]FileInfo, error) {
	cfg := aws.GetAwsConfig()
	svc := s3.NewFromConfig(cfg)

	path := getAwsS3PathPrefix()
	input := &s3.ListObjectsV2Input{
		Bucket: &config.S3Bucket,
		Prefix: &path,
	}

	result, err := svc.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list files in S3 bucket: %w", err)
	}

	var files = make([]FileInfo, 0, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, extension) {
			continue
		}

		files = append(files, FileInfo{
			Name:         strings.TrimPrefix(*obj.Key, path),
			Hash:         *obj.ETag,
			LastModified: *obj.LastModified,
		})
	}

	return files, nil
}

func (s *awsStorage) getFileContents(ctx context.Context, name string) ([]byte, error) {
	cfg := aws.GetAwsConfig()
	svc := s3.NewFromConfig(cfg)

	key := getAwsS3PathPrefix() + name
	input := &s3.GetObjectInput{
		Bucket: &config.S3Bucket,
		Key:    &key,
	}

	obj, err := svc.GetObject(ctx, input)
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

func getAwsS3PathPrefix() string {
	path := strings.TrimRight(config.S3Path, "/")
	if path == "" {
		return path
	}
	return path + "/"
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"
	"io"
	"strings"

	"hmruntime/config"
	"hmruntime/logger"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3RetrievalHelper(ctx context.Context) (*s3.ListObjectsV2Output, error) {
	if !useS3PluginStorage {
		return nil, fmt.Errorf("unable too retrieve from S3, S3 plugin storage is disabled")
	}

	path := getPathPrefix()
	input := &s3.ListObjectsV2Input{
		Bucket: &config.S3Bucket,
		Prefix: &path,
	}

	svc := s3.NewFromConfig(awsConfig)
	result, err := svc.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting object from S3: %w", err)
	}
	return result, nil
}

func ListJsons(ctx context.Context) (map[string]string, error) {
	result, err := S3RetrievalHelper(ctx)
	if err != nil {
		return nil, err
	}

	var jsons = make(map[string]string, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, ".json") {
			continue
		}

		name := strings.TrimSuffix(strings.TrimPrefix(*obj.Key, getPathPrefix()), ".json")
		jsons[name] = *obj.ETag
	}

	return jsons, nil
}

func ListPluginsFiles(ctx context.Context) (map[string]string, error) {

	result, err := S3RetrievalHelper(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing objects from S3 bucket: %w", err)
	}

	var files = make(map[string]string, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, ".wasm") {
			continue
		}

		name := strings.TrimPrefix(*obj.Key, getPathPrefix())
		files[name] = *obj.ETag
	}

	return files, nil
}

func GetJsonBytes(ctx context.Context, name string) ([]byte, error) {
	if !useS3PluginStorage {
		return nil, fmt.Errorf("unable to retrieve JSON because S3 plugin storage is disabled")
	}

	key := getPathPrefix() + name + ".json"
	input := &s3.GetObjectInput{
		Bucket: &config.S3Bucket,
		Key:    &key,
	}

	svc := s3.NewFromConfig(awsConfig)
	obj, err := svc.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting object for JSON '%s' from S3: %w", name, err)
	}

	defer obj.Body.Close()
	bytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading content stream of JSON '%s' from S3: %w", name, err)
	}

	logger.Info(ctx).
		Str("key", key).
		Msg(fmt.Sprintf("Retrieved JSON '%s' from S3.", name))

	return bytes, nil
}

func GetPluginBytes(ctx context.Context, path string) ([]byte, error) {

	if !useS3PluginStorage {
		return nil, fmt.Errorf("unable to retrieve plugin because S3 plugin storage is disabled")
	}

	key := getPathPrefix() + path
	input := &s3.GetObjectInput{
		Bucket: &config.S3Bucket,
		Key:    &key,
	}

	svc := s3.NewFromConfig(awsConfig)
	obj, err := svc.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting object %s from S3: %w", path, err)
	}

	defer obj.Body.Close()
	bytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading content stream of %s from S3: %w", path, err)
	}

	return bytes, nil
}

func getPathPrefix() string {
	path := strings.TrimRight(config.PluginsPath, "/") + "/"
	if path == "/" {
		path = ""
	}
	return path
}

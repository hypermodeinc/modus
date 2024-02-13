/*
 * Copyright 2024 Hypermode, Inc.
 */

package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"hmruntime/config"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func UseAwsForPluginStorage() bool {
	return awsEnabled && config.S3Bucket != ""
}

func ListPlugins(ctx context.Context) ([]string, error) {

	// If we're not using AWS, or the S3 bucket is not set, return nil instead of error.
	// The caller will then use the local filesystem to load plugins instead.
	if !UseAwsForPluginStorage() {
		return nil, nil
	}

	path := config.PluginsPath
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// TODO: If we ever expect a single backend to have more than 1000 plugins, we'll need to handle pagination.

	svc := s3.NewFromConfig(awsConfig)
	input := &s3.ListObjectsV2Input{
		Bucket: &config.S3Bucket,
		Prefix: &path,
	}

	result, err := svc.ListObjectsV2(ctx, input)
	if err != nil {
		log.Fatalln(fmt.Errorf("error listing plugins: %w", err))
	}

	var plugins = make([]string, 0, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, ".wasm") {
			continue
		}
		plugin := strings.TrimSuffix(strings.TrimPrefix(*obj.Key, path), ".wasm")
		plugins = append(plugins, plugin)
	}

	return plugins, nil
}

func GetPluginBytes(ctx context.Context, name string) ([]byte, error) {
	key := config.PluginsPath + "/" + name + ".wasm"
	input := &s3.GetObjectInput{
		Bucket: &config.S3Bucket,
		Key:    &key,
	}

	svc := s3.NewFromConfig(awsConfig)
	obj, err := svc.GetObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error getting object for plugin '%s' from S3: %w", name, err)
	}

	defer obj.Body.Close()
	bytes, err := io.ReadAll(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading content stream of plugin '%s' from S3: %w", name, err)
	}

	fmt.Printf("Loaded from S3 object with key %s\n", key)

	return bytes, nil
}

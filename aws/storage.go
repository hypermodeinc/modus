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

func ListPlugins(ctx context.Context) (map[string]string, error) {

	if !useS3PluginStorage {
		return nil, fmt.Errorf("unable to list plugins because S3 plugin storage is disabled")
	}

	path := getPathPrefix()
	input := &s3.ListObjectsV2Input{
		Bucket: &config.S3Bucket,
		Prefix: &path,
	}

	// TODO: If we ever expect a single backend to have more than 1000 plugins, we'll need to handle pagination.

	svc := s3.NewFromConfig(awsConfig)
	result, err := svc.ListObjectsV2(ctx, input)
	if err != nil {
		log.Fatalln(fmt.Errorf("error listing plugins from S3: %w", err))
	}

	var plugins = make(map[string]string, *result.KeyCount)
	for _, obj := range result.Contents {
		if !strings.HasSuffix(*obj.Key, ".wasm") {
			continue
		}

		name := strings.TrimSuffix(strings.TrimPrefix(*obj.Key, path), ".wasm")
		plugins[name] = *obj.ETag
	}

	return plugins, nil
}

func GetPluginBytes(ctx context.Context, name string) ([]byte, error) {

	if !useS3PluginStorage {
		return nil, fmt.Errorf("unable to retrieve plugin because S3 plugin storage is disabled")
	}

	key := getPathPrefix() + name + ".wasm"
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

func getPathPrefix() string {
	path := strings.TrimRight(config.PluginsPath, "/") + "/"
	if path == "/" {
		path = ""
	}
	return path
}

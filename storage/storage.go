/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"hmruntime/config"
	"time"

	"hmruntime/utils"
)

var provider storageProvider

type storageProvider interface {
	initialize(ctx context.Context)
	listFiles(ctx context.Context, extension string) ([]FileInfo, error)
	getFileContents(ctx context.Context, name string) ([]byte, error)
	writeFile(ctx context.Context, name string, contents []byte) error
}

type FileInfo struct {
	Name         string
	Hash         string
	LastModified time.Time
}

func Initialize(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	if config.UseAwsStorage {
		provider = &awsStorageProvider{}
	} else {
		provider = &localStorageProvider{}
	}

	provider.initialize(ctx)
}

func ListFiles(ctx context.Context, extension string) ([]FileInfo, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return provider.listFiles(ctx, extension)
}

func GetFileContents(ctx context.Context, name string) ([]byte, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return provider.getFileContents(ctx, name)
}

func WriteFile(ctx context.Context, name string, contents []byte) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return provider.writeFile(ctx, name, contents)
}

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

var impl storageImplementation

type storageImplementation interface {
	initialize()
	listFiles(ctx context.Context, extension string) ([]FileInfo, error)
	getFileContents(ctx context.Context, name string) ([]byte, error)
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
		impl = &awsStorage{}
	} else {
		impl = &localStorage{}
	}

	impl.initialize()
}

func ListFiles(ctx context.Context, extension string) ([]FileInfo, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return impl.listFiles(ctx, extension)
}

func GetFileContents(ctx context.Context, name string) ([]byte, error) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return impl.getFileContents(ctx, name)
}

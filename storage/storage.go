/*
 * Copyright 2024 Hypermode, Inc.
 */

package storage

import (
	"context"
	"hmruntime/config"
	"time"
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

func Initialize() {
	if config.UseAwsStorage {
		impl = &awsStorage{}
	} else {
		impl = &localStorage{}
	}

	impl.initialize()
}

func ListFiles(ctx context.Context, extension string) ([]FileInfo, error) {
	return impl.listFiles(ctx, extension)
}

func GetFileContents(ctx context.Context, name string) ([]byte, error) {
	return impl.getFileContents(ctx, name)
}

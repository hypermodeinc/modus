/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package storage

import (
	"context"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
)

var provider storageProvider

type storageProvider interface {
	initialize(ctx context.Context)
	listFiles(ctx context.Context, patterns ...string) ([]FileInfo, error)
	getFileContents(ctx context.Context, name string) ([]byte, error)
}

type FileInfo struct {
	Name         string
	Hash         string
	LastModified time.Time
}

func Initialize(ctx context.Context) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	if app.Config().UseAwsStorage() {
		provider = &awsStorageProvider{}
	} else {
		provider = &localStorageProvider{}
	}

	provider.initialize(ctx)
}

func ListFiles(ctx context.Context, patterns ...string) ([]FileInfo, error) {
	return provider.listFiles(ctx, patterns...)
}

func GetFileContents(ctx context.Context, name string) ([]byte, error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	return provider.getFileContents(ctx, name)
}

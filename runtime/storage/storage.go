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
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/utils"
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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
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
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	return provider.getFileContents(ctx, name)
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import (
	"context"
	"fmt"

	"github.com/puzpuzpuz/xsync/v4"
)

var dsr = newDSRegistry()

type dsRegistry struct {
	cache *xsync.Map[string, dataSource]
}

func newDSRegistry() *dsRegistry {
	return &dsRegistry{
		cache: xsync.NewMap[string, dataSource](),
	}
}

func (r *dsRegistry) shutdown() {
	r.cache.Range(func(key string, _ dataSource) bool {
		if ds, ok := dsr.cache.LoadAndDelete(key); ok {
			ds.Shutdown()
		}
		return true
	})
}

func (r *dsRegistry) getDataSource(ctx context.Context, dsName, dsType string) (dataSource, error) {
	var creationErr error
	ds, _ := r.cache.LoadOrCompute(dsName, func() (dataSource, bool) {
		switch dsType {
		case "postgresql":
			if ds, err := newPostgresqlDS(ctx, dsName); err != nil {
				creationErr = err
				return nil, true
			} else {
				return ds, false
			}
		case "mysql":
			if ds, err := newMysqlDS(ctx, dsName); err != nil {
				creationErr = err
				return nil, true
			} else {
				return ds, false
			}
		default:
			creationErr = fmt.Errorf("unsupported data source type: %s", dsType)
			return nil, true
		}
	})
	return ds, creationErr
}

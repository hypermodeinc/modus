/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraphclient

import (
	"context"

	"github.com/hypermodeinc/modus/runtime/manifestdata"
)

func Initialize() {
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		ShutdownConns()
		return nil
	})
}

func ExecuteQuery(ctx context.Context, hostName string, req *Request) (*Response, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return nil, err
	}

	return dc.execute(ctx, req)
}

func AlterSchema(ctx context.Context, hostName, schema string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.alterSchema(ctx, schema)
}

func DropAttribute(ctx context.Context, hostName, attr string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.dropAttr(ctx, attr)
}

func DropAllData(ctx context.Context, hostName string) (string, error) {
	dc, err := dgr.getDgraphConnector(ctx, hostName)
	if err != nil {
		return "", err
	}

	return dc.dropAll(ctx)
}

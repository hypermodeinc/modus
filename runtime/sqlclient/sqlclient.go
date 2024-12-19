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

	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func Initialize() {
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		ShutdownPGPools()
		return nil
	})
}

func ExecuteQuery(ctx context.Context, connectionName, dbType, statement, paramsJson string) (*HostQueryResponse, error) {
	var params []any
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		return nil, fmt.Errorf("error deserializing database query parameters: %w", err)
	}

	dbResponse, err := doExecuteQuery(ctx, connectionName, dbType, statement, params)
	if err != nil {
		return nil, err
	}

	var resultJson []byte
	if dbResponse.Result != nil {
		var err error
		resultJson, err = utils.JsonSerialize(dbResponse.Result)
		if err != nil {
			return nil, fmt.Errorf("error serializing result: %w", err)
		}
	}

	response := &HostQueryResponse{
		Error:        dbResponse.Error,
		RowsAffected: dbResponse.RowsAffected,
	}

	if len(resultJson) > 0 {
		s := string(resultJson)
		response.ResultJson = &s
	}

	return response, nil
}

func doExecuteQuery(ctx context.Context, dsName, dsType, stmt string, params []any) (*dbResponse, error) {
	switch dsType {
	case "postgresql":
		ds, err := dsr.getPostgresDS(ctx, dsName)
		if err != nil {
			return nil, err
		}

		return ds.query(ctx, stmt, params)

	default:
		return nil, fmt.Errorf("connection [%s] has an unsupported type: %s", dsName, dsType)
	}
}

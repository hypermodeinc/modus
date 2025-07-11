/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/hypermodeinc/modus/runtime/utils"
)

func Initialize() {
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		dsr.shutdown()
		return nil
	})
}

func Shutdown() {
	dsr.shutdown()
}

func ExecuteQuery(ctx context.Context, connectionName, dbType, statement, paramsJson string) (*HostQueryResponse, error) {

	// Small Hack:
	// The `paramsJson` argument might be prefixed with flags set by the SDK.
	// Example: "flag1,flag2:[]" (where `[]` is the actual JSON string)
	//
	// Currently supported flags:
	// - "exec" - No rows are expected to be returned. Just execute the provided statement.
	//
	// Passing the flags in this manner avoids a breaking change in the host function signature.
	//
	var execOnly bool
	paramsJson, execOnly = strings.CutPrefix(paramsJson, "exec:")

	var params []any
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		return nil, fmt.Errorf("error deserializing database query parameters: %w", err)
	}

	dbResponse, err := doExecuteQuery(ctx, connectionName, dbType, statement, params, execOnly)
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
		LastInsertID: dbResponse.LastInsertID,
	}

	if len(resultJson) > 0 {
		s := string(resultJson)
		response.ResultJson = &s
	}

	return response, nil
}

func doExecuteQuery(ctx context.Context, dsName, dsType, stmt string, params []any, execOnly bool) (*dbResponse, error) {
	ds, err := dsr.getDataSource(ctx, dsName, dsType)
	if err != nil {
		return nil, err
	}

	return ds.query(ctx, stmt, params, execOnly)
}

func getConnectionString(ctx context.Context, dsName string, connType manifest.ConnectionType) (string, error) {
	man := manifestdata.GetManifest()
	if man == nil {
		return "", fmt.Errorf("manifest not loaded")
	}

	info, ok := man.Connections[dsName]
	if !ok {
		return "", fmt.Errorf("connection [%s] not found", dsName)
	}
	if info.ConnectionType() != connType {
		return "", fmt.Errorf("[%s] is not a %s connection", dsName, connType)
	}

	var connStr string
	switch t := info.(type) {
	case manifest.MysqlConnectionInfo:
		connStr = t.ConnStr
	case manifest.PostgresqlConnectionInfo:
		connStr = t.ConnStr
	}

	if connStr == "" {
		return "", fmt.Errorf("connection [%s] has empty or missing connection string", dsName)
	}

	if cs, err := secrets.ApplySecretsToString(ctx, info, connStr); err != nil {
		return "", fmt.Errorf("failed to apply secrets for connection [%s]: %w", dsName, err)
	} else {
		return cs, nil
	}
}

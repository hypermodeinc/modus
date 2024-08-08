/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

import (
	"context"
	"fmt"

	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

func ExecuteQuery(ctx context.Context, hostName, dbType, statement, paramsJson string) (*HostQueryResponse, error) {
	var params []any
	if err := utils.JsonDeserialize([]byte(paramsJson), &params); err != nil {
		return nil, fmt.Errorf("error deserializing database query parameters: %w", err)
	}

	dbResponse, err := doExecuteQuery(ctx, hostName, dbType, statement, params)
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

func doExecuteQuery(ctx context.Context, dsname, dsType, stmt string, params []any) (*dbResponse, error) {
	switch dsType {
	case manifest.HostTypePostgresql:
		ds, err := dsr.getPGPool(ctx, dsname)
		if err != nil {
			return nil, err
		}

		return ds.query(ctx, stmt, params)

	default:
		return nil, fmt.Errorf("host %s has an unsupported type: %s", dsname, dsType)
	}
}

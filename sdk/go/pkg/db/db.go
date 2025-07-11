/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

type HostQueryResponse struct {
	Error        *string
	ResultJson   *string
	RowsAffected uint32
	LastInsertId uint64
}

// Represents the result of a database query that does not return data.
type Response struct {

	// The number of rows affected by the query.
	RowsAffected uint32

	// The ID of the last row inserted by the query, if applicable.
	// Note that not all databases support this feature.
	LastInsertId uint64
}

// Represents the result of a database query that returns rows.
type QueryResponse[T any] struct {

	// The number of rows affected by the query, which typically corresponds to the number of rows returned.
	RowsAffected uint32

	// The ID of the last row inserted by the query, if applicable.
	// Note that not all databases support this feature.
	LastInsertId uint64

	// The rows returned by the query.
	Rows []T
}

// Represents the result of a database query that returns a single scalar value.
type ScalarResponse[T any] struct {

	// The number of rows affected by the query, which is typically 1 unless the query failed,
	// or had side effects that modified more than one row.
	RowsAffected uint32

	// The ID of the last row inserted by the query, if applicable.
	// Note that not all databases support this feature.
	LastInsertId uint64

	// The scalar value returned by the query.
	Value T
}

// Executes a database query that does not return rows.
func Execute(connection, dbType, statement string, params ...any) (*Response, error) {
	r, err := doQuery(connection, dbType, statement, params, true)
	if err != nil {
		return nil, err
	}

	return &Response{
		RowsAffected: r.RowsAffected,
		LastInsertId: r.LastInsertId,
	}, nil
}

// Executes a database query that returns rows.
// The structure of the rows is determined by the type parameter.
func Query[T any](connection, dbType, statement string, params ...any) (*QueryResponse[T], error) {
	r, err := doQuery(connection, dbType, statement, params, false)
	if err != nil {
		return nil, err
	}

	var rows []T
	if r.ResultJson != nil {
		if err := utils.JsonDeserialize([]byte(*r.ResultJson), &rows); err != nil {
			return nil, fmt.Errorf("could not JSON deserialize database response: %v", err)
		}
	}

	return &QueryResponse[T]{
		RowsAffected: r.RowsAffected,
		LastInsertId: r.LastInsertId,
		Rows:         rows,
	}, nil
}

// Executes a database query that returns a single scalar value.
// The type parameter determines the type of the scalar value.
func QueryScalar[T any](connection, dbType, statement string, params ...any) (*ScalarResponse[T], error) {
	r, err := Query[map[string]any](connection, dbType, statement, params...)
	if err != nil {
		return nil, err
	}

	if len(r.Rows) == 1 {
		fields := r.Rows[0]
		if len(fields) > 1 {
			return nil, fmt.Errorf("expected a single column from a scalar database query, but received %d", len(fields))
		}

		for _, value := range fields {
			result, err := utils.ConvertInterfaceTo[T](value)
			if err != nil {
				var zero T
				return nil, fmt.Errorf("could not convert database result to %T: %v", zero, err)
			}
			return &ScalarResponse[T]{
				RowsAffected: r.RowsAffected,
				LastInsertId: r.LastInsertId,
				Value:        result,
			}, nil
		}
	} else if len(r.Rows) > 1 {
		return nil, fmt.Errorf("expected a single row from a scalar database query, but received %d", len(r.Rows))
	}

	return nil, errors.New("no result returned from database query")
}

func doQuery(connection, dbType, statement string, params []any, execOnly bool) (*HostQueryResponse, error) {
	paramsJson := "[]"
	if len(params) > 0 {
		bytes, err := utils.JsonSerialize(params)
		if err != nil {
			return nil, fmt.Errorf("could not JSON serialize query parameters: %v", err)
		}
		paramsJson = string(bytes)
	}

	if execOnly {
		// This flag instructs the host function not to return rows, but to simply execute the statement.
		paramsJson = "exec:" + paramsJson
	}

	statement = strings.TrimSpace(statement)
	response := hostExecuteQuery(&connection, &dbType, &statement, &paramsJson)
	if response == nil {
		return nil, errors.New("no response received from database query")
	}

	if response.Error != nil {
		return nil, fmt.Errorf("database returned an error: %s", *response.Error)
	}

	return response, nil
}

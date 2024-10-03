/*
 * Copyright 2024 Hypermode, Inc.
 */

package db

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hypermodeAI/functions-go/pkg/utils"
)

type HostQueryResponse struct {
	Error        *string
	ResultJson   *string
	RowsAffected uint32
}

func Execute(hostName, dbType, statement string, params ...any) (uint, error) {
	_, affected, err := doQuery(hostName, dbType, statement, params...)
	return affected, err
}

func Query[T any](hostName, dbType, statement string, params ...any) ([]T, uint, error) {
	resultJson, affected, err := doQuery(hostName, dbType, statement, params...)
	if err != nil {
		return nil, affected, err
	}

	var rows []T
	if resultJson != nil {
		if err := utils.JsonDeserialize([]byte(*resultJson), &rows); err != nil {
			return nil, affected, fmt.Errorf("could not JSON deserialize database response: %v", err)
		}
	}

	return rows, affected, nil
}

func QueryScalar[T any](hostName, dbType, statement string, params ...any) (T, uint, error) {
	var zero T

	rows, affected, err := Query[map[string]any](hostName, dbType, statement, params...)
	if err != nil {
		return zero, affected, err
	}

	if len(rows) == 1 {
		fields := rows[0]
		if len(fields) > 1 {
			return zero, affected, fmt.Errorf("expected a single column from a scalar database query, but received %d", len(fields))
		}

		for _, value := range fields {
			result, err := utils.ConvertInterfaceTo[T](value)
			if err != nil {
				return zero, affected, fmt.Errorf("could not convert database result to %T: %v", zero, err)
			}
			return result, affected, nil
		}
	} else if len(rows) > 1 {
		return zero, affected, fmt.Errorf("expected a single row from a scalar database query, but received %d", len(rows))
	}

	return zero, affected, errors.New("no result returned from database query")
}

func doQuery(hostName, dbType, statement string, params ...any) (*string, uint, error) {
	paramsJson := "[]"
	if len(params) > 0 {
		bytes, err := utils.JsonSerialize(params)
		if err != nil {
			return nil, 0, fmt.Errorf("could not JSON serialize query parameters: %v", err)
		}
		paramsJson = string(bytes)
	}

	statement = strings.TrimSpace(statement)
	response := databaseQuery(&hostName, &dbType, &statement, &paramsJson)
	if response == nil {
		return nil, 0, errors.New("no response received from database query")
	}

	affected := uint(response.RowsAffected)

	if response.Error != nil {
		return nil, affected, fmt.Errorf("database returned an error: %s", *response.Error)
	}

	return response.ResultJson, affected, nil
}

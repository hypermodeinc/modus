/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import (
	"errors"
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/console"
	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

func Execute[T any](connection, statement string, variables map[string]any) (*Response[T], error) {
	bytes, err := utils.JsonSerialize(variables)
	if err != nil {
		console.Error(err.Error())
		return nil, err
	}

	varsStr := string(bytes)

	response := hostExecuteQuery(&connection, &statement, &varsStr)

	if response == nil {
		return nil, errors.New("failed to execute the GQL query")
	}

	var result Response[T]
	err = utils.JsonDeserialize([]byte(*response), &result)
	if err != nil {
		console.Error(err.Error())
		return nil, err
	}

	if len(result.Errors) > 0 {
		errBytes, err := utils.JsonSerialize(result.Errors)
		if err != nil {
			console.Error(err.Error())
			return nil, err
		}
		console.Error(fmt.Sprint("GraphQL API Errors:" + string(errBytes)))
	}
	return &result, nil
}

type Response[T any] struct {
	Errors []ErrorResult `json:"errors"`
	Data   *T            `json:"data"`
}

type ErrorResult struct {
	Message   string         `json:"message"`
	Locations []CodeLocation `json:"locations"`
	Path      []string       `json:"path"`
}

type CodeLocation struct {
	Line   uint32 `json:"line"`
	Column uint32 `json:"column"`
}

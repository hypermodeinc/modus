/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import (
	"errors"
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/console"
	"github.com/hypermodeAI/functions-go/pkg/utils"
)

func Execute[T any](hostName, statement string, variables map[string]any) (*Response[T], error) {
	bytes, err := utils.JsonSerialize(variables)
	if err != nil {
		console.Error(err.Error())
		return nil, err
	}

	varsStr := string(bytes)

	response := hostExecuteGQL(&hostName, &statement, &varsStr)

	if response == nil {
		return nil, errors.New("Failed to execute the GQL query.")
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

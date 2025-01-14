/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4jclient

import (
	"context"
	"encoding/json"

	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func Initialize() {
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		CloseDrivers(ctx)
		return nil
	})
}

func ExecuteQuery(ctx context.Context, hostName, dbName, query string, parametersJson string) (*EagerResult, error) {
	driver, err := n4j.getDriver(ctx, hostName)
	if err != nil {
		return nil, err
	}

	parameters := make(map[string]any)
	if err := json.Unmarshal([]byte(parametersJson), &parameters); err != nil {
		return nil, err
	}

	res, err := neo4j.ExecuteQuery(ctx, driver, query, parameters, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(dbName))
	if err != nil {
		return nil, err
	}

	records := make([]*Record, len(res.Records))
	for i, record := range res.Records {
		vals := make([]string, len(record.Values))
		for j, val := range record.Values {
			valBytes, err := utils.JsonSerialize(val)
			if err != nil {
				return nil, err
			}
			vals[j] = string(valBytes)
		}
		records[i] = &Record{
			Values: vals,
			Keys:   record.Keys,
		}
	}
	return &EagerResult{
		Keys:    res.Keys,
		Records: records,
	}, nil
}

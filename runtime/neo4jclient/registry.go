/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4jclient

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/puzpuzpuz/xsync/v3"
)

var n4j = newNeo4jRegistry()

type neo4jRegistry struct {
	cache *xsync.MapOf[string, neo4j.DriverWithContext]
}

func newNeo4jRegistry() *neo4jRegistry {
	return &neo4jRegistry{
		cache: xsync.NewMapOf[string, neo4j.DriverWithContext](),
	}
}

func CloseDrivers(ctx context.Context) {
	n4j.cache.Range(func(key string, _ neo4j.DriverWithContext) bool {
		if driver, ok := n4j.cache.LoadAndDelete(key); ok {
			driver.Close(ctx)
		}
		return true
	})
}

func (nr *neo4jRegistry) getDriver(ctx context.Context, n4jName string) (neo4j.DriverWithContext, error) {
	var creationErr error
	driver, _ := n4j.cache.LoadOrTryCompute(n4jName, func() (neo4j.DriverWithContext, bool) {
		driver, err := createDriver(ctx, n4jName)
		if err != nil {
			creationErr = err
			return nil, true
		}
		return driver, false
	})
	return driver, creationErr
}

func createDriver(ctx context.Context, n4jName string) (neo4j.DriverWithContext, error) {
	man := manifestdata.GetManifest()
	info, ok := man.Connections[n4jName]
	if !ok {
		return nil, fmt.Errorf("Neo4j connection [%s] not found", n4jName)
	}

	if info.ConnectionType() != manifest.ConnectionTypeNeo4j {
		return nil, fmt.Errorf("[%s] is not a Neo4j connection", n4jName)
	}

	connection := info.(manifest.Neo4jConnectionInfo)
	if err := validateNeo4jConnection(connection); err != nil {
		return nil, err
	}

	dbUri, err := secrets.ApplySecretsToString(ctx, info, connection.DbUri)
	if err != nil {
		return nil, err
	}

	username, err := secrets.ApplySecretsToString(ctx, info, connection.Username)
	if err != nil {
		return nil, err
	}

	password, err := secrets.ApplySecretsToString(ctx, info, connection.Password)
	if err != nil {
		return nil, err
	}

	driver, err := neo4j.NewDriverWithContext(
		dbUri,
		neo4j.BasicAuth(username, password, ""),
	)
	if err != nil {
		return nil, err
	}

	return driver, nil
}

func validateNeo4jConnection(connection manifest.Neo4jConnectionInfo) error {
	var emptyFields []string

	if connection.DbUri == "" {
		emptyFields = append(emptyFields, "DbUri")
	}
	if connection.Username == "" {
		emptyFields = append(emptyFields, "Username")
	}
	if connection.Password == "" {
		emptyFields = append(emptyFields, "Password")
	}

	if len(emptyFields) > 0 {
		return fmt.Errorf("[%s] has empty required fields: %v",
			connection.Name,
			emptyFields)
	}

	return nil
}

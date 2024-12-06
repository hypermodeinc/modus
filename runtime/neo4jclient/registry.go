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
	"sync"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/secrets"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var n4j = newNeo4jRegistry()

type neo4jRegistry struct {
	sync.RWMutex
	neo4jDriverCache map[string]neo4j.DriverWithContext
}

func newNeo4jRegistry() *neo4jRegistry {
	return &neo4jRegistry{
		neo4jDriverCache: make(map[string]neo4j.DriverWithContext),
	}
}

func CloseDrivers(ctx context.Context) {
	n4j.Lock()
	defer n4j.Unlock()

	removed := make([]string, 0)

	for key, driver := range n4j.neo4jDriverCache {
		driver.Close(ctx)
		removed = append(removed, key)
	}

	for _, key := range removed {
		delete(n4j.neo4jDriverCache, key)
	}
}

func (nr *neo4jRegistry) getDriver(ctx context.Context, n4jName string) (neo4j.DriverWithContext, error) {
	nr.RLock()
	ds, ok := nr.neo4jDriverCache[n4jName]
	nr.RUnlock()
	if ok {
		return ds, nil
	}
	nr.Lock()
	defer nr.Unlock()

	if driver, ok := nr.neo4jDriverCache[n4jName]; ok {
		return driver, nil
	}

	info, ok := manifestdata.GetManifest().Connections[n4jName]
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

	nr.neo4jDriverCache[n4jName] = driver

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

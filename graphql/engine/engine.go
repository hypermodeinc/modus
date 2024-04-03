/*
 * Copyright 2024 Hypermode, Inc.
 */

package engine

import (
	"sync"

	gql "github.com/wundergraph/graphql-go-tools/v2/pkg/graphql"
)

var instance *gql.ExecutionEngineV2
var mutex sync.RWMutex

// GetEngine provides thread-safe access to the current GraphQL execution engine.
func GetEngine() *gql.ExecutionEngineV2 {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func setEngine(engine *gql.ExecutionEngineV2) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = engine
}

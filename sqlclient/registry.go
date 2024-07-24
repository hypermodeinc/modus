/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

import (
	"sync"
)

var dsr = newDSRegistry()

type dsRegistry struct {
	sync.RWMutex
	pgCache map[string]*postgresqlDS
}

func newDSRegistry() *dsRegistry {
	return &dsRegistry{
		pgCache: make(map[string]*postgresqlDS),
	}
}

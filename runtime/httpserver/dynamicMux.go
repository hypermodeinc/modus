/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpserver

import (
	"net/http"
	"sync"
)

type dynamicMux struct {
	mux    *http.ServeMux
	routes map[string]http.Handler
	mu     sync.RWMutex
}

func newDynamicMux(routes map[string]http.Handler) *dynamicMux {
	return &dynamicMux{
		mux:    http.NewServeMux(),
		routes: routes,
	}
}

func (dm *dynamicMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	if handler, ok := dm.routes[r.URL.Path]; ok {
		handler.ServeHTTP(w, r)
	} else {
		http.Error(w, "Endpoint not found", http.StatusNotFound)
	}
}

func (dm *dynamicMux) ReplaceRoutes(routes map[string]http.Handler) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.routes = routes
}

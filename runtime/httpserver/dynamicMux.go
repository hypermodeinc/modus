/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpserver

import (
	"net/http"
	"strings"
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
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	if path != r.URL.Path {
		if _, ok := dm.routes[path]; ok {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
			return
		}
	}

	for route, handler := range dm.routes {
		if len(route) > 1 && strings.HasSuffix(route, "/") &&
			(strings.HasPrefix(r.URL.Path, route) || route == r.URL.Path+"/") {
			handler.ServeHTTP(w, r)
			return
		}
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func (dm *dynamicMux) ReplaceRoutes(routes map[string]http.Handler) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.routes = routes
}

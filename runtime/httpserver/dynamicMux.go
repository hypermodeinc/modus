/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpserver

import (
	"maps"
	"net/http"
	"strings"

	"github.com/puzpuzpuz/xsync/v4"
)

type dynamicMux struct {
	mux    *http.ServeMux
	routes *xsync.Map[string, http.Handler]
}

func newDynamicMux(routes map[string]http.Handler) *dynamicMux {

	mapRoutes := xsync.NewMap[string, http.Handler](xsync.WithPresize(len(routes)))
	for path, handler := range routes {
		mapRoutes.Store(path, handler)
	}

	return &dynamicMux{
		mux:    http.NewServeMux(),
		routes: mapRoutes,
	}
}

func (dm *dynamicMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if handler, ok := dm.routes.Load(r.URL.Path); ok {
		handler.ServeHTTP(w, r)
		return
	}

	path := strings.TrimSuffix(r.URL.Path, "/")
	if path != r.URL.Path {
		if _, ok := dm.routes.Load(path); ok {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
			return
		}
	}

	var found bool
	dm.routes.Range(func(route string, handler http.Handler) bool {
		if len(route) > 1 && strings.HasSuffix(route, "/") &&
			(strings.HasPrefix(r.URL.Path, route) || route == r.URL.Path+"/") {
			handler.ServeHTTP(w, r)
			found = true
			return false
		} else {
			return true
		}
	})

	if !found {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}

func (dm *dynamicMux) ReplaceRoutes(routes map[string]http.Handler) {
	temp := maps.Clone(routes)
	dm.routes.Range(func(route string, _ http.Handler) bool {
		if handler, ok := temp[route]; ok {
			dm.routes.Store(route, handler)
			delete(temp, route)
		} else {
			dm.routes.Delete(route)
		}
		return true
	})
	for route, handler := range temp {
		dm.routes.Store(route, handler)
	}
}

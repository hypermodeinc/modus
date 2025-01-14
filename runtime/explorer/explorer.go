/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package explorer

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/utils"
)

//go:embed content
var content embed.FS
var contentRoot, _ = fs.Sub(content, "content/dist")

var ExplorerHandler = http.HandlerFunc(explorerHandler)

func explorerHandler(w http.ResponseWriter, r *http.Request) {

	mux := http.NewServeMux()
	mux.Handle("/explorer/", http.StripPrefix("/explorer/", http.FileServerFS(contentRoot)))
	mux.HandleFunc("/explorer/api/endpoints", endpointsHandler)

	mux.ServeHTTP(w, r)
}

func endpointsHandler(w http.ResponseWriter, r *http.Request) {

	type endpoint struct {
		ApiType string `json:"type"`
		Name    string `json:"name"`
		Path    string `json:"path"`
	}

	m := manifestdata.GetManifest()

	endpoints := make([]endpoint, 0, len(m.Endpoints))
	for name, ep := range m.Endpoints {
		switch ep.EndpointType() {
		case manifest.EndpointTypeGraphQL:
			info := ep.(manifest.GraphqlEndpointInfo)
			endpoints = append(endpoints, endpoint{"GraphQL", name, info.Path})
		}
	}

	utils.WriteJsonContentHeader(w)
	j, _ := utils.JsonSerialize(endpoints)
	_, _ = w.Write(j)
}

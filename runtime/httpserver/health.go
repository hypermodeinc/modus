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
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hypermodeinc/modus/runtime/actors"
	"github.com/hypermodeinc/modus/runtime/app"
)

var healthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	env := app.Config().Environment()
	ver := app.VersionNumber()
	agents := actors.ListLocalAgents(r.Context())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// custom format the JSON response for easy readability
	_, _ = w.Write([]byte(`{
  "status": "ok",
  "environment": "` + env + `",
  "version": "` + ver + `",
`))

	if len(agents) == 0 {
		_, _ = w.Write([]byte(`  "agents": []` + "\n"))
	} else {
		_, _ = w.Write([]byte(`  "agents": [` + "\n"))
		for i, agent := range agents {
			if i > 0 {
				_, _ = w.Write([]byte(",\n"))
			}
			name, _ := json.Marshal(agent.Name)
			_, _ = w.Write(fmt.Appendf(nil, `    {"id": "%s", "name": %s, "status": "%s"}`, agent.Id, name, agent.Status))
		}
		_, _ = w.Write([]byte("\n  ]\n"))
	}

	_, _ = w.Write([]byte("}\n"))

})

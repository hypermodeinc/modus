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

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/utils"
)

var healthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	env := app.Config().Environment()
	ver := app.VersionNumber()
	w.WriteHeader(http.StatusOK)
	utils.WriteJsonContentHeader(w)
	_, _ = w.Write([]byte(`{"status":"ok","environment":"` + env + `","version":"` + ver + `"}`))
})

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpserver

import (
	"net/http"
	"runtime"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/utils"
)

const msg = "Failed to serialize health check response."

var healthHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	data := []utils.KeyValuePair{
		{Key: "status", Value: "ok"},
		{Key: "environment", Value: app.Config().Environment()},
		{Key: "app_version", Value: app.VersionNumber()},
		{Key: "go_version", Value: runtime.Version()},
	}
	if ns, ok := app.KubernetesNamespace(); ok {
		data = append(data, utils.KeyValuePair{Key: "kubernetes_namespace", Value: ns})
	}

	jsonBytes, err := utils.MakeJsonObject(data, true)
	if err != nil {
		var ctx = r.Context()
		logger.Error(ctx, err).Msg(msg)
		sentryutils.CaptureError(ctx, err, msg)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonBytes)
})

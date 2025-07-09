/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/graphql/engine"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/sentryutils"
	"github.com/hypermodeinc/modus/runtime/timezones"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
	"github.com/puzpuzpuz/xsync/v4"
	"github.com/rs/xid"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	eng "github.com/wundergraph/graphql-go-tools/execution/engine"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/graphqlerrors"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/operationreport"
)

var GraphQLRequestHandler = http.HandlerFunc(handleGraphQLRequest)

var cancelFuncs = xsync.NewMap[string, context.CancelFunc]()

func CancelSubscriptions() {
	cancelFuncs.Range(func(key string, cancel context.CancelFunc) bool {
		cancel()
		return true
	})
	cancelFuncs.Clear()
}

func Initialize() {
	// The GraphQL engine's Activate function should be called when a plugin is loaded.
	pluginmanager.RegisterPluginLoadedCallback(engine.Activate)

	// It should also be called when the manifest changes, since the manifest can affect function filtering.
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		plugins := pluginmanager.GetRegisteredPlugins()
		if len(plugins) == 0 {
			// No plugins are loaded, so there's nothing to do.
			// This is expected during startup, because the manifest loads before the plugins.
			return nil
		}

		if len(plugins) > 1 {
			// TODO: We should support multiple plugins in the future.
			logger.Warn(ctx).Msg("Multiple plugins loaded.  Only the first plugin will be used.")
		}

		return engine.Activate(ctx, plugins[0])
	})
}

func handleGraphQLRequest(w http.ResponseWriter, r *http.Request) {

	// In dev, redirect non-GraphQL requests to the explorer
	if app.IsDevEnvironment() &&
		r.Method == http.MethodGet &&
		!strings.Contains(r.Header.Get("Accept"), "application/json") {
		http.Redirect(w, r, "/explorer", http.StatusTemporaryRedirect)
		return
	}

	ctx := r.Context()

	// Read the incoming GraphQL request
	var gqlRequest gql.Request
	if err := gql.UnmarshalHttpRequest(r, &gqlRequest); err != nil {
		// TODO: we should capture metrics here
		msg := "Failed to parse GraphQL request."
		http.Error(w, msg, http.StatusBadRequest)

		// NOTE: We only log these in dev, to avoid a bad actor spamming the logs in prod.
		if app.IsDevEnvironment() {
			logger.Warn(ctx, err).Msg(msg)
		}
		return
	}

	// Get the active GraphQL engine, if there is one.
	engine := engine.GetEngine()
	if engine == nil {
		msg := "There is no active GraphQL schema.  Please load a Modus plugin."
		logger.Warn(ctx).Msg(msg)
		w.Header().Set("Content-Type", "application/json")
		if ok, _ := gqlRequest.IsIntrospectionQuery(); ok {
			fmt.Fprint(w, `{"data":{"__schema":{"types":[]}}}`)
		} else {
			fmt.Fprintf(w, `{"errors":[{"message":"%s"}]}`, msg)
		}
		return
	}

	// Create the output map
	output := xsync.NewMap[string, wasmhost.ExecutionInfo]()
	ctx = context.WithValue(ctx, utils.FunctionOutputContextKey, output)

	// Set time zone in the context
	var timeZone string
	if tz := r.Header.Get("X-Time-Zone"); tz != "" {
		// If the X-Time-Zone header is set in the request, use that time zone.
		timeZone = tz
	} else {
		// Otherwise, use the host's local time zone.
		// Note, the TZ environment variable can be set to override the actual local time zone.
		timeZone = timezones.GetLocalTimeZone()
	}
	ctx = context.WithValue(ctx, utils.TimeZoneContextKey, timeZone)

	// Set tracing options
	var options = []eng.ExecutionOptions{}
	if utils.TraceModeEnabled() {
		var traceOpts resolve.TraceOptions
		traceOpts.Enable = true
		traceOpts.IncludeTraceOutputInResponseExtensions = true
		options = append(options, eng.WithRequestTraceOptions(traceOpts))
	}

	// Prepare the result writer
	streaming := false
	resultWriter := gql.NewEngineResultWriter()
	if operationType, err := gqlRequest.OperationType(); err != nil {
		msg := "Failed to determine operation type from GraphQL request."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Error(ctx, err).Msg(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	} else if operationType == gql.OperationTypeSubscription {
		if !isSSERequest(r) {
			msg := "Subscriptions use SSE (Server-Sent Events). Requests must accept 'text/event-stream' for SSE responses."
			logger.Warn(ctx).Msg(msg)
			http.Error(w, msg, http.StatusBadRequest)
			return
		}

		flusher := w.(http.Flusher)
		streaming = true

		// We're following the GraphQL SSE draft spec for subscriptions.  References:
		//   https://the-guild.dev/graphql/sse
		//   https://github.com/enisdenjo/graphql-sse/blob/master/PROTOCOL.md
		//
		// Clients should be implemented similar to these examples:
		//   https://the-guild.dev/graphql/sse/recipes#client-usage

		h := w.Header()
		h.Set("Content-Type", "text/event-stream")
		h.Set("Cache-Control", "no-cache")
		h.Set("Connection", "keep-alive")
		flusher.Flush()

		resultWriter.SetFlushCallback(func(data []byte) {
			// graphql subscription data and errors are pushed over a "next" event, per the graphql-sse draft spec.
			fmt.Fprintf(w, "event: next\ndata: %s\n\n", data)
			flusher.Flush()
		})

		id := xid.New().String()
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		cancelFuncs.Store(id, cancel)
		defer cancelFuncs.Delete(id)
	}

	// Execute the GraphQL operation
	if err := engine.Execute(ctx, &gqlRequest, &resultWriter, options...); err != nil {
		if report, ok := err.(operationreport.Report); ok {
			if len(report.InternalErrors) > 0 {
				// Log internal errors, but don't return them to the client
				msg := "Failed to execute GraphQL operation."
				sentryutils.CaptureError(ctx, err, msg)
				logger.Error(ctx, err).Msg(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		if requestErrors := graphqlerrors.RequestErrorsFromError(err); len(requestErrors) > 0 {
			// TODO: we should capture metrics here

			if streaming {
				fmt.Fprint(w, "event: next\ndata: ")
				_, _ = requestErrors.WriteResponse(w)
				fmt.Fprint(w, "\n\nevent: complete\ndata: \n\n")
			} else {
				w.Header().Set("Content-Type", "application/json")
				_, _ = requestErrors.WriteResponse(w)
			}

			// NOTE: We only log these in dev, to avoid a bad actor spamming the logs in prod.
			if app.IsDevEnvironment() {
				// cleanup empty arrays from error message before logging
				errMsg := strings.Replace(err.Error(), ", locations: []", "", 1)
				errMsg = strings.Replace(errMsg, ", path: []", "", 1)
				logger.Warn(ctx).Str("error", errMsg).Msg("Failed to execute GraphQL operation.")
			}
		} else {
			msg := "Failed to execute GraphQL operation."
			sentryutils.CaptureError(ctx, err, msg)
			logger.Error(ctx, err).Msg(msg)
			http.Error(w, fmt.Sprintf("%s\n%v", msg, err), http.StatusInternalServerError)
		}
		return
	}

	if streaming {
		// In case the connection is still open, we send a final "complete" event, per the graphql-sse draft spec.
		fmt.Fprint(w, "event: complete\ndata: \n\n")
		return
	}

	if response, err := addOutputToResponse(resultWriter.Bytes(), xsync.ToPlainMap(output)); err != nil {
		msg := "Failed to add function output to response."
		sentryutils.CaptureError(ctx, err, msg)
		logger.Error(ctx, err).Msg(msg)
		http.Error(w, fmt.Sprintf("%s\n%v", msg, err), http.StatusInternalServerError)
	} else {

		// An introspection query will always return a Query type, but if only mutations were defined,
		// the fields of the Query type will be null.  That will fail the introspection query, so we need
		// to replace the null with an empty array.
		if ok, _ := gqlRequest.IsIntrospectionQuery(); ok {
			if q := gjson.GetBytes(response, `data.__schema.types.#(name="Query")`); q.Exists() {
				if f := q.Get("fields"); f.Exists() && f.Type == gjson.Null {
					response[f.Index] = '['
					response[f.Index+1] = ']'
					response = slices.Delete(response, f.Index+2, f.Index+4)
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(response)
	}
}

func addOutputToResponse(response []byte, output map[string]wasmhost.ExecutionInfo) ([]byte, error) {

	// NOTE: JSON serialization should be as efficient as possible, as it is called on every GraphQL response.

	jsonOptions := &sjson.Options{
		Optimistic:     true,
		ReplaceInPlace: true,
	}

	var invocations []byte

	for key, item := range output {

		if b, err := sjson.SetBytesOptions(invocations, key+".executionId", item.ExecutionId(), jsonOptions); err != nil {
			return nil, err
		} else {
			invocations = b
		}

		logMessages := utils.TransformConsoleOutput(item.Buffers())
		if len(logMessages) == 0 {
			continue
		}

		i := 0
		for _, logMessage := range logMessages {
			// Only include non-error messages here.
			// Error messages are already included in the response as GraphQL errors.
			if !logMessage.IsError() {
				path := key + ".logs." + strconv.Itoa(i)
				if len(logMessage.Level) > 0 {
					if b, err := sjson.SetBytesOptions(invocations, path+".level", logMessage.Level, jsonOptions); err != nil {
						return nil, err
					} else {
						invocations = b
					}
				}
				if b, err := sjson.SetBytesOptions(invocations, path+".message", logMessage.Message, jsonOptions); err != nil {
					return nil, err
				} else {
					invocations = b
				}
				i++
			}
		}
	}

	if len(invocations) > 0 {
		return sjson.SetRawBytesOptions(response, "extensions.invocations", invocations, jsonOptions)
	}

	return response, nil
}

func isSSERequest(r *http.Request) bool {
	for _, accept := range r.Header.Values("Accept") {
		for value := range strings.SplitSeq(accept, ",") {
			if strings.EqualFold(strings.TrimSpace(strings.SplitN(value, ";", 2)[0]), "text/event-stream") {
				return true
			}
		}
	}
	return false
}

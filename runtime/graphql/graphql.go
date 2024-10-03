/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package graphql

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hypermodeinc/modus/runtime/graphql/engine"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/pluginmanager"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	"github.com/tidwall/sjson"
	eng "github.com/wundergraph/graphql-go-tools/execution/engine"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/graphqlerrors"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/operationreport"
)

var GraphQLRequestHandler = http.HandlerFunc(handleGraphQLRequest)

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

		return engine.Activate(ctx, plugins[0].Metadata)
	})
}

func handleGraphQLRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Read the incoming GraphQL request
	var gqlRequest gql.Request
	err := gql.UnmarshalHttpRequest(r, &gqlRequest)
	if err != nil {
		// NOTE: we intentionally don't log this, to avoid a bad actor spamming the logs
		// TODO: we should capture metrics here though
		msg := "Failed to parse GraphQL request."
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// Get the active GraphQL engine, if there is one.
	engine := engine.GetEngine()
	if engine == nil {
		msg := "There is no active GraphQL schema.  Please load a Modus plugin."
		logger.Warn(ctx).Msg(msg)
		utils.WriteJsonContentHeader(w)
		if ok, _ := gqlRequest.IsIntrospectionQuery(); ok {
			_, _ = w.Write([]byte(`{"data":{"__schema":{"types":[]}}}`))
		} else {
			_, _ = w.Write([]byte(fmt.Sprintf(`{"errors":[{"message":"%s"}]}`, msg)))
		}
		return
	}

	// Create the output map
	output := make(map[string]wasmhost.ExecutionInfo)
	ctx = context.WithValue(ctx, utils.FunctionOutputContextKey, output)

	// Set tracing options
	var options = []eng.ExecutionOptions{}
	if utils.TraceModeEnabled() {
		var traceOpts resolve.TraceOptions
		traceOpts.Enable = true
		traceOpts.IncludeTraceOutputInResponseExtensions = true
		options = append(options, eng.WithRequestTraceOptions(traceOpts))
	}

	// Execute the GraphQL query
	resultWriter := gql.NewEngineResultWriter()
	err = engine.Execute(ctx, &gqlRequest, &resultWriter, options...)
	if err != nil {

		if report, ok := err.(operationreport.Report); ok {
			if len(report.InternalErrors) > 0 {
				// Log internal errors, but don't return them to the client
				msg := "Failed to execute GraphQL query."
				logger.Err(ctx, err).Msg(msg)
				http.Error(w, msg, http.StatusInternalServerError)
				return
			}
		}

		requestErrors := graphqlerrors.RequestErrorsFromError(err)
		if len(requestErrors) > 0 {
			// NOTE: we intentionally don't log this, to avoid a bad actor spamming the logs
			// TODO: we should capture metrics here though
			utils.WriteJsonContentHeader(w)
			_, _ = requestErrors.WriteResponse(w)
		} else {
			msg := "Failed to execute GraphQL query."
			logger.Err(ctx, err).Msg(msg)
			http.Error(w, fmt.Sprintf("%s\n%v", msg, err), http.StatusInternalServerError)
		}
		return
	}

	response := resultWriter.Bytes()
	response, err = addOutputToResponse(response, output)
	if err != nil {
		msg := "Failed to add function output to response."
		logger.Err(ctx, err).Msg(msg)
		http.Error(w, fmt.Sprintf("%s\n%v", msg, err), http.StatusInternalServerError)
	}

	// Return the response
	utils.WriteJsonContentHeader(w)
	_, _ = w.Write(response)
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

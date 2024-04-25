/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"hmruntime/graphql/datasource"
	"hmruntime/graphql/engine"
	"hmruntime/logger"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/buger/jsonparser"
	eng "github.com/wundergraph/graphql-go-tools/execution/engine"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/graphqlerrors"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/operationreport"
)

func Initialize() {
	wasmhost.RegisterPluginLoadedCallback(engine.Activate)
}

func HandleGraphQLRequest(w http.ResponseWriter, r *http.Request) {
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
		msg := "There is no active GraphQL schema.  Please load a Hypermode plugin."
		logger.Warn(ctx).Msg(msg)
		writeJsonContentHeader(w)
		if ok, _ := gqlRequest.IsIntrospectionQuery(); ok {
			w.Write([]byte(`{"data":{"__schema":{"types":[]}}}`))
		} else {
			w.Write([]byte(fmt.Sprintf(`{"errors":[{"message":"%s"}]}`, msg)))
		}
		return
	}

	// Create the output map
	output := map[string]datasource.FunctionOutput{}
	ctx = context.WithValue(ctx, utils.FunctionOutputContextKey, output)

	// Set tracing options
	var options = []eng.ExecutionOptions{}
	if utils.HypermodeTraceEnabled() {
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
			writeJsonContentHeader(w)
			requestErrors.WriteResponse(w)
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
	writeJsonContentHeader(w)
	w.Write(response)
}

func writeJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func addOutputToResponse(response []byte, output map[string]datasource.FunctionOutput) ([]byte, error) {

	type invocationInfo struct {
		ExecutionId string             `json:"executionId"`
		Logs        []utils.LogMessage `json:"logs,omitempty"`
	}

	invocations := make(map[string]invocationInfo, len(output))
	for key, item := range output {
		invocation := invocationInfo{
			ExecutionId: item.ExecutionId,
		}

		l := utils.TransformConsoleOutput(item.Buffers)
		a := make([]utils.LogMessage, 0, len(l))
		for _, m := range l {
			// Only include non-error messages here.
			// Error messages are already included in the response as GraphQL errors.
			if !m.IsError() {
				a = append(a, m)
			}
		}
		if len(a) > 0 {
			invocation.Logs = a
		}

		invocations[key] = invocation
	}

	if len(invocations) == 0 {
		return response, nil
	}

	extensions, jsonType, _, err := jsonparser.Get(response, "extensions")
	if jsonType == jsonparser.NotExist {
		extensions = []byte("{}")
	} else if err != nil {
		return nil, err
	}

	invocationData, err := json.Marshal(invocations)
	if err != nil {
		return nil, err
	}

	extensions, err = jsonparser.Set(extensions, invocationData, "invocations")
	if err != nil {
		return nil, err
	}

	response, err = jsonparser.Set(response, extensions, "extensions")
	if err != nil {
		return nil, err
	}

	return response, nil
}

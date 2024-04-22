/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"hmruntime/graphql/engine"
	"hmruntime/logger"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	"github.com/buger/jsonparser"
	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/graphqlerrors"
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

	// Create the output buffers map
	buffers := map[string]utils.OutputBuffers{}
	ctx = context.WithValue(ctx, utils.FunctionOutputBuffersContextKey, buffers)

	// Execute the GraphQL query
	resultWriter := gql.NewEngineResultWriter()
	err = engine.Execute(ctx, &gqlRequest, &resultWriter)
	if err != nil {
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
	response, err = addLogsToResponse(response, buffers)
	if err != nil {
		msg := "Failed to add logs to response."
		logger.Err(ctx, err).Msg(msg)
		http.Error(w, fmt.Sprintf("%s\n%v", msg, err), http.StatusInternalServerError)
	}

	// Return the response
	writeJsonContentHeader(w)
	w.Write(adjustResponse(response))
}

func writeJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func addLogsToResponse(response []byte, buffers map[string]utils.OutputBuffers) ([]byte, error) {

	logs := make(map[string][]utils.LogMessage, len(buffers))
	for key, buf := range buffers {
		l := utils.TransformConsoleOutput(buf)
		a := make([]utils.LogMessage, 0, len(l))
		for _, m := range l {
			// Only include non-error messages here.
			// Error messages are already included in the response as GraphQL errors.
			if !m.IsError() {
				a = append(a, m)
			}
		}

		if len(a) > 0 {
			logs[key] = a
		}
	}

	if len(logs) == 0 {
		return response, nil
	}

	extensions, jsonType, _, err := jsonparser.Get(response, "extensions")
	if jsonType == jsonparser.NotExist {
		extensions = []byte("{}")
	} else if err != nil {
		return nil, err
	}

	logBytes, err := json.Marshal(logs)
	if err != nil {
		return nil, err
	}

	extensions, err = jsonparser.Set(extensions, logBytes, "logs")
	if err != nil {
		return nil, err
	}

	response, err = jsonparser.Set(response, extensions, "extensions")
	if err != nil {
		return nil, err
	}

	return response, nil
}

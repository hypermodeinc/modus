/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import (
	"fmt"
	"net/http"

	"hmruntime/graphql/engine"
	"hmruntime/logger"
	"hmruntime/wasmhost"

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

	// Execute the GraphQL query
	result := gql.NewEngineResultWriter()
	err = engine.Execute(ctx, &gqlRequest, &result)
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

	// Return the response
	writeJsonContentHeader(w)
	w.Write(adjustResponse(result.Bytes()))
}

func writeJsonContentHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

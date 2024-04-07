/*
 * Copyright 2024 Hypermode, Inc.
 */

package graphql

import (
	"net/http"

	"hmruntime/functions"
	"hmruntime/graphql/engine"
	"hmruntime/logger"

	gql "github.com/wundergraph/graphql-go-tools/execution/graphql"
)

func Initialize() {
	functions.RegisterSchemaLoadedCallback(engine.ActivateSchema)
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

	// Execute the GraphQL query
	engine := engine.GetEngine()
	result := gql.NewEngineResultWriter()
	err = engine.Execute(ctx, &gqlRequest, &result)
	if err != nil {
		msg := "Failed to execute GraphQL query."
		logger.Err(ctx, err).Msg(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	// Return the response
	w.Header().Set("Content-Type", "application/json")
	w.Write(result.Bytes())
}

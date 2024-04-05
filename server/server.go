/*
 * Copyright 2024 Hypermode, Inc.
 */

package server

import (
	"context"
	"fmt"
	"net/http"

	"hmruntime/config"
	"hmruntime/graphql"
	"hmruntime/lambda"
	"hmruntime/logger"

	"github.com/rs/xid"
)

func Start(ctx context.Context) error {
	logger.Info(ctx).
		Int("port", config.Port).
		Msg("Listening for incoming requests.")
	http.HandleFunc("/graphql", graphql.HandleGraphQLRequest)
	http.HandleFunc("/graphql-worker", lambda.HandleDgraphLambdaRequest)
	http.HandleFunc("/admin", handleAdminRequest)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func assignExecutionId(w http.ResponseWriter, r *http.Request) (context.Context, *http.Request) {
	executionId := xid.New().String()
	w.Header().Add("X-Hypermode-ExecutionID", executionId)

	ctx := context.WithValue(r.Context(), logger.ExecutionIdContextKey, executionId)
	return ctx, r.WithContext(ctx)
}

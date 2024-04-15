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
	"hmruntime/logger"
)

func Start(ctx context.Context) error {
	logger.Info(ctx).
		Int("port", config.Port).
		Msg("Listening for incoming requests.")
	http.HandleFunc("/graphql", graphql.HandleGraphQLRequest)
	// http.HandleFunc("/graphql-worker", lambda.HandleDgraphLambdaRequest)
	http.HandleFunc("/admin", handleAdminRequest)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

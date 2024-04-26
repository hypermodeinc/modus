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
	"hmruntime/metrics"
)

func Start(ctx context.Context) error {
	logger.Info(ctx).
		Str("url", fmt.Sprintf("http://localhost:%d/graphql", config.Port)).
		Msg("Listening for incoming requests.")

	http.Handle("/graphql", metrics.InstrumentHandler(graphql.HandleGraphQLRequest, "graphql"))
	http.Handle("/admin", metrics.InstrumentHandler(handleAdminRequest, "admin"))
	http.HandleFunc("/health", healthHandler)
	// register metrics endpoint following prometheus scrape protocol. we do not instrument
	// it with the InstrumentHandler so that scrapper (otel) hitting the server doesn't count.
	http.Handle("/metrics", metrics.MetricsHandler)

	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

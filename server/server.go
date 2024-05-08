/*
 * Copyright 2024 Hypermode, Inc.
 */

package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"hmruntime/config"
	"hmruntime/graphql"
	"hmruntime/logger"
	"hmruntime/metrics"
	"hmruntime/utils"

	"github.com/rs/cors"
)

// shutdownTimeout is the time to wait for the server to shutdown gracefully.
const shutdownTimeout = 5 * time.Second

func Start(ctx context.Context) {

	// Create the configuration for the server.
	mux := GetHandlerMux()
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Port),
		Handler: mux,
	}

	// Start the server in a goroutine so we can listen for signals to shutdown.
	shutdownChan := make(chan bool, 1)
	go func() {
		logger.Info(ctx).
			Str("url", fmt.Sprintf("http://localhost:%d/graphql", config.Port)).
			Msg("Listening for incoming requests.")

		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ctx).Err(err).Msg("HTTP server error.  Exiting.")
		}

		logger.Info(ctx).Msg("Shutdown requested.  Stopping HTTP server...")

		// TODO: If we have other goroutines we want to stop gracefully, we can do so here.

		// Signal shutdown is complete.
		shutdownChan <- true
	}()

	// Wait for a signal to shutdown the server.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Shutdown the server gracefully.
	shutdownCtx, shutdownRelease := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdownRelease()
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("HTTP server shutdown error.")
	}

	// Wait for the server to shutdown completely.
	<-shutdownChan
	logger.Info(ctx).Msg("Shutdown complete.")
}

func GetHandlerMux() http.Handler {
	mux := http.NewServeMux()

	// Register our main endpoints with instrumentation.
	mux.Handle("/graphql", metrics.InstrumentHandler(graphql.HandleGraphQLRequest, "graphql"))
	mux.Handle("/admin", metrics.InstrumentHandler(handleAdminRequest, "admin"))

	// Register metrics endpoint which uses the Prometheus scraping protocol.
	// We do not instrument it with the InstrumentHandler so that any scraper (eg. OTel)
	// hitting the server doesn't count.
	mux.Handle("/metrics", metrics.MetricsHandler)

	// Also register the health endpoint, un-instrumented.
	mux.HandleFunc("/health", healthHandler)

	// Restrict the HTTP methods for all above handlers to GET and POST.
	handler := restrictHttpMethods(mux)

	// Add CORS support to all endpoints.
	c := cors.New(cors.Options{
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	return c.Handler(handler)
}

func restrictHttpMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodPost:
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	env := config.GetEnvironmentName()
	ver := config.GetVersionNumber()
	w.WriteHeader(http.StatusOK)
	utils.WriteJsonContentHeader(w)
	_, _ = w.Write([]byte(`{"status":"ok","environment":"` + env + `","version":"` + ver + `"}`))
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package httpserver

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/graphql"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/metrics"
	"github.com/hypermodeinc/modus/runtime/middleware"
	"github.com/rs/cors"
)

// shutdownTimeout is the time to wait for the server to shutdown gracefully.
const shutdownTimeout = 5 * time.Second

func Start(ctx context.Context, local bool) {
	if local {
		// If we are running locally, only listen on localhost.
		// This prevents getting nagged for firewall permissions each launch.
		// Listen on IPv4, and also on IPv6 if available.
		addresses := []string{fmt.Sprintf("127.0.0.1:%d", config.Port)}
		if isIPv6Available() {
			addresses = append(addresses, fmt.Sprintf("[::1]:%d", config.Port))
		}
		startHttpServer(ctx, addresses...)
	} else {
		// Otherwise, listen on all interfaces.
		addr := fmt.Sprintf(":%d", config.Port)
		startHttpServer(ctx, addr)
	}
}

func startHttpServer(ctx context.Context, addresses ...string) {

	// Setup a server for each address.
	mux := GetMainHandler()
	servers := make([]*http.Server, len(addresses))
	for i, addr := range addresses {
		servers[i] = &http.Server{Handler: mux, Addr: addr}
	}

	// Start a goroutine for each server.
	shutdownChan := make(chan bool, len(addresses))
	for _, server := range servers {
		go func() {
			if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Fatal(ctx).Err(err).Msg("HTTP server error.  Exiting.")
			}
			shutdownChan <- true
		}()
	}

	// Wait for a signal to shutdown the servers.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info(ctx).Msg("Context canceled.  Stopping HTTP server...")
	case sig := <-sigChan:
		switch sig {
		case syscall.SIGINT:
			fmt.Print("\b\b") // erase the ^C
			logger.Info(ctx).Msg("Interrupt signal received.  Stopping HTTP server...")
		case syscall.SIGTERM:
			logger.Info(ctx).Msg("Terminate signal received.  Stopping HTTP server...")
		}
	}

	// Shutdown all servers gracefully.
	for _, server := range servers {
		shutdownCtx, shutdownRelease := context.WithTimeout(ctx, shutdownTimeout)
		defer shutdownRelease()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Fatal(ctx).Err(err).Msg("HTTP server shutdown error.")
		}
	}

	// Wait for the servers to shutdown completely.
	for range servers {
		<-shutdownChan
	}

	logger.Info(ctx).Msg("Shutdown complete.")
}

func WithDefaultGraphQLHandler() func(routes map[string]http.Handler) {
	return func(routes map[string]http.Handler) {
		routes["/graphql"] = metrics.InstrumentHandler(graphql.GraphQLRequestHandler, "default")
	}
}

func GetMainHandler(options ...func(map[string]http.Handler)) http.Handler {

	// Create default routes.
	defaultRoutes := map[string]http.Handler{
		"/health":  healthHandler,
		"/metrics": metrics.MetricsHandler,
	}
	for _, opt := range options {
		opt(defaultRoutes)
	}

	// Create a dynamic mux to handle the routing.
	mux := newDynamicMux(defaultRoutes)

	// Dynamically add routes as they are loaded from the manifest.
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		routes := maps.Clone(defaultRoutes)

		m := manifestdata.GetManifest()
		for name, ep := range m.Endpoints {
			switch ep.EndpointType() {
			case manifest.EndpointTypeGraphQL:
				info := ep.(manifest.GraphqlEndpointInfo)
				var handler http.Handler = graphql.GraphQLRequestHandler

				switch info.Auth {
				case manifest.EndpointAuthNone:
					// No auth required.
				case manifest.EndpointAuthBearerToken:
					handler = middleware.HandleJWT(handler)
				default:
					logger.Warn(ctx).Str("endpoint", name).Msg("Unsupported auth type.")
					continue
				}

				routes[info.Path] = metrics.InstrumentHandler(handler, name)

				logger.Info(ctx).
					Str("url", fmt.Sprintf("http://localhost:%d%s", config.Port, info.Path)).
					Msg("Registered GraphQL endpoint.")

			default:
				logger.Warn(ctx).Str("endpoint", name).Msg("Unsupported endpoint type.")
			}
		}

		mux.ReplaceRoutes(routes)

		return nil
	})

	// Restrict the HTTP methods for all handlers to GET and POST.
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

func isIPv6Available() bool {
	addr := &net.UDPAddr{IP: net.ParseIP("::1")}
	conn, err := net.ListenUDP("udp6", addr)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

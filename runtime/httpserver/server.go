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
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/explorer"
	"github.com/hypermodeinc/modus/runtime/graphql"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/metrics"
	"github.com/hypermodeinc/modus/runtime/middleware"

	"github.com/fatih/color"
	"github.com/rs/cors"
)

var titleColor = color.New(color.FgHiGreen, color.Bold)
var itemColor = color.New(color.FgHiBlue)
var urlColor = color.New(color.FgHiCyan)
var noticeColor = color.New(color.FgGreen, color.Italic)
var warningColor = color.New(color.FgYellow)

// ShutdownTimeout is the time to wait for the server to shutdown gracefully.
const shutdownTimeout = 5 * time.Second

func Start(ctx context.Context, local bool) {

	port := app.Config().Port()

	if local {
		// If we are running locally, only listen on localhost.
		// This prevents getting nagged for firewall permissions each launch.
		// Listen on IPv4, and also on IPv6 if available.
		addresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}
		if isIPv6Available() {
			addresses = append(addresses, fmt.Sprintf("[::1]:%d", port))
		}
		startHttpServer(ctx, addresses...)
	} else {
		// Otherwise, listen on all interfaces.
		addr := fmt.Sprintf(":%d", port)
		startHttpServer(ctx, addr)
	}
}

func startHttpServer(ctx context.Context, addresses ...string) {

	// Get the main handler for the server.
	// Note: This must be done first, because it registers for callback events.
	mux := GetMainHandler()

	// Initialize our middleware before starting the server.
	middleware.Init(ctx)

	// Setup a server for each address.
	servers := make([]*http.Server, len(addresses))
	for i, addr := range addresses {
		servers[i] = &http.Server{Handler: mux, Addr: addr}
	}

	// Start a goroutine for each server.
	shutdownChan := make(chan bool, len(addresses))
	for _, server := range servers {
		go func() {
			err := server.ListenAndServe()
			app.SetShuttingDown()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
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

	cfg := app.Config()
	if cfg.IsDevEnvironment() {
		defaultRoutes["/explorer/"] = explorer.ExplorerHandler
		defaultRoutes["/"] = http.RedirectHandler("/explorer/", http.StatusSeeOther)
	}
	port := cfg.Port()

	for _, opt := range options {
		opt(defaultRoutes)
	}

	// Create a dynamic mux to handle the routing.
	mux := newDynamicMux(defaultRoutes)

	// Dynamically add routes as they are loaded from the manifest.
	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		routes := maps.Clone(defaultRoutes)

		type endpoint struct {
			apiType string
			name    string
			url     string
		}

		var endpoints []endpoint

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

				url := fmt.Sprintf("http://localhost:%d%s", port, info.Path)
				logger.Info(ctx).Str("url", url).Msg("Registered GraphQL endpoint.")
				endpoints = append(endpoints, endpoint{"GraphQL", name, url})

			default:
				logger.Warn(ctx).Str("endpoint", name).Msg("Unsupported endpoint type.")
			}
		}

		mux.ReplaceRoutes(routes)

		if app.Config().IsDevEnvironment() {
			fmt.Fprintln(os.Stderr)

			switch len(endpoints) {
			case 0:
				warningColor.Fprintln(os.Stderr, "No local endpoints are configured.")
				warningColor.Fprintln(os.Stderr, "Please add one or more endpoints to your modus.json file.")
			case 1:
				ep := endpoints[0]
				titleColor.Fprintln(os.Stderr, "Your local endpoint is ready!")
				itemColor.Fprintf(os.Stderr, "• %s (%s): ", ep.apiType, ep.name)
				urlColor.Fprintln(os.Stderr, ep.url)

				explorerURL := fmt.Sprintf("http://localhost:%d/explorer", port)
				titleColor.Fprintf(os.Stderr, "\nView endpoint: ")
				urlColor.Fprintln(os.Stderr, explorerURL)

			default:
				titleColor.Fprintln(os.Stderr, "Your local endpoints are ready!")
				for _, ep := range endpoints {
					itemColor.Fprintf(os.Stderr, "• %s (%s): ", ep.apiType, ep.name)
					urlColor.Fprintln(os.Stderr, ep.url)
				}

				explorerURL := fmt.Sprintf("http://localhost:%d/explorer", port)
				titleColor.Fprintf(os.Stderr, "\nView your endpoints at: ")
				urlColor.Fprintln(os.Stderr, explorerURL)
			}

			fmt.Fprintln(os.Stderr)
			noticeColor.Fprintln(os.Stderr, "Changes will automatically be applied when you save your files.")
			noticeColor.Fprintln(os.Stderr, "Press Ctrl+C at any time to stop the server.")
			fmt.Fprintln(os.Stderr)
		}

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

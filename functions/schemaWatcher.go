/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"errors"
	"net/url"
	"time"

	"hmruntime/config"
	"hmruntime/dgraph"
	"hmruntime/host"
	"hmruntime/logger"
)

// Holds the current GraphQL schema
var gqlSchema string

func MonitorGqlSchema(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		var urlErrorLogged, schemaErrorLogged bool
		for {
			schema, err := dgraph.GetGQLSchema(ctx)
			if err != nil {
				var urlErr *url.Error
				if errors.As(err, &urlErr) && !urlErrorLogged {
					logger.Err(ctx, urlErr).Msg("Failed to connect to Dgraph.")
					urlErrorLogged = true
				} else if !schemaErrorLogged {
					logger.Err(ctx, err).Msg("Failed to retrieve GraphQL schema.")
					schemaErrorLogged = true
				}
			} else if schema != gqlSchema {
				urlErrorLogged = false
				schemaErrorLogged = false
				if gqlSchema == "" {
					logger.Info(ctx).Msg("Schema loaded.")
				} else {
					logger.Info(ctx).Msg("Schema changed.")
				}

				// Signal that we need to register functions
				host.RegistrationRequest <- true

				gqlSchema = schema
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				return
			}
		}
	}()
}

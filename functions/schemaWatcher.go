/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"errors"
	"hmruntime/config"
	"hmruntime/dgraph"
	"hmruntime/host"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
)

// Holds the current GraphQL schema
var gqlSchema string

func MonitorGqlSchema(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		for {
			schema, err := dgraph.GetGQLSchema(ctx)
			if err != nil {
				var urlErr *url.Error
				if errors.As(err, &urlErr) {
					log.Err(urlErr).Msg("Failed to connect to Dgraph.")
				} else {
					log.Err(err).Msg("Failed to retrieve GraphQL schema.")
				}
			} else if schema != gqlSchema {
				if gqlSchema == "" {
					log.Info().Msg("Schema loaded.")
				} else {
					log.Info().Msg("Schema changed.")
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

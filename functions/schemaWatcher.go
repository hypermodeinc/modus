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
	"hmruntime/host"
	"hmruntime/logger"
	"hmruntime/storage"
)

// Holds the current GraphQL schema
var gqlSchema string

func GqlSchema() string {
	return gqlSchema
}

var schemaLoaded func(ctx context.Context, schema string) error

func RegisterSchemaLoadedCallback(callback func(ctx context.Context, schema string) error) {
	schemaLoaded = callback
}

func MonitorGqlSchema(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(config.RefreshInterval)
		defer ticker.Stop()

		var urlErrorLogged, schemaErrorLogged bool
		for {
			// schema, err := dgraph.GetGQLSchema(ctx)
			schema, err := getSchema(ctx)
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

				err := schemaLoaded(ctx, schema)
				if err != nil {
					logger.Err(ctx, err).Msg("GraphQL schema callback failed.")
				}

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

func getSchema(ctx context.Context) (string, error) {
	contents, err := storage.GetFileContents(ctx, "api.graphql")
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

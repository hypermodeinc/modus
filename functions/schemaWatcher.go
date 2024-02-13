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
	"log"
	"net/url"
	"time"
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
					log.Printf("Failed to connect to Dgraph: %v", urlErr)
				} else {
					log.Printf("Failed to retrieve GraphQL schema: %v", err)
				}
			} else if schema != gqlSchema {
				if gqlSchema == "" {
					log.Printf("Schema loaded")
				} else {
					log.Printf("Schema changed")
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

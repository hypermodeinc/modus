/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"errors"
	"log"
	"net/url"
	"time"
)

// Polling interval to check Dgraph for GraphQL schema changes
const schemaRefreshInterval time.Duration = time.Second * 5

// Holds the current GraphQL schema
var gqlSchema string

func monitorGqlSchema(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(schemaRefreshInterval)
		defer ticker.Stop()

		for {
			schema, err := getGQLSchema(ctx)
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
				register <- true

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

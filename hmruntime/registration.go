/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"context"
	"log"
)

var register chan bool = make(chan bool)

func monitorRegistration(ctx context.Context) {
	go func() {
		for {
			select {
			case <-register:
				log.Printf("Registering functions")
				err := registerFunctions(gqlSchema)
				if err != nil {
					log.Printf("Failed to register functions: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

/*
 * Copyright 2023 Hypermode, Inc.
 */

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"hmruntime/config"
	"hmruntime/functions"
	"hmruntime/host"
	"hmruntime/logger"

	"github.com/rs/xid"
)

type HMRequest struct {
	AccessToken string `json:"X-Dgraph-AccessToken"`
	AuthHeader  struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"authHeader"`
	Args     map[string]any   `json:"args"`
	Parents  []map[string]any `json:"parents"`
	Resolver string           `json:"resolver"`
}

type HMErrorResponse struct {
	Errors []HMError `json:"errors"`
}

type HMError struct {
	Message string `json:"message"`
}

type AdminRequest struct {
	Action string `json:"action"`
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	// Assign an Execution ID to the request context
	ctx, r := assignExecutionId(w, r)

	// Decode the request body
	var req HMRequest
	dec := json.NewDecoder(r.Body)
	dec.UseNumber()
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Err(ctx, err).Msg("Failed to decode request body.")
		return
	}

	// Get the function info for the resolver
	info, ok := functions.FunctionsMap[req.Resolver]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(ctx).
			Str("resolver", req.Resolver).
			Msg("No function registered for resolver.")
		return
	}

	// Add plugin details to the context
	ctx = context.WithValue(ctx, logger.PluginNameContextKey, info.Plugin.Name())
	ctx = context.WithValue(ctx, logger.BuildIdContextKey, info.Plugin.BuildId())

	// Get a module instance for this request.
	// Each request will get its own instance of the plugin module,
	// so that we can run multiple requests in parallel without risk
	// of corrupting the module's memory.
	mod, buf, err := host.GetModuleInstance(ctx, info.Plugin)
	if err != nil {
		logger.Err(ctx, err).Msg("Failed to get module instance.")
		err := writeErrorResponse(w, err)
		if err != nil {
			logger.Err(ctx, err).Msg("Failed to write error response.")
		}
		return
	}
	defer mod.Close(ctx)

	fnName := info.FunctionName()
	if req.Args != nil {

		// Call the function, passing in the args from the request
		result, err := functions.CallFunction(ctx, mod, info, req.Args)
		if err != nil {
			err := fmt.Errorf("error calling function '%s': %w", fnName, err)
			err = writeErrorResponse(w, err, buf.Stdout.String(), buf.Stderr.String())
			if err != nil {
				logger.Err(ctx, err).Msg("Failed to write error response.")
			}
			return
		}

		// Handle no result due to void return type
		if result == nil {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Determine if the result is already JSON
		isJson := false
		fieldType := info.Schema.FieldDef.Type.NamedType
		if _, ok := result.(string); ok && fieldType != "String" {
			isJson = true
		}

		// Write the result
		err = writeDataAsJson(w, result, isJson)
		if err != nil {
			logger.Err(ctx, err).Msg("Failed to write result data to response stream.")
		}

	} else if req.Parents != nil {

		results := make([]any, len(req.Parents))

		// Call the function for each parent
		for i, parent := range req.Parents {
			results[i], err = functions.CallFunction(ctx, mod, info, parent)
			if err != nil {
				err := fmt.Errorf("error calling function '%s': %w", fnName, err)
				err = writeErrorResponse(w, err, buf.Stdout.String(), buf.Stderr.String())
				if err != nil {
					logger.Err(ctx, err).Msg("Failed to write error response.")
				}
				return
			}
		}

		// Write the result
		isJson := info.Schema.FieldDef.Type.NamedType == ""
		err = writeDataAsJson(w, results, isJson)
		if err != nil {
			logger.Err(ctx, err).Msg("Failed to write result data to response stream.")
		}

	} else {
		w.WriteHeader(http.StatusBadRequest)
		logger.Error(ctx).Msg("Request must have either args or parents.")
	}
}

func writeErrorResponse(w http.ResponseWriter, err error, msgs ...string) error {
	w.WriteHeader(http.StatusInternalServerError)

	// Dgraph lambda expects a JSON response similar to a GraphQL error response
	w.Header().Set("Content-Type", "application/json")
	resp := HMErrorResponse{Errors: []HMError{}}

	// Emit messages first
	for _, msg := range msgs {
		for _, line := range strings.Split(msg, "\n") {
			if len(line) > 0 {
				resp.Errors = append(resp.Errors, HMError{Message: line})
			}
		}
	}

	// Emit the error last
	resp.Errors = append(resp.Errors, HMError{Message: err.Error()})

	return json.NewEncoder(w).Encode(resp)
}

func writeDataAsJson(w http.ResponseWriter, data any, isJson bool) error {
	if isJson {
		switch data := data.(type) {
		case string:
			w.Header().Set("Content-Type", "application/json")
			n, err := w.Write([]byte(data))
			if err != nil || n != len(data) {
				return fmt.Errorf("failed to write result data: %w", err)
			}
		case []string:
			w.Header().Set("Content-Type", "application/json")
			n, err := w.Write([]byte{'['})
			if err != nil || n != 1 {
				return fmt.Errorf("failed to write result data: %w", err)
			}
			for i, s := range data {
				if i > 0 {
					n, err = w.Write([]byte{','})
					if err != nil || n != 1 {
						return fmt.Errorf("failed to write result data: %w", err)
					}
				}
				n, err = w.Write([]byte(s))
				if err != nil || n != len(s) {
					return fmt.Errorf("failed to write result data: %w", err)
				}
			}
			n, err = w.Write([]byte{']'})
			if err != nil || n != 1 {
				return fmt.Errorf("failed to write result data: %w", err)
			}
		default:
			return fmt.Errorf("unexpected result type: %T", data)
		}

		return nil
	}

	output, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to serialize result data: %s", err)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write result data: %s", err)
	}

	return nil
}

func handleAdminRequest(w http.ResponseWriter, r *http.Request) {

	// Assign an Execution ID to the request context
	ctx, r := assignExecutionId(w, r)

	// Decode the request body
	var req AdminRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logger.Err(ctx, err).Msg("Failed to decode request body.")
		return
	}

	// Perform the requested action
	switch req.Action {
	// TODO: Add admin actions here
	default:
		err = fmt.Errorf("unknown action: %s", req.Action)
	}

	// Write the response
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logger.Err(ctx, err).Msg("Failed to perform admin action.")
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func startServer(ctx context.Context) error {
	logger.Info(ctx).
		Int("port", config.Port).
		Msg("Listening for incoming requests.")
	http.HandleFunc("/graphql-worker", handleRequest)
	http.HandleFunc("/admin", handleAdminRequest)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

func assignExecutionId(w http.ResponseWriter, r *http.Request) (context.Context, *http.Request) {
	executionId := xid.New().String()
	w.Header().Add("X-Hypermode-ExecutionID", executionId)

	ctx := context.WithValue(r.Context(), logger.ExecutionIdContextKey, executionId)
	return ctx, r.WithContext(ctx)
}

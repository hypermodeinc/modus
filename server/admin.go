/*
 * Copyright 2024 Hypermode, Inc.
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"hmruntime/logger"
)

type AdminRequest struct {
	Action string `json:"action"`
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

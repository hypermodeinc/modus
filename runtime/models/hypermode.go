/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/secrets"
)

var _hypermodeModelHost string

func getHypermodeModelEndpointUrl(model *manifest.ModelInfo) (string, error) {
	// In development, use the shared Hypermode model server.
	// Note: Authentication via the Hypermode CLI is required.
	if app.IsDevEnvironment() {
		endpoint := fmt.Sprintf("https://models.hypermode.host/%s", strings.ToLower(model.SourceModel))
		return endpoint, nil
	}

	// In production, use the Hypermode internal model endpoint.
	// Access is protected by the Hypermode internal network.
	if _hypermodeModelHost == "" {
		_hypermodeModelHost = os.Getenv("HYPERMODE_MODEL_HOST")
		if _hypermodeModelHost == "" {
			return "", fmt.Errorf("hypermode hosted models are not available in this environment")
		}
	}
	endpoint := fmt.Sprintf("http://%s.%s/%[1]s:predict", strings.ToLower(model.Name), _hypermodeModelHost)
	return endpoint, nil
}

func authenticateHypermodeModelRequest(ctx context.Context, req *http.Request, connection *manifest.HTTPConnectionInfo) error {
	// In development, Hypermode models require authentication.
	if app.IsDevEnvironment() {
		return secrets.ApplyAuthToLocalHypermodeModelRequest(ctx, connection, req)
	}

	// In production, the Hypermode infrastructure protects the model server.
	return nil
}

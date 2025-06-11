/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/fatih/color"
	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
)

var provider secretsProvider
var errLocalAuthFailed = fmt.Errorf("local authentication failed")

type secretsProvider interface {
	initialize(ctx context.Context)
	hasSecret(ctx context.Context, name string) bool
	getSecretValue(ctx context.Context, name string) (string, error)
	getConnectionSecrets(ctx context.Context, connection manifest.ConnectionInfo) (map[string]string, error)
}

func Initialize(ctx context.Context) {
	if app.Config().UseKubernetesSecret() {
		provider = &kubernetesSecretsProvider{}
	} else {
		provider = &localSecretsProvider{}
	}
	provider.initialize(ctx)
}

func HasSecret(ctx context.Context, name string) bool {
	return provider.hasSecret(ctx, name)
}

func GetSecretValue(ctx context.Context, name string) (string, error) {
	return provider.getSecretValue(ctx, name)
}

func GetConnectionSecrets(ctx context.Context, connection manifest.ConnectionInfo) (map[string]string, error) {
	return provider.getConnectionSecrets(ctx, connection)
}

func GetConnectionSecret(ctx context.Context, connection manifest.ConnectionInfo, secretName string) (string, error) {
	secrets, err := GetConnectionSecrets(ctx, connection)
	if err != nil {
		return "", err
	}

	if val, ok := secrets[secretName]; ok {
		return val, nil
	}

	return "", fmt.Errorf("could not find secret '%s' for connection '%s'", secretName, connection.ConnectionName())
}

// ApplySecretsToHttpRequest evaluates the given request and replaces any placeholders
// present in the query parameters and headers with their secret values for the given connection.
func ApplySecretsToHttpRequest(ctx context.Context, connection *manifest.HTTPConnectionInfo, req *http.Request) error {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	// get secrets for the connection
	secrets, err := GetConnectionSecrets(ctx, connection)
	if err != nil {
		return err
	}

	// apply query parameters from manifest
	q := req.URL.Query()
	for k, v := range connection.QueryParameters {
		q.Add(k, applySecretsToString(ctx, secrets, v))
	}
	req.URL.RawQuery = q.Encode()

	// apply headers from manifest
	for k, v := range connection.Headers {
		req.Header.Add(k, applySecretsToString(ctx, secrets, v))
	}

	return nil
}

func ApplyAuthToLocalHypermodeModelRequest(ctx context.Context, connection manifest.ConnectionInfo, req *http.Request) error {

	apiKey := os.Getenv("HYP_API_KEY")
	workspaceId := os.Getenv("HYP_WORKSPACE_ID")

	warningColor := color.New(color.FgHiYellow, color.Bold)

	if apiKey == "" || workspaceId == "" {
		fmt.Fprintln(os.Stderr)
		warningColor.Fprintln(os.Stderr, "Warning: Local authentication not found. Please login using `hyp login`")
		fmt.Fprintln(os.Stderr)
		return errLocalAuthFailed
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("HYP-WORKSPACE-ID", workspaceId)

	return nil
}

// ApplySecretsToString evaluates the given string and replaces any placeholders
// present in the string with their secret values for the given connection.
func ApplySecretsToString(ctx context.Context, connection manifest.ConnectionInfo, str string) (string, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	secrets, err := GetConnectionSecrets(ctx, connection)
	if err != nil {
		return "", err
	}
	return applySecretsToString(ctx, secrets, str), nil
}

var templateRegex = regexp.MustCompile(`{{\s*(?:base64\((.+?):(.+?)\)|(.+?))\s*}}`)

func applySecretsToString(ctx context.Context, secrets map[string]string, s string) string {
	return templateRegex.ReplaceAllStringFunc(s, func(match string) string {
		submatches := templateRegex.FindStringSubmatch(match)
		if len(submatches) != 4 {
			logger.Warn(ctx).Str("match", match).Msg("Invalid template.")
			return match
		}

		// standard secret template
		nameKey := submatches[3]
		if nameKey != "" {
			val, ok := secrets[nameKey]
			if ok {
				return val
			}

			logger.Warn(ctx).Str("name", nameKey).Msg("Secret not found.")
			return ""
		}

		// base64 secret template
		userKey := submatches[1]
		passKey := submatches[2]
		if userKey != "" && passKey != "" {
			user, userOk := secrets[userKey]
			if !userOk {
				logger.Warn(ctx).Str("name", userKey).Msg("Secret not found.")
			}

			pass, passOk := secrets[passKey]
			if !passOk {
				logger.Warn(ctx).Str("name", passKey).Msg("Secret not found.")
			}

			if userOk && passOk {
				return base64.StdEncoding.EncodeToString([]byte(user + ":" + pass))
			} else {
				return ""
			}
		}

		logger.Warn(ctx).Str("template", match).Msg("Invalid secret variable template.")
		return match
	})
}

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

	"github.com/hypermodeinc/modus/lib/manifest"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"
)

var provider secretsProvider

type secretsProvider interface {
	initialize(ctx context.Context)
	hasSecret(name string) bool
	getSecretValue(name string) (string, error)
	getConnectionSecrets(connection manifest.ConnectionInfo) (map[string]string, error)
}

func Initialize(ctx context.Context) {
	provider = &localSecretsProvider{}
	provider.initialize(ctx)
}

func HasSecret(name string) bool {
	return provider.hasSecret(name)
}

func GetSecretValue(name string) (string, error) {
	return provider.getSecretValue(name)
}

func GetConnectionSecrets(connection manifest.ConnectionInfo) (map[string]string, error) {
	return provider.getConnectionSecrets(connection)
}

func GetConnectionSecret(connection manifest.ConnectionInfo, secretName string) (string, error) {
	secrets, err := GetConnectionSecrets(connection)
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
	secrets, err := GetConnectionSecrets(connection)
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

	jwt := os.Getenv("HYP_JWT")
	orgId := os.Getenv("HYP_ORG_ID")

	if jwt == "" || orgId == "" {
		return fmt.Errorf("missing HYP_JWT or HYP_ORG_ID environment variables, login to Hypermode using 'hyp login'")
	}

	req.Header.Set("Authorization", "Bearer "+jwt)
	req.Header.Set("HYP-ORG-ID", orgId)

	return nil
}

// ApplySecretsToString evaluates the given string and replaces any placeholders
// present in the string with their secret values for the given connection.
func ApplySecretsToString(ctx context.Context, connection manifest.ConnectionInfo, str string) (string, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	secrets, err := GetConnectionSecrets(connection)
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

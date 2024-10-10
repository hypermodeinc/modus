/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package middleware

import (
	"context"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"net/http"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaimsKey string

const jwtClaims jwtClaimsKey = "claims"

var authPublicKeys map[string]any

func Init(ctx context.Context) {
	publicKeysJson := os.Getenv("MODUS_PEMS")
	if publicKeysJson == "" {
		return
	}
	var publicKeyStrings map[string]string
	err := json.Unmarshal([]byte(publicKeysJson), &publicKeyStrings)
	if err != nil {
		if config.IsDevEnvironment() {
			logger.Fatal(ctx).Err(err).Msg("Auth public keys deserializing error")
		}
		logger.Error(ctx).Err(err).Msg("Auth public keys deserializing error")
		return
	}
	authPublicKeys = make(map[string]any)
	for key, value := range publicKeyStrings {
		block, _ := pem.Decode([]byte(value))
		if block == nil {
			if config.IsDevEnvironment() {
				logger.Fatal(ctx).Msg("Invalid PEM block for key: " + key)
			}
			logger.Error(ctx).Msg("Invalid PEM block for key: " + key)
			return
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			if config.IsDevEnvironment() {
				logger.Fatal(ctx).Err(err).Msg("JWT public key parsing error for key: " + key)
			}
			logger.Error(ctx).Err(err).Msg("JWT public key parsing error for key: " + key)
			return
		}
		authPublicKeys[key] = pubKey
	}
}

func HandleJWT(next http.Handler) http.Handler {
	var jwtParser = jwt.NewParser()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context = r.Context()
		tokenStr := r.Header.Get("Authorization")
		if tokenStr != "" {
			if s, found := strings.CutPrefix(tokenStr, "Bearer "); found {
				tokenStr = s
			} else {
				logger.Error(ctx).Msg("Invalid JWT token format, Bearer required")
				http.Error(w, "Invalid JWT token format, Bearer required", http.StatusUnauthorized)
				return
			}
		}

		if config.IsDevEnvironment() {
			if len(authPublicKeys) == 0 {
				if tokenStr == "" {
					next.ServeHTTP(w, r)
					return
				}
				token, _, err := jwtParser.ParseUnverified(tokenStr, jwt.MapClaims{})
				if err != nil {
					logger.Warn(ctx).Err(err).Msg("Error parsing JWT token. Continuing since running in development")
					next.ServeHTTP(w, r)
					return
				}
				if claims, ok := token.Claims.(jwt.MapClaims); ok {
					ctx = addClaimsToContext(ctx, claims)
				}
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
		} else {
			if len(authPublicKeys) == 0 {
				next.ServeHTTP(w, r)
				return
			}
		}

		var token *jwt.Token
		var err error

		for _, publicKey := range authPublicKeys {
			token, err = jwtParser.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return publicKey, nil
			})
			if err == nil {
				logger.Info(ctx).Msg("JWT token parsed successfully")
				break
			}
		}
		if err != nil {
			if config.IsDevEnvironment() {
				logger.Debug(r.Context()).Err(err).Msg("JWT parse error")
				next.ServeHTTP(w, r)
				return
			}
			logger.Error(r.Context()).Err(err).Msg("JWT parse error")
			http.Error(w, "Invalid JWT token", http.StatusUnauthorized)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			ctx = addClaimsToContext(ctx, claims)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addClaimsToContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	claimsJson, err := utils.JsonSerialize(claims)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("JWT claims serialization error")
		return ctx
	}
	return context.WithValue(ctx, jwtClaims, string(claimsJson))
}

func GetJWTClaims(ctx context.Context) string {
	if claims, ok := ctx.Value(jwtClaims).(string); ok {
		return claims
	}
	return ""
}

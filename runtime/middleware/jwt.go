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
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/runtime/config"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/golang-jwt/jwt/v5"
)

type jwtClaimsKey string

const jwtClaims jwtClaimsKey = "jwt_claims"

var PrivKeys map[string]string

func Init() {
	privKeysStr := os.Getenv("MODUS_RSA_PEMS")
	if privKeysStr == "" {
		return
	}
	err := json.Unmarshal([]byte(privKeysStr), &PrivKeys)
	if err != nil {
		logger.Error(context.Background()).Err(err).Msg("JWT private keys unmarshalling error")
		return
	}
}

func HandleJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ctx context.Context = r.Context()
		tokenStr := r.Header.Get("Authorization")
		trimmedTokenStr := strings.TrimPrefix(tokenStr, "Bearer ")

		if trimmedTokenStr == tokenStr {
			logger.Error(ctx).Msg("Invalid JWT token format, Bearer required")
			http.Error(w, "Invalid JWT token format, Bearer required", http.StatusUnauthorized)
			return
		}

		if PrivKeys == nil {
			if !config.IsDevEnvironment() || trimmedTokenStr == "" {
				next.ServeHTTP(w, r)
				return
			}
			token, _, err := new(jwt.Parser).ParseUnverified(trimmedTokenStr, jwt.MapClaims{})
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

		var token *jwt.Token
		var err error

		for _, privKey := range PrivKeys {
			token, err = jwt.Parse(trimmedTokenStr, func(token *jwt.Token) (interface{}, error) {
				return jwt.ParseRSAPublicKeyFromPEM([]byte(privKey))
			})
			if err != nil {
				if config.IsDevEnvironment() {
					logger.Debug(r.Context()).Err(err).Msg("JWT parse error")
					next.ServeHTTP(w, r)
					return
				}
				logger.Error(r.Context()).Err(err).Msg("JWT parse error")
				continue
			} else {
				break
			}
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

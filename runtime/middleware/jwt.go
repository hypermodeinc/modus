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

const jwtClaims jwtClaimsKey = "jwt_claims"

var jwtPublicKeys map[string]any

func Init(ctx context.Context) {
	privKeysJson := os.Getenv("MODUS_RSA_PEMS")
	if privKeysJson == "" {
		return
	}
	var publicKeyStrings map[string]string
	err := json.Unmarshal([]byte(privKeysJson), &publicKeyStrings)
	if err != nil {
		logger.Error(ctx).Err(err).Msg("JWT private keys unmarshalling error")
		return
	}
	jwtPublicKeys = make(map[string]any)
	for key, value := range publicKeyStrings {
		block, _ := pem.Decode([]byte(value))
		if block == nil {
			logger.Error(ctx).Msg("Invalid PEM block")
			return
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			logger.Error(ctx).Err(err).Msg("JWT public key parsing error")
			return
		}
		jwtPublicKeys[key] = pubKey
	}
}

func HandleJWT(next http.Handler) http.Handler {
	var jwtParser = new(jwt.Parser)
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

		if len(jwtPublicKeys) == 0 {
			if !config.IsDevEnvironment() || tokenStr == "" {
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

		var token *jwt.Token
		var err error

		for _, publicKey := range jwtPublicKeys {
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

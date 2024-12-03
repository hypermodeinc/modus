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
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/envfiles"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/jwk"
)

type jwtClaimsKey string

const jwtClaims jwtClaimsKey = "claims"

func Init(ctx context.Context) {
	globalAuthKeys = newAuthKeys()
	go globalAuthKeys.worker(ctx)
	envfiles.RegisterEnvFilesLoadedCallback(initKeys)
	initKeys(ctx)
}

func initKeys(ctx context.Context) {
	publicPemKeysJson := os.Getenv("MODUS_PEMS")
	jwksEndpointsJson := os.Getenv("MODUS_JWKS_ENDPOINTS")
	if publicPemKeysJson == "" && jwksEndpointsJson == "" {
		return
	}

	if publicPemKeysJson != "" {
		keys, err := publicPemKeysJsonToKeys(publicPemKeysJson)
		if err != nil {
			if app.Config().IsDevEnvironment() {
				logger.Fatal(ctx).Err(err).Msg("Auth PEM public keys deserializing error")
			}
			logger.Error(ctx).Err(err).Msg("Auth PEM public keys deserializing error")
			return
		}
		globalAuthKeys.setPemPublicKeys(keys)
	}
	if jwksEndpointsJson != "" {
		keys, err := jwksEndpointsJsonToKeys(ctx, jwksEndpointsJson)
		if err != nil {
			if app.Config().IsDevEnvironment() {
				logger.Fatal(ctx).Err(err).Msg("Auth JWKS public keys deserializing error")
			}
			logger.Error(ctx).Err(err).Msg("Auth JWKS public keys deserializing error")
			return
		}
		globalAuthKeys.setJwksPublicKeys(keys)
	}
}

func Shutdown() {
	close(globalAuthKeys.quit)
	<-globalAuthKeys.done
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
				http.Error(w, "Invalid JWT token format, Bearer required", http.StatusBadRequest)
				return
			}
		}

		if len(globalAuthKeys.getPemPublicKeys()) == 0 && len(globalAuthKeys.getJwksPublicKeys()) == 0 {
			if app.Config().IsDevEnvironment() {
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

			} else {
				next.ServeHTTP(w, r)
				return
			}
		}

		if tokenStr == "" {
			logger.Error(ctx).Msg("JWT token not found")
			http.Error(w, "Access Denied", http.StatusUnauthorized)
			return
		}

		var token *jwt.Token
		var err error
		var found bool

		for _, pemPublicKey := range globalAuthKeys.getPemPublicKeys() {
			token, err = jwtParser.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return pemPublicKey, nil
			})
			if err == nil {
				if utils.DebugModeEnabled() {
					logger.Debug(ctx).Msg("JWT token parsed successfully")
				}
				found = true
				break
			}
		}
		if !found {
			for _, jwksPublicKey := range globalAuthKeys.getJwksPublicKeys() {
				token, err = jwtParser.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
					return jwksPublicKey, nil
				})
				if err == nil {
					if utils.DebugModeEnabled() {
						logger.Debug(ctx).Msg("JWT token parsed successfully")
					}
					found = true
					break
				}
			}
		}
		if !found {
			logger.Error(ctx).Err(err).Msg("JWT parse error")
			http.Error(w, "Access Denied", http.StatusUnauthorized)
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

func publicPemKeysJsonToKeys(publicPemKeysJson string) (map[string]any, error) {
	var publicKeyStrings map[string]string
	if err := json.Unmarshal([]byte(publicPemKeysJson), &publicKeyStrings); err != nil {
		return nil, err
	}
	keys := make(map[string]any)
	for key, value := range publicKeyStrings {
		block, _ := pem.Decode([]byte(value))
		if block == nil {
			return nil, errors.New("Invalid PEM block for key: " + key)
		}

		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, err
		}
		keys[key] = pubKey
	}
	return keys, nil
}

func jwksEndpointsJsonToKeys(ctx context.Context, jwksEndpointsJson string) (map[string]any, error) {
	var jwksEndpoints map[string]string
	if err := json.Unmarshal([]byte(jwksEndpointsJson), &jwksEndpoints); err != nil {
		return nil, err
	}
	keys := make(map[string]any)
	for key, value := range jwksEndpoints {
		jwks, err := jwk.Fetch(ctx, value)
		if err != nil {
			return nil, err
		}

		jwkKey, exists := jwks.Get(0)
		if !exists {
			return nil, errors.New("No keys found in JWKS for key: " + key)
		}

		var rawKey any
		err = jwkKey.Raw(&rawKey)
		if err != nil {
			return nil, err
		}

		// Marshal the raw key into DER-encoded PKIX format
		derBytes, err := x509.MarshalPKIXPublicKey(rawKey)
		if err != nil {
			return nil, err
		}

		pubKey, err := x509.ParsePKIXPublicKey(derBytes)
		if err != nil {
			return nil, err
		}
		keys[key] = pubKey
	}
	return keys, nil
}

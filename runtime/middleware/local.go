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
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ParseJWTUnverified parses a JWT without verifying its signature and returns the claims.
func ParseJWTUnverified(tokenString string) (jwt.MapClaims, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("failed to extract claims from JWT")
}

// CheckJWTExpiration checks if the JWT has expired based on the 'exp' claim.
func CheckJWTExpiration(claims jwt.MapClaims) (bool, error) {
	exp, ok := claims["exp"].(float64)
	if !ok {
		return false, fmt.Errorf("exp claim is missing or not a number")
	}

	expirationTime := time.Unix(int64(exp), 0)
	return time.Now().After(expirationTime), nil
}

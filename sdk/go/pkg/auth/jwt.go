/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package auth

import (
	"errors"
	"os"

	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
)

func GetJWTClaims[T any]() (T, error) {
	var claims T
	claimsStr := os.Getenv("CLAIMS")
	if claimsStr == "" {
		return claims, errors.New("JWT claims not found")
	}
	err := utils.JsonDeserialize([]byte(claimsStr), &claims)
	if err != nil {
		return claims, err
	}

	return claims, nil
}

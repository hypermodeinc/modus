/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

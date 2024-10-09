/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/auth"
)

// This is a sample struct that represents the claims in a JWT token.
type ExampleClaims struct {
	Exp    int64  `json:"exp"`
	Iat    int64  `json:"iat"`
	Iss    string `json:"iss"`
	Jti    string `json:"jti"`
	Nbf    int64  `json:"nbf"`
	Sub    string `json:"sub"`
	UserId string `json:"user-id"`
}

// GetJWTClaims is a function that returns the claims in a JWT token, and parses them into the Claims struct.
func GetJWTClaims() (*ExampleClaims, error) {
	return auth.GetJWTClaims[*ExampleClaims]()
}

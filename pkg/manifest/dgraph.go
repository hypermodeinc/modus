/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

const (
	HostTypeDgraph string = "dgraph"
)

type DgraphHostInfo struct {
	Name       string `json:"-"`
	Type       string `json:"type"`
	GrpcTarget string `json:"grpcTarget"`
	Key        string `json:"key"`
}

func (p DgraphHostInfo) HostName() string {
	return p.Name
}

func (DgraphHostInfo) HostType() string {
	return HostTypeDgraph
}

func (h DgraphHostInfo) GetVariables() []string {
	return extractVariables(h.Key)
}

func (h DgraphHostInfo) Hash() string {
	// Concatenate the attributes into a single string
	data := fmt.Sprintf("%v|%v|%v", h.Name, h.Type, h.GrpcTarget)

	// Compute the SHA-256 hash
	hash := sha256.Sum256([]byte(data))

	// Convert the hash to a hexadecimal string
	hashStr := hex.EncodeToString(hash[:])

	return hashStr
}

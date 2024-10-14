/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

const ConnectionTypeDgraph string = "dgraph"

type DgraphConnectionInfo struct {
	Name       string `json:"-"`
	Type       string `json:"type"`
	GrpcTarget string `json:"grpcTarget"`
	Key        string `json:"key"`
}

func (info DgraphConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info DgraphConnectionInfo) ConnectionType() string {
	return info.Type
}

func (info DgraphConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.GrpcTarget)
}

func (info DgraphConnectionInfo) Variables() []string {
	return extractVariables(info.Key)
}

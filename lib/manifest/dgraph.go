/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

const ConnectionTypeDgraph ConnectionType = "dgraph"

type DgraphConnectionInfo struct {
	Name       string         `json:"-"`
	Type       ConnectionType `json:"type"`
	ConnStr    string         `json:"connString"`
	GrpcTarget string         `json:"grpcTarget"`
	Key        string         `json:"key"`
}

func (info DgraphConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info DgraphConnectionInfo) ConnectionType() ConnectionType {
	return info.Type
}

func (info DgraphConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.GrpcTarget)
}

func (info DgraphConnectionInfo) Variables() []string {
	return extractVariables(info.Key)
}

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

const ConnectionTypeNeo4j ConnectionType = "neo4j"

type Neo4jConnectionInfo struct {
	Name     string         `json:"-"`
	Type     ConnectionType `json:"type"`
	DbUri    string         `json:"dbUri"`
	Username string         `json:"username"`
	Password string         `json:"password"`
}

func (info Neo4jConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info Neo4jConnectionInfo) ConnectionType() ConnectionType {
	return info.Type
}

func (info Neo4jConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.DbUri)
}

func (info Neo4jConnectionInfo) Variables() []string {
	return append(extractVariables(info.Username), extractVariables(info.Password)...)
}

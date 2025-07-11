/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

const ConnectionTypeMysql ConnectionType = "mysql"

type MysqlConnectionInfo struct {
	Name    string         `json:"-"`
	Type    ConnectionType `json:"type"`
	ConnStr string         `json:"connString"`
}

func (info MysqlConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info MysqlConnectionInfo) ConnectionType() ConnectionType {
	return info.Type
}

func (info MysqlConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.ConnStr)
}

func (info MysqlConnectionInfo) Variables() []string {
	return extractVariables(info.ConnStr)
}

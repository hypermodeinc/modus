/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

const ConnectionTypePostgresql ConnectionType = "postgresql"

type PostgresqlConnectionInfo struct {
	Name    string         `json:"-"`
	Type    ConnectionType `json:"type"`
	ConnStr string         `json:"connString"`
}

func (info PostgresqlConnectionInfo) ConnectionName() string {
	return info.Name
}

func (info PostgresqlConnectionInfo) ConnectionType() ConnectionType {
	return info.Type
}

func (info PostgresqlConnectionInfo) Hash() string {
	return computeHash(info.Name, info.Type, info.ConnStr)
}

func (info PostgresqlConnectionInfo) Variables() []string {
	return extractVariables(info.ConnStr)
}

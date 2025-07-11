/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import "context"

type dataSource interface {
	Shutdown()
	query(ctx context.Context, stmt string, params []any, execOnly bool) (*dbResponse, error)
}

type dbResponse struct {
	Error        *string
	Result       any
	RowsAffected uint32
	LastInsertID uint64
}

type HostQueryResponse struct {
	Error        *string
	ResultJson   *string
	RowsAffected uint32
	LastInsertID uint64
}

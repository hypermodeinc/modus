/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

type dbResponse struct {
	Error        *string
	Result       any
	RowsAffected uint32
}

type hostQueryResponse struct {
	Error        *string
	ResultJson   *string
	RowsAffected uint32
}

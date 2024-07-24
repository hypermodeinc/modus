/*
 * Copyright 2024 Hypermode, Inc.
 */

package sqlclient

import (
	"hmruntime/plugins"
)

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

func (r *hostQueryResponse) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "HostQueryResponse",
		Path: "~lib/@hypermode/functions-as/assembly/database/HostQueryResponse",
	}
}

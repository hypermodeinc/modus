/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

type dbResponse struct {
	Error        *string
	Result       any
	RowsAffected uint32
}

type HostQueryResponse struct {
	Error        *string
	ResultJson   *string
	RowsAffected uint32
}

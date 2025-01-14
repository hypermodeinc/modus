/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sqlclient

import "github.com/puzpuzpuz/xsync/v3"

var dsr = newDSRegistry()

type dsRegistry struct {
	cache *xsync.MapOf[string, *postgresqlDS]
}

func newDSRegistry() *dsRegistry {
	return &dsRegistry{
		cache: xsync.NewMapOf[string, *postgresqlDS](),
	}
}

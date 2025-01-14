/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package neo4jclient

type EagerResult struct {
	Keys    []string
	Records []*Record
	// Summary ResultSummary
}

type Record struct {
	Values []string
	Keys   []string
}

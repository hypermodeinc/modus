/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package dgraphclient

type Request struct {
	Query     Query
	Mutations []Mutation
}

type Query struct {
	Query     string
	Variables map[string]string
}

type Mutation struct {
	SetJson   string
	DelJson   string
	SetNquads string
	DelNquads string
	Condition string
}

type Response struct {
	Json string
	Uids map[string]string
}

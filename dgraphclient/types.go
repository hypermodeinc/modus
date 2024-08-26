/*
 * Copyright 2024 Hypermode, Inc.
 */

package dgraphclient

type Request struct {
	Query     Query
	Mutations []Mutation
	CommitNow bool
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

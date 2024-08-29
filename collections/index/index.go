/*
 * Copyright 2024 Hypermode, Inc.
 */

package index

// SearchFilter defines a predicate function that we will use to determine
// whether or not a given vector is "interesting". When used in the context
// of VectorIndex.Search, a true result means that we want to keep the result
// in the returned list, and a false result implies we should skip.
type SearchFilter func(query, resultVal []float32, resultUID string) bool

// AcceptAll implements SearchFilter by way of accepting all results.
func AcceptAll(_, _ []float32, _ string) bool { return true }

// AcceptNone implements SearchFilter by way of rejecting all results.
func AcceptNone(_, _ []float32, _ string) bool { return false }

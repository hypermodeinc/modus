package index

import (
	"context"
)

// VectorSource defines the process of getting a sequence of Vectors for
// use in defining a new VectorIndex. Each returned vector is assumed
// to be associated with a unique identifier specified by a string.
// Iteration proceeds until Next() or NextChunk() returns an indicator
// that there are no more values.
type VectorSource interface {
	// Next will return a vector (specified by []Float) to be indexed,
	// a unique identifier (string), an indicator of whether or not the
	// value is valid, and possibly an error if there are troubles reaching
	// an underlying persistent storage.
	// Note that if the bool return value is true, you should accept the
	// returned value, but a false value indicates that the returned float
	// and string should be ignored.
	// If returned error is not nil, you should expect that the bool returned
	// value is false, and that the sourcing process is not recoverable
	// (i.e., if you want to retry, you need to start from the beginning).
	Next() ([]float64, string, bool, error)

	// NextChunk behaves like Next(), but returns multiple floats at the same
	// time. The choice of how many to return is based on the VectorSource
	// implementation.
	NextChunk() ([][]float64, []string, bool, error)
}

type emptyVectorSource struct{}

// EmptyVectorSource returns an implementation of the VectorSource interface
// that supplies no elements. This is useful when creating an initially empty
// vector index.
func EmptyVectorSource() VectorSource {
	return &emptyVectorSource{}
}

// Next() is part of VectorSource interface implementation.
func (vs *emptyVectorSource) Next() ([]float64, string, bool, error) {
	return nil, "", false, nil
}

// NextChunk() is part of VectorSource interface implementation.
func (vs *emptyVectorSource) NextChunk() ([][]float64, []string, bool, error) {
	return nil, nil, false, nil
}

// SearchFilter defines a predicate function that we will use to determine
// whether or not a given vector is "interesting". When used in the context
// of VectorIndex.Search, a true result means that we want to keep the result
// in the returned list, and a false result implies we should skip.
type SearchFilter func(query, resultVal []float64, resultUID string) bool

// AcceptAll implements SearchFilter by way of accepting all results.
func AcceptAll(_, _ []float64, _ string) bool { return true }

// AcceptNone implements SearchFilter by way of rejecting all results.
func AcceptNone(_, _ []float64, _ string) bool { return false }

// A Txn is an interface representation of a persistent storage transaction,
// where multiple operations are performed on a database
type Txn interface {
	// StartTs gets the exact time that the transaction started, returned in string format
	StartTs() string
	// Get uses a []byte key to return the Value corresponding to the key
	Get(key []byte) (rval Value, rerr error)
	// GetWithLockHeld uses a []byte key to return the Value corresponding to the key with a mutex lock held
	GetWithLockHeld(key []byte) (rval Value, rerr error)
	Find(prefix []byte, filter func(val []byte) bool) (string, error)
	// Adds a mutation operation on a index.Txn interface, where the mutation
	// is represented in the form of an index.DirectedEdge
	AddMutation(ctx context.Context, key []byte, t *KeyValue) error
	// Same as AddMutation but with a mutex lock held
	AddMutationWithLockHeld(ctx context.Context, key []byte, t *KeyValue) error
	// mutex lock
	LockKey(key []byte)
	// mutex unlock
	UnlockKey(key []byte)
}

// Local cache is an interface representation of the local cache of a persistent storage system
type LocalCache interface {
	// Get uses a []byte key to return the Value corresponding to the key
	Get(key []byte) (rval Value, rerr error)
	// GetWithLockHeld uses a []byte key to return the Value corresponding to the key with a mutex lock held
	GetWithLockHeld(key []byte) (rval Value, rerr error)
	Find(prefix []byte, filter func(val []byte) bool) (string, error)
}

// Value is an interface representation of the value of a persistent storage system
type Value interface{}

// CacheType is an interface representation of the cache of a persistent storage system
type CacheType interface {
	Get(key []byte) (rval Value, rerr error)
	Ts() string
	Find(prefix []byte, filter func(val []byte) bool) (string, error)
}

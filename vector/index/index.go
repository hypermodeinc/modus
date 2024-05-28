package index

import (
	"context"

	c "hmruntime/vector/constraints"

	"hmruntime/vector/options"
)

// VectorSource defines the process of getting a sequence of Vectors for
// use in defining a new VectorIndex. Each returned vector is assumed
// to be associated with a unique identifier specified by a uint64.
// Iteration proceeds until Next() or NextChunk() returns an indicator
// that there are no more values.
type VectorSource[T c.Float] interface {
	// Next will return a vector (specified by []Float) to be indexed,
	// a unique identifier (uint64), an indicator of whether or not the
	// value is valid, and possibly an error if there are troubles reaching
	// an underlying persistent storage.
	// Note that if the bool return value is true, you should accept the
	// returned value, but a false value indicates that the returned float
	// and uint64 should be ignored.
	// If returned error is not nil, you should expect that the bool returned
	// value is false, and that the sourcing process is not recoverable
	// (i.e., if you want to retry, you need to start from the beginning).
	Next() ([]T, uint64, bool, error)

	// NextChunk behaves like Next(), but returns multiple floats at the same
	// time. The choice of how many to return is based on the VectorSource
	// implementation.
	NextChunk() ([][]T, []uint64, bool, error)
}

type emptyVectorSource[T c.Float] struct{}

// EmptyVectorSource returns an implementation of the VectorSource interface
// that supplies no elements. This is useful when creating an initially empty
// vector index.
func EmptyVectorSource[T c.Float]() VectorSource[T] {
	return &emptyVectorSource[T]{}
}

// Next() is part of VectorSource interface implementation.
func (vs *emptyVectorSource[T]) Next() ([]T, uint64, bool, error) {
	return nil, 0, false, nil
}

// NextChunk() is part of VectorSource interface implementation.
func (vs *emptyVectorSource[T]) NextChunk() ([][]T, []uint64, bool, error) {
	return nil, nil, false, nil
}

// SearchFilter defines a predicate function that we will use to determine
// whether or not a given vector is "interesting". When used in the context
// of VectorIndex.Search, a true result means that we want to keep the result
// in the returned list, and a false result implies we should skip.
type SearchFilter[T c.Float] func(query, resultVal []T, resultUID uint64) bool

// AcceptAll implements SearchFilter by way of accepting all results.
func AcceptAll[T c.Float](_, _ []T, _ uint64) bool { return true }

// AcceptNone implements SearchFilter by way of rejecting all results.
func AcceptNone[T c.Float](_, _ []T, _ uint64) bool { return false }

// OptionalIndexSupport defines abilities that might not be universally
// supported by all VectorIndex types. A VectorIndex will technically
// define the functions required by OptionalIndexSupport, but may do so
// by way of simply returning an errors.ErrUnsupported result.
type OptionalIndexSupport[T c.Float] interface {
	// SearchWithPath(ctx, c, query, maxResults, filter) is similar to
	// Search(ctx, c, query, maxResults, filter), but returns an extended
	// set of content in the search results.
	// The full contents returned are indicated by the SearchPathResult.
	// See the description there for more info.
	SearchWithPath(
		ctx context.Context,
		c CacheType,
		query []T,
		maxResults int,
		filter SearchFilter[T]) (*SearchPathResult, error)
}

// A VectorIndex can be used to Search for vectors and add vectors to an index.
type VectorIndex[T c.Float] interface {
	OptionalIndexSupport[T]

	// Specifies the set of allowed options and a corresponding means to
	// parse a string version of those options.
	AllowedOptions() options.AllowedOptions

	ApplyOptions(o options.Options) error

	// Search will find the uids for a given set of vectors based on the
	// input query, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	Search(ctx context.Context, c CacheType, query []T,
		maxResults int,
		filter SearchFilter[T]) ([]uint64, error)

	// SearchWithUid will find the uids for a given set of vectors based on the
	// input queryUid, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	SearchWithUid(ctx context.Context, c CacheType, queryUid uint64,
		maxResults int,
		filter SearchFilter[T]) ([]uint64, error)

	// Insert will add a vector and uuid into the existing VectorIndex. If
	// uuid already exists, it should throw an error to not insert duplicate uuids
	Insert(ctx context.Context, c CacheType, uuid uint64, vec []T) ([]*KeyValue, error)

	// InsertDataNode will add a data node to the existing VectorIndex
	InsertDataNode(ctx context.Context, c CacheType, uuid uint64, data string) error

	// OptionalFeatures() returns a collection of optional features that
	// may be supported by this class. By default, every implementation
	// should attempt to faithfully define each of the above functions,
	// but an index type that cannot support some optional capability
	// can be allowed to "implement" the feature by returning
	// an errors.NotSupported result.
	// This result should almost always match the result returned by
	// the factory used to generate this.
	// (This condition may be relaxed in the future, but for now,
	// it is safer to make this assumption).
}

// A Txn is an interface representation of a persistent storage transaction,
// where multiple operations are performed on a database
type Txn interface {
	// StartTs gets the exact time that the transaction started, returned in uint64 format
	StartTs() uint64
	// Get uses a []byte key to return the Value corresponding to the key
	Get(key []byte) (rval Value, rerr error)
	// GetWithLockHeld uses a []byte key to return the Value corresponding to the key with a mutex lock held
	GetWithLockHeld(key []byte) (rval Value, rerr error)
	Find(prefix []byte, filter func(val []byte) bool) (uint64, error)
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
	Find(prefix []byte, filter func(val []byte) bool) (uint64, error)
}

// Value is an interface representation of the value of a persistent storage system
type Value interface{}

// CacheType is an interface representation of the cache of a persistent storage system
type CacheType interface {
	Get(key []byte) (rval Value, rerr error)
	Ts() uint64
	Find(prefix []byte, filter func(val []byte) bool) (uint64, error)
}

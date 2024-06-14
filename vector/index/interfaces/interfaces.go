package interfaces

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/vector/in_mem/sequential"
	"hmruntime/vector/index"
	"hmruntime/vector/utils"
)

var (
	ErrInvalidVectorIndexType = fmt.Errorf("invalid vector index type")
)

const (
	SequentialManifestType = "sequential"
	HnswManifestType       = "hnsw"
)

type VectorIndexWrapper struct {
	Type string `json:"Type"`
	VectorIndex
}

type UnmarshallSequentialVectorIndex struct {
	VectorIndex sequential.SequentialVectorIndex `json:"VectorIndex"`
}

func (v *VectorIndexWrapper) UnmarshalJSON(data []byte) error {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return fmt.Errorf("error unmarshalling JSON data: %w", err)
	}

	rawType, ok := m["Type"]
	if !ok {
		return fmt.Errorf("type field not found in JSON data")
	}

	var t string
	if err := json.Unmarshal(rawType, &t); err != nil {
		return fmt.Errorf("error unmarshalling type field: %w", err)
	}

	switch t {
	case sequential.SequentialVectorIndexType:
		rawVectorIndex, ok := m["VectorIndex"]
		if !ok {
			return fmt.Errorf("VectorIndex field not found in JSON data")
		}

		var index sequential.SequentialVectorIndex
		if err := json.Unmarshal(rawVectorIndex, &index); err != nil {
			return fmt.Errorf("error unmarshalling VectorIndex field: %w", err)
		}
		v.Type = t
		v.VectorIndex = &index
	default:
		return fmt.Errorf("invalid vector index type: %s", t)
	}

	return nil
}

type TextIndex interface {
	// GetVectorIndexMap returns the map of searchMethod to VectorIndex
	GetVectorIndexMap() map[string]*VectorIndexWrapper

	// GetVectorIndex returns the VectorIndex for a given searchMethod
	GetVectorIndex(searchMethod string) (*VectorIndexWrapper, error)

	// SetVectorIndex sets the VectorIndex for a given searchMethod
	SetVectorIndex(searchMethod string, index *VectorIndexWrapper) error

	// DeleteVectorIndex deletes the VectorIndex for a given searchMethod
	DeleteVectorIndex(searchMethod string) error

	// InsertText will add a text and uuid into the existing VectorIndex
	InsertText(ctx context.Context, uuid string, text string) ([]*index.KeyValue, error)

	// DeleteText will remove a text and uuid from the existing VectorIndex
	DeleteText(ctx context.Context, uuid string) error

	// GetText will return the text for a given uuid
	GetText(ctx context.Context, uuid string) (string, error)

	// GetTextMap returns the map of uuid to text
	GetTextMap() map[string]string
}

// A VectorIndex can be used to Search for vectors and add vectors to an index.
type VectorIndex interface {

	// Search will find the uids for a given set of vectors based on the
	// input query, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	Search(ctx context.Context, query []float64,
		maxResults int,
		filter index.SearchFilter) (utils.MinTupleHeap, error)

	// SearchWithUid will find the uids for a given set of vectors based on the
	// input queryUid, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	SearchWithUid(ctx context.Context, queryUid string,
		maxResults int,
		filter index.SearchFilter) (utils.MinTupleHeap, error)

	// Insert will add a vector and uuid into the existing VectorIndex. If
	// uuid already exists, it should throw an error to not insert duplicate uuids
	InsertVector(ctx context.Context, uuid string, vec []float64) ([]*index.KeyValue, error)

	// Delete will remove a vector and uuid from the existing VectorIndex. If
	// uuid does not exist, it should throw an error to not delete non-existent uuids
	DeleteVector(ctx context.Context, uuid string) error
}

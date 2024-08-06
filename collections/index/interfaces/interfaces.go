package interfaces

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/collections/in_mem/sequential"
	"hmruntime/collections/index"
	"hmruntime/collections/utils"
)

var (
	ErrInvalidVectorIndexType = fmt.Errorf("invalid vector index type")
	ErrInvalidVectorIndexName = fmt.Errorf("invalid vector index name")
)

const (
	SequentialManifestType = "sequential"
	HnswManifestType       = "hnsw"
)

type VectorIndexWrapper struct {
	Type string `json:"Type"`
	VectorIndex
}

type UnmarshalSequentialVectorIndex struct {
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

type Collection interface {
	GetCollectionName() string

	// GetVectorIndexMap returns the map of searchMethod to VectorIndex
	GetVectorIndexMap() map[string]*VectorIndexWrapper

	// GetVectorIndex returns the VectorIndex for a given searchMethod
	GetVectorIndex(ctx context.Context, searchMethod string) (*VectorIndexWrapper, error)

	// SetVectorIndex sets the VectorIndex for a given searchMethod
	SetVectorIndex(ctx context.Context, searchMethod string, index *VectorIndexWrapper) error

	// DeleteVectorIndex deletes the VectorIndex for a given searchMethod
	DeleteVectorIndex(ctx context.Context, searchMethod string) error

	// InsertTexts will add texts and keys into the existing VectorIndex
	InsertTexts(ctx context.Context, keys []string, texts []string) error

	// InsertText will add a text and key into the existing VectorIndex
	InsertText(ctx context.Context, key string, text string) error

	InsertTextsToMemory(ctx context.Context, ids []int64, keys []string, texts []string) error

	InsertTextToMemory(ctx context.Context, id int64, key string, text string) error

	// DeleteText will remove a text and key from the existing VectorIndex
	DeleteText(ctx context.Context, key string) error

	// GetText will return the text for a given key
	GetText(ctx context.Context, key string) (string, error)

	// GetTextMap returns the map of key to text
	GetTextMap(ctx context.Context) (map[string]string, error)

	// GetExternalId returns the external id for a given key
	GetExternalId(ctx context.Context, key string) (int64, error)

	GetCheckpointId(ctx context.Context) (int64, error)
}

// A VectorIndex can be used to Search for vectors and add vectors to an index.
type VectorIndex interface {
	GetSearchMethodName() string

	SetEmbedderName(embedderName string) error

	GetEmbedderName() string

	// Search will find the keys for a given set of vectors based on the
	// input query, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	Search(ctx context.Context, query []float32,
		maxResults int,
		filter index.SearchFilter) (utils.MaxTupleHeap, error)

	// SearchWithKey will find the keys for a given set of vectors based on the
	// input queryKey, limiting to the specified maximum number of results.
	// The filter parameter indicates that we might discard certain parameters
	// based on some input criteria. The maxResults count is counted *after*
	// being filtered. In other words, we only count those results that had not
	// been filtered out.
	SearchWithKey(ctx context.Context, queryKey string,
		maxResults int,
		filter index.SearchFilter) (utils.MaxTupleHeap, error)

	// Insert Vectors will add vectors and keys into the existing VectorIndex
	InsertVectors(ctx context.Context, textIds []int64, vecs [][]float32) error

	// Insert will add a vector and key into the existing VectorIndex. If
	// key already exists, it should throw an error to not insert duplicate keys
	InsertVector(ctx context.Context, textId int64, vec []float32) error

	// InsertVectorToMemory will add a vector and key into the existing VectorIndex. If
	// key already exists, it should throw an error to not insert duplicate keys
	InsertVectorToMemory(ctx context.Context, textId, vectorId int64, key string, vec []float32) error

	// Delete will remove a vector and key from the existing VectorIndex. If
	// key does not exist, it should throw an error to not delete non-existent keys
	DeleteVector(ctx context.Context, textId int64, key string) error

	// GetVector will return the vector for a given key
	GetVector(ctx context.Context, key string) ([]float32, error)

	GetCheckpointId(ctx context.Context) (int64, error)

	GetLastIndexedTextId(ctx context.Context) (int64, error)
}

package hnsw

import (
	"container/heap"
	"context"
	"fmt"
	"hmruntime/collections/index"
	"hmruntime/collections/utils"
	"hmruntime/db"
	"hmruntime/logger"
	"time"

	"github.com/hypermodeAI/hnsw"
)

const (
	HnswVectorIndexType = "HnswVectorIndex"
)

type HnswVectorIndex struct {
	// mu                sync.RWMutex
	searchMethodName  string
	embedderName      string
	lastInsertedID    int64
	lastIndexedTextID int64
	HnswIndex         *hnsw.Graph[string]
}

func NewHnswVectorIndex(collection, searchMethod, embedder string) *HnswVectorIndex {
	return &HnswVectorIndex{
		searchMethodName: searchMethod,
		embedderName:     embedder,
		HnswIndex:        hnsw.NewGraph[string](),
	}
}

func (ims *HnswVectorIndex) GetSearchMethodName() string {
	return ims.searchMethodName
}

func (ims *HnswVectorIndex) SetEmbedderName(embedderName string) error {
	// ims.mu.Lock()
	// defer ims.mu.Unlock()
	ims.embedderName = embedderName
	return nil
}

func (ims *HnswVectorIndex) GetEmbedderName() string {
	// ims.mu.RLock()
	// defer ims.mu.RUnlock()
	return ims.embedderName
}

func (ims *HnswVectorIndex) GetVectorNodesMap() map[string][]float32 {
	return map[string][]float32{}
}

func (ims *HnswVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	// calculate cosine similarity and return top maxResults results
	start := time.Now()
	logger.Info(ctx).Msg(start.String())
	// ims.mu.RLock()
	// defer ims.mu.RUnlock()
	if ims.HnswIndex == nil {
		return nil, fmt.Errorf("vector index is not initialized")
	}
	fmt.Println(ims.HnswIndex.Len())
	neighbors, err := ims.HnswIndex.Search(query, maxResults)
	if err != nil {
		return nil, err
	}
	keys, distances := make([]string, len(neighbors)), make([]float64, len(neighbors))
	for i, neighbor := range neighbors {
		keys[i] = string(neighbor.Key)
		distances[i] = float64(neighbor.Distance)
	}
	var results utils.MinTupleHeap
	heap.Init(&results)

	for i, key := range keys {
		heap.Push(&results, utils.InitHeapElement(distances[i], key, false))
	}

	// Return top maxResults results
	var finalResults utils.MinTupleHeap
	for results.Len() > 0 {
		finalResults = append(finalResults, heap.Pop(&results).(utils.MinHeapElement))
	}
	// Reverse the finalResults to get the highest similarity first
	for i, j := 0, len(finalResults)-1; i < j; i, j = i+1, j-1 {
		finalResults[i], finalResults[j] = finalResults[j], finalResults[i]
	}

	logger.Info(ctx).Msg(fmt.Sprintf("Search took %v", time.Since(start)))

	return finalResults, nil
}

func (ims *HnswVectorIndex) SearchWithKey(ctx context.Context, queryKey string, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	// ims.mu.RLock()
	query, found := ims.HnswIndex.Lookup(queryKey)
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	// ims.mu.RUnlock()
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, query, maxResults, filter)
}

func (ims *HnswVectorIndex) InsertVectors(ctx context.Context, textIds []int64, vecs [][]float32) error {
	// ims.mu.Lock()
	// defer ims.mu.Unlock()
	if len(textIds) != len(vecs) {
		return fmt.Errorf("textIds and vecs must have the same length")
	}
	vectorIds, keys, err := db.WriteCollectionVectors(ctx, ims.searchMethodName, textIds, vecs)
	if err != nil {
		return err
	}

	return ims.InsertVectorsToMemory(ctx, textIds, vectorIds, keys, vecs)
}

func (ims *HnswVectorIndex) InsertVector(ctx context.Context, textId int64, vec []float32) error {

	// Write vector to database, this textId is now the last inserted textId
	vectorId, key, err := db.WriteCollectionVector(ctx, ims.searchMethodName, textId, vec)
	if err != nil {
		return err
	}

	return ims.InsertVectorToMemory(ctx, textId, vectorId, key, vec)

}

func (ims *HnswVectorIndex) InsertVectorsToMemory(ctx context.Context, textIds []int64, vectorIds []int64, keys []string, vecs [][]float32) error {
	nodes, err := hnsw.MakeNodes(keys, vecs)
	if err != nil {
		return err
	}
	err = ims.HnswIndex.Add(nodes...)
	if err != nil {
		return err
	}
	ims.lastInsertedID = vectorIds[len(vectorIds)-1]
	ims.lastIndexedTextID = textIds[len(textIds)-1]
	return nil
}

func (ims *HnswVectorIndex) InsertVectorToMemory(ctx context.Context, textId, vectorId int64, key string, vec []float32) error {
	// ims.mu.Lock()
	// defer ims.mu.Unlock()
	err := ims.HnswIndex.Add(hnsw.MakeNode(key, vec))
	if err != nil {
		return err
	}
	ims.lastInsertedID = vectorId
	ims.lastIndexedTextID = textId
	return nil
}

func (ims *HnswVectorIndex) DeleteVector(ctx context.Context, textId int64, key string) error {
	// ims.mu.Lock()
	// defer ims.mu.Unlock()
	err := db.DeleteCollectionVector(ctx, ims.searchMethodName, textId)
	if err != nil {
		return err
	}
	ims.HnswIndex.Delete(key)
	return nil
}

func (ims *HnswVectorIndex) GetVector(ctx context.Context, key string) ([]float32, error) {
	// ims.mu.RLock()
	// defer ims.mu.RUnlock()
	vec, _ := ims.HnswIndex.Lookup(key)
	return vec, nil
}

func (ims *HnswVectorIndex) GetCheckpointId(ctx context.Context) (int64, error) {
	// ims.mu.RLock()
	// defer ims.mu.RUnlock()
	return ims.lastInsertedID, nil
}

func (ims *HnswVectorIndex) GetLastIndexedTextId(ctx context.Context) (int64, error) {
	// ims.mu.RLock()
	// defer ims.mu.RUnlock()
	return ims.lastIndexedTextID, nil
}

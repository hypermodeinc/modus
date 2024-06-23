package sequential

import (
	"container/heap"
	"context"
	"hmruntime/collections/index"
	"hmruntime/collections/utils"
	"hmruntime/db"
	"sync"
)

const (
	SequentialVectorIndexType = "SequentialVectorIndex"
)

type SequentialVectorIndex struct {
	mu                sync.RWMutex
	searchMethodName  string
	embedderName      string
	lastInsertedID    int64
	lastIndexedTextID int64
	VectorMap         map[string][]float32 // key: vector
}

func NewSequentialVectorIndex(collection, searchMethod, embedder string) *SequentialVectorIndex {
	return &SequentialVectorIndex{
		searchMethodName: searchMethod,
		embedderName:     embedder,
		VectorMap:        make(map[string][]float32),
	}
}

func (ims *SequentialVectorIndex) GetSearchMethodName() string {
	return ims.searchMethodName
}

func (ims *SequentialVectorIndex) GetEmbedderName() string {
	return ims.embedderName
}

func (ims *SequentialVectorIndex) GetVectorNodesMap() map[string][]float32 {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorMap
}

func (ims *SequentialVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	// calculate cosine similarity and return top maxResults results
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	var results utils.MinTupleHeap
	heap.Init(&results)
	for key, vector := range ims.VectorMap {
		if filter != nil && !filter(query, vector, key) {
			continue
		}
		similarity, err := utils.CosineSimilarity(query, vector)
		if err != nil {
			return nil, err
		}
		if results.Len() < maxResults {
			heap.Push(&results, utils.InitHeapElement(similarity, key, false))
		} else if utils.IsBetterScoreForSimilarity(similarity, results[0].GetValue()) {
			heap.Pop(&results)
			heap.Push(&results, utils.InitHeapElement(similarity, key, false))
		}
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
	return finalResults, nil
}

func (ims *SequentialVectorIndex) SearchWithKey(ctx context.Context, queryKey string, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	ims.mu.RLock()
	query := ims.VectorMap[queryKey]
	ims.mu.RUnlock()
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, query, maxResults, filter)
}

func (ims *SequentialVectorIndex) InsertVector(ctx context.Context, textId int64, vec []float32) error {

	// Write vector to database, this textId is now the last inserted textId
	vectorId, key, err := db.WriteCollectionVector(ctx, ims.searchMethodName, textId, vec)
	if err != nil {
		return err
	}

	return ims.InsertVectorToMemory(ctx, textId, vectorId, key, vec)

}

func (ims *SequentialVectorIndex) InsertVectorToMemory(ctx context.Context, textId, vectorId int64, key string, vec []float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.VectorMap[key] = vec
	ims.lastInsertedID = vectorId
	ims.lastIndexedTextID = textId
	return nil
}

func (ims *SequentialVectorIndex) DeleteVector(ctx context.Context, textId int64, key string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	err := db.DeleteCollectionVector(ctx, ims.searchMethodName, textId)
	if err != nil {
		return err
	}
	delete(ims.VectorMap, key)
	return nil
}

func (ims *SequentialVectorIndex) GetVector(ctx context.Context, key string) ([]float32, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorMap[key], nil
}

func (ims *SequentialVectorIndex) GetCheckpointId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastInsertedID, nil
}

func (ims *SequentialVectorIndex) GetLastIndexedTextId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastIndexedTextID, nil
}

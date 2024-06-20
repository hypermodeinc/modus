package sequential

import (
	"container/heap"
	"context"
	"hmruntime/collections/index"
	"hmruntime/collections/utils"
	"sync"
)

const (
	SequentialVectorIndexType = "SequentialVectorIndex"
)

type SequentialVectorIndex struct {
	// vectorNodes is a map of string to []float32
	mu          sync.RWMutex
	VectorNodes map[string][]float32
}

func NewSequentialVectorIndex() *SequentialVectorIndex {
	return &SequentialVectorIndex{
		VectorNodes: make(map[string][]float32),
	}
}

func (ims *SequentialVectorIndex) GetVectorNodesMap() map[string][]float32 {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorNodes
}

func (ims *SequentialVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	// calculate cosine similarity and return top maxResults results
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	var results utils.MinTupleHeap
	heap.Init(&results)
	for uid, vector := range ims.VectorNodes {
		if filter != nil && !filter(query, vector, uid) {
			continue
		}
		similarity, err := utils.CosineSimilarity(query, vector)
		if err != nil {
			return nil, err
		}
		if results.Len() < maxResults {
			heap.Push(&results, utils.InitHeapElement(similarity, uid, false))
		} else if utils.IsBetterScoreForSimilarity(similarity, results[0].GetValue()) {
			heap.Pop(&results)
			heap.Push(&results, utils.InitHeapElement(similarity, uid, false))
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

func (ims *SequentialVectorIndex) SearchWithUid(ctx context.Context, queryUid string, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	ims.mu.RLock()
	query := ims.VectorNodes[queryUid]
	ims.mu.RUnlock()
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, query, maxResults, filter)
}

func (ims *SequentialVectorIndex) InsertVector(ctx context.Context, uid string, vector []float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.VectorNodes[uid] = vector
	return nil
}

func (ims *SequentialVectorIndex) DeleteVector(ctx context.Context, uid string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.VectorNodes, uid)
	return nil
}

func (ims *SequentialVectorIndex) GetVector(ctx context.Context, uid string) ([]float32, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorNodes[uid], nil
}

package in_mem

import (
	"container/heap"
	"context"
	"encoding/gob"
	"hmruntime/vector/index"
	"hmruntime/vector/options"
	"hmruntime/vector/utils"
	"os"
	"sync"
)

type SequentialVectorIndex struct {
	// vectorNodes is a map of string to []float64
	mu          sync.RWMutex
	vectorNodes map[string][]float64
}

func NewSequentialVectorIndex() *SequentialVectorIndex {
	return &SequentialVectorIndex{
		vectorNodes: make(map[string][]float64),
	}
}

func (ims *SequentialVectorIndex) GetVectorNodesMap() map[string][]float64 {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.vectorNodes
}

func (ims *SequentialVectorIndex) AllowedOptions() options.AllowedOptions {
	return nil
}

func (ims *SequentialVectorIndex) ApplyOptions(o options.Options) error {
	return nil
}

func (ims *SequentialVectorIndex) emptyFinalResultWithError(e error) (
	*index.SearchPathResult, error) {
	return index.NewSearchPathResult(), e
}

func (ims *SequentialVectorIndex) Search(ctx context.Context, c index.CacheType, query []float64, maxResults int, filter index.SearchFilter[float64]) (utils.MinTupleHeap[float64], error) {
	// calculate cosine similarity and return top maxResults results
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	var results utils.MinTupleHeap[float64]
	heap.Init(&results)
	for uid, vector := range ims.vectorNodes {
		similarity, err := utils.CosineSimilarity[float64](query, vector, 64)
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
	var finalResults utils.MinTupleHeap[float64]
	for results.Len() > 0 {
		finalResults = append(finalResults, heap.Pop(&results).(utils.MinHeapElement[float64]))
	}
	// Reverse the finalResults to get the highest similarity first
	for i, j := 0, len(finalResults)-1; i < j; i, j = i+1, j-1 {
		finalResults[i], finalResults[j] = finalResults[j], finalResults[i]
	}
	return finalResults, nil
}

func (ims *SequentialVectorIndex) SearchWithUid(ctx context.Context, c index.CacheType, queryUid string, maxResults int, filter index.SearchFilter[float64]) (utils.MinTupleHeap[float64], error) {
	query := ims.vectorNodes[queryUid]
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, c, query, maxResults, filter)
}

func (ims *SequentialVectorIndex) SearchWithPath(ctx context.Context, c index.CacheType, query []float64, maxResults int, filter index.SearchFilter[float64]) (*index.SearchPathResult, error) {
	return ims.emptyFinalResultWithError(nil)
}

func (ims *SequentialVectorIndex) InsertVector(ctx context.Context, c index.CacheType, uid string, vector []float64) ([]*index.KeyValue, error) {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.vectorNodes[uid] = vector
	return nil, nil
}

func (ims *SequentialVectorIndex) DeleteVector(ctx context.Context, c index.CacheType, uid string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	delete(ims.vectorNodes, uid)
	return nil
}

func (ims *SequentialVectorIndex) WriteToWAL(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(ims.vectorNodes); err != nil {
		return err
	}

	return nil
}

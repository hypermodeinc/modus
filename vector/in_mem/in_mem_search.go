package in_mem

import (
	"container/heap"
	"context"
	"encoding/gob"
	"hmruntime/vector/index"
	"hmruntime/vector/options"
	"hmruntime/vector/utils"
	"os"
)

type InMemBruteForceIndex struct {
	// vectorNodes is a map of uint64 to []float64
	vectorNodes map[uint64][]float64
	dataNodes   map[uint64]string
}

func (ims *InMemBruteForceIndex) AllowedOptions() options.AllowedOptions {
	return nil
}

func (ims *InMemBruteForceIndex) ApplyOptions(o options.Options) error {
	return nil
}

func (ims *InMemBruteForceIndex) emptyFinalResultWithError(e error) (
	*index.SearchPathResult, error) {
	return index.NewSearchPathResult(), e
}

func (ims *InMemBruteForceIndex) Search(ctx context.Context, c index.CacheType, query []float64, maxResults int, filter index.SearchFilter[float64]) ([]uint64, error) {
	// calculate cosine similarity and return top maxResults results
	var results utils.MinTupleHeap[float64]
	heap.Init(&results)
	for uid, vector := range ims.vectorNodes {
		similarity, err := utils.CosineSimilarity[float64](query, vector, 64)
		if err != nil {
			return nil, err
		}
		if results.Len() < maxResults {
			results.Push(utils.InitHeapElement(similarity, uid, false))
		} else if utils.IsBetterScoreForSimilarity(similarity, results[0].GetValue()) {
			results.Pop()
			results.Push(utils.InitHeapElement(similarity, uid, false))
		}
	}

	// Return top maxResults results
	var uids []uint64
	for len(results) > 0 {
		top := heap.Pop(&results).(*utils.MinHeapElement[float64])
		uids = append(uids, top.GetIndex())
	}
	// Reverse the uids to get the highest similarity first
	for i, j := 0, len(uids)-1; i < j; i, j = i+1, j-1 {
		uids[i], uids[j] = uids[j], uids[i]
	}
	return uids, nil
}

func (ims *InMemBruteForceIndex) SearchWithUid(ctx context.Context, c index.CacheType, queryUid uint64, maxResults int, filter index.SearchFilter[float64]) ([]uint64, error) {
	query := ims.vectorNodes[queryUid]
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, c, query, maxResults, filter)
}

func (ims *InMemBruteForceIndex) SearchWithPath(ctx context.Context, c index.CacheType, query []float64, maxResults int, filter index.SearchFilter[float64]) (*index.SearchPathResult, error) {
	return ims.emptyFinalResultWithError(nil)
}

func (ims *InMemBruteForceIndex) Insert(ctx context.Context, c index.CacheType, uid uint64, vector []float64) ([]*index.KeyValue, error) {
	ims.vectorNodes[uid] = vector
	return nil, nil
}

func (ims *InMemBruteForceIndex) InsertDataNode(ctx context.Context, c index.CacheType, uid uint64, data string) error {
	ims.dataNodes[uid] = data
	return nil
}

func (ims *InMemBruteForceIndex) GetDataNode(ctx context.Context, c index.CacheType, uid uint64) string {
	return ims.dataNodes[uid]
}

func (ims *InMemBruteForceIndex) WriteToWAL(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	if err := encoder.Encode(ims.vectorNodes); err != nil {
		return err
	}

	if err := encoder.Encode(ims.dataNodes); err != nil {
		return err
	}

	return nil
}

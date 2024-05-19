package in_mem

import (
	"context"
	"hmruntime/vector/index"
	"hmruntime/vector/options"
	"hmruntime/vector/utils"
	"sort"
)

type InMemBruteForceIndex struct {
	pred string
	// vectorNodes is a map of uint64 to []float64
	vectorNodes map[uint64][]float64
	dataNodes   map[uint64]string
}

func (ims *InMemBruteForceIndex) applyOptions(o options.Options) error {
	return nil
}

func (ims *InMemBruteForceIndex) emptyFinalResultWithError(e error) (
	*index.SearchPathResult, error) {
	return index.NewSearchPathResult(), e
}

func (ims *InMemBruteForceIndex) Search(ctx context.Context, c index.CacheType, query []float64, maxResults int, filter index.SearchFilter[float64]) ([]uint64, error) {
	// calculate cosine similarity and return top maxResults results
	type result struct {
		uid        uint64
		similarity float64
	}
	var results []result
	for uid, vector := range ims.vectorNodes {
		similarity, err := utils.CosineSimilarity[float64](query, vector, 64)
		if err != nil {
			return nil, err
		}
		results = append(results, result{uid, similarity})
	}

	// Sort results by similarity in descending order
	sort.Slice(results, func(i, j int) bool {
		return results[i].similarity > results[j].similarity
	})

	// Return top maxResults results
	var uids []uint64
	for i := 0; i < maxResults && i < len(results); i++ {
		uids = append(uids, results[i].uid)
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

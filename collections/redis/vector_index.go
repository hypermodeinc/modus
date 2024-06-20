package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"hmruntime/collections/utils"
	"sync"
)

const (
	RedisVectorIndexType = "RedisVectorIndex"
)

type RedisVectorIndex struct {
	name            string
	indexType       string
	containsVectors bool
	mu              sync.RWMutex
}

func NewRedisVectorIndex(name string, indexType string) *RedisVectorIndex {
	return &RedisVectorIndex{
		name:            name, //collection:searchMethod
		indexType:       indexType,
		containsVectors: false}
}

type SearchResult struct {
	ID     string                 `json:"id"`
	Score  float64                `json:"score"`
	Fields map[string]interface{} `json:"fields"`
}

func (rvi *RedisVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	queryVectorJson, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	searchQueryStr := fmt.Sprintf("*=>[KNN %s @vector_field $vec AS score]", fmt.Sprint(maxResults))
	indexName := fmt.Sprintf("vector:%s", rvi.name)
	response, err := RedisClient.Do(ctx, "FT.SEARCH", indexName, searchQueryStr, "SORTBY", "score", "PARAMS", "2", "vec", queryVectorJson).Result()
	if err != nil {
		return nil, err
	}
	var searchResults []SearchResult
	if data, ok := response.(string); ok {
		if err := json.Unmarshal([]byte(data), &searchResults); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unexpected response type from Redis")
	}

	results := make(utils.MinTupleHeap, 0, len(searchResults))
	for _, result := range searchResults {
		results = append(results, utils.InitHeapElement(result.Score, result.Fields["uuid"].(string), false))
	}
	return results, nil
}

func (rvi *RedisVectorIndex) SearchWithUid(ctx context.Context, queryUid string, maxResults int, filter index.SearchFilter) (utils.MinTupleHeap, error) {
	queryVec, err := rvi.GetVector(ctx, queryUid)
	if err != nil {
		return nil, err
	}
	return rvi.Search(ctx, queryVec, maxResults, filter)
}

func (rvi *RedisVectorIndex) CreateRedisVectorIndex(ctx context.Context, vector []float32) error {

	var redisIndexType string
	switch rvi.indexType {
	case interfaces.SequentialManifestType:
		redisIndexType = "FLAT"
	case interfaces.HnswManifestType:
		redisIndexType = "HNSW"
	default:
		return interfaces.ErrInvalidVectorIndexType
	}

	indexName := fmt.Sprintf("vector:%s", rvi.name)
	_, err := RedisClient.Do(ctx, "FT.CREATE", indexName, "SCHEMA", "HASH", "PREFIX", "1",
		rvi.name, "SCHEMA", "vector_field", "VECTOR", redisIndexType,
		"6",
		"TYPE", "FLOAT32",
		"DIM", fmt.Sprint(len(vector)),
		"DISTANCE_METRIC", "L2").Result()
	if err != nil {
		return err
	}
	return nil
}

func (rvi *RedisVectorIndex) InsertVector(ctx context.Context, uuid string, vector []float32) error {
	rvi.mu.Lock()
	defer rvi.mu.Unlock()
	if !rvi.containsVectors {
		err := rvi.CreateRedisVectorIndex(ctx, vector)
		if err != nil {
			return err
		}
		rvi.containsVectors = true
	}

	vectorJson, err := json.Marshal(vector)
	if err != nil {
		return err
	}

	// Insert vector into Redis
	_, err = RedisClient.HSet(ctx, rvi.name, map[string]interface{}{
		"uuid":         uuid,
		"vector_field": vectorJson,
	}).Result()

	if err != nil {
		return err
	}

	return nil
}

func (rvi *RedisVectorIndex) DeleteVector(ctx context.Context, uuid string) error {
	rvi.mu.Lock()
	defer rvi.mu.Unlock()

	// Delete the vector fields
	err := RedisClient.HDel(ctx, rvi.name, "uuid", "vector_field").Err()
	if err != nil {
		return err
	}

	// Check if the hash is empty
	length, err := RedisClient.HLen(ctx, rvi.name).Result()
	if err != nil {
		return err
	}

	if length == 0 {
		rvi.containsVectors = false
	}

	return nil
}

func (rvi *RedisVectorIndex) GetVector(ctx context.Context, uuid string) ([]float32, error) {
	hash := fmt.Sprintf("vector:%s", uuid)
	result, err := RedisClient.HGet(ctx, hash, "vector_field").Result()
	if err != nil {
		return nil, err
	}
	var vector []float32
	err = json.Unmarshal([]byte(result), &vector)
	if err != nil {
		return nil, err
	}
	return vector, nil
}

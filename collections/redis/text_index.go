package redis

import (
	"context"
	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"sync"
)

type RedisCollection struct {
	name           string
	mu             sync.RWMutex
	VectorIndexMap map[string]*interfaces.VectorIndexWrapper
}

func NewCollection(name string) *RedisCollection {
	return &RedisCollection{
		name:           name,
		VectorIndexMap: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (rc *RedisCollection) GetName() string {
	return rc.name
}

func (rc *RedisCollection) GetVectorIndexMap() map[string]*interfaces.VectorIndexWrapper {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.VectorIndexMap
}

func (rc *RedisCollection) GetVectorIndex(name string) (*interfaces.VectorIndexWrapper, error) {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	if ind, ok := rc.VectorIndexMap[name]; !ok {
		return nil, index.ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (rc *RedisCollection) SetVectorIndex(name string, vectorIndex *interfaces.VectorIndexWrapper) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.VectorIndexMap == nil {
		rc.VectorIndexMap = map[string]*interfaces.VectorIndexWrapper{}
	}
	if _, ok := rc.VectorIndexMap[name]; ok {
		return index.ErrVectorIndexAlreadyExists
	}
	rc.VectorIndexMap[name] = vectorIndex
	return nil
}

func (rc *RedisCollection) DeleteVectorIndex(name string) error {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	delete(rc.VectorIndexMap, name)
	return nil
}

func (rc *RedisCollection) InsertText(ctx context.Context, uuid string, text string) error {
	return RedisClient.HSet(ctx, rc.name, uuid, text).Err()
}

func (rc *RedisCollection) DeleteText(ctx context.Context, uuid string) error {
	return RedisClient.HDel(ctx, rc.name, uuid).Err()
}

func (rc *RedisCollection) GetText(ctx context.Context, uuid string) (string, error) {
	return RedisClient.HGet(ctx, rc.name, uuid).Result()
}

func (rc *RedisCollection) GetTextMap(ctx context.Context) (map[string]string, error) {
	return RedisClient.HGetAll(ctx, rc.name).Result()
}

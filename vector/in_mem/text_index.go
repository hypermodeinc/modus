package in_mem

import (
	"context"
	"fmt"
	c "hmruntime/vector/constraints"
	"hmruntime/vector/index"
	"sync"
)

var (
	ErrVectorIndexAlreadyExists = fmt.Errorf("vector index already exists")
	ErrVectorIndexNotFound      = fmt.Errorf("vector index not found")
)

type InMemTextIndex[T c.Float] struct {
	mu            sync.RWMutex
	textMap       map[string]string
	vectorIndexes map[string]index.VectorIndex[T]
}

func (ti *InMemTextIndex[T]) GetVectorIndexMap() map[string]index.VectorIndex[T] {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.vectorIndexes
}

func (ti *InMemTextIndex[T]) GetVectorIndex(name string) (index.VectorIndex[T], error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.vectorIndexes[name]; !ok {
		return nil, ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemTextIndex[T]) SetVectorIndex(name string, index index.VectorIndex[T]) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if _, ok := ti.vectorIndexes[name]; ok {
		return ErrVectorIndexAlreadyExists
	}
	ti.vectorIndexes[name] = index
	return nil
}

func (ti *InMemTextIndex[T]) DeleteVectorIndex(name string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.vectorIndexes, name)
	return nil
}

func (ti *InMemTextIndex[T]) InsertText(ctx context.Context, c index.CacheType, uuid string, text string) ([]*index.KeyValue, error) {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.textMap[uuid] = text
	return nil, nil
}

func (ti *InMemTextIndex[T]) DeleteText(ctx context.Context, c index.CacheType, uuid string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.textMap, uuid)
	return nil
}

func (ti *InMemTextIndex[T]) GetText(ctx context.Context, c index.CacheType, uuid string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.textMap[uuid], nil
}

func (ti *InMemTextIndex[T]) GetTextMap() map[string]string {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.textMap
}

package in_mem

import (
	"context"
	"fmt"

	"hmruntime/vector/index"
	"sync"
)

var (
	ErrVectorIndexAlreadyExists = fmt.Errorf("vector index already exists")
	ErrVectorIndexNotFound      = fmt.Errorf("vector index not found")
)

type InMemTextIndex struct {
	mu            sync.RWMutex
	textMap       map[string]string
	vectorIndexes map[string]index.VectorIndex
}

func NewTextIndex() *InMemTextIndex {
	return &InMemTextIndex{
		textMap:       map[string]string{},
		vectorIndexes: map[string]index.VectorIndex{},
	}
}

func (ti *InMemTextIndex) GetVectorIndexMap() map[string]index.VectorIndex {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.vectorIndexes
}

func (ti *InMemTextIndex) GetVectorIndex(name string) (index.VectorIndex, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.vectorIndexes[name]; !ok {
		return nil, ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemTextIndex) SetVectorIndex(name string, index index.VectorIndex) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if _, ok := ti.vectorIndexes[name]; ok {
		return ErrVectorIndexAlreadyExists
	}
	ti.vectorIndexes[name] = index
	return nil
}

func (ti *InMemTextIndex) DeleteVectorIndex(name string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.vectorIndexes, name)
	return nil
}

func (ti *InMemTextIndex) InsertText(ctx context.Context, c index.CacheType, uuid string, text string) ([]*index.KeyValue, error) {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.textMap[uuid] = text
	return nil, nil
}

func (ti *InMemTextIndex) DeleteText(ctx context.Context, c index.CacheType, uuid string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.textMap, uuid)
	return nil
}

func (ti *InMemTextIndex) GetText(ctx context.Context, c index.CacheType, uuid string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.textMap[uuid], nil
}

func (ti *InMemTextIndex) GetTextMap() map[string]string {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.textMap
}

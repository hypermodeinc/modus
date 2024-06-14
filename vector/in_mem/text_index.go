package in_mem

import (
	"context"
	"fmt"

	"hmruntime/vector/index"
	"hmruntime/vector/index/interfaces"
	"sync"
)

var (
	ErrVectorIndexAlreadyExists = fmt.Errorf("vector index already exists")
	ErrVectorIndexNotFound      = fmt.Errorf("vector index not found")
)

type InMemTextIndex struct {
	mu            sync.RWMutex
	TextMap       map[string]string
	VectorIndexes map[string]*interfaces.VectorIndexWrapper
}

func NewTextIndex() *InMemTextIndex {
	return &InMemTextIndex{
		TextMap:       map[string]string{},
		VectorIndexes: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (ti *InMemTextIndex) GetVectorIndexMap() map[string]*interfaces.VectorIndexWrapper {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.VectorIndexes
}

func (ti *InMemTextIndex) GetVectorIndex(name string) (*interfaces.VectorIndexWrapper, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.VectorIndexes[name]; !ok {
		return nil, ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemTextIndex) SetVectorIndex(name string, index *interfaces.VectorIndexWrapper) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if _, ok := ti.VectorIndexes[name]; ok {
		return ErrVectorIndexAlreadyExists
	}
	ti.VectorIndexes[name] = index
	return nil
}

func (ti *InMemTextIndex) DeleteVectorIndex(name string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.VectorIndexes, name)
	return nil
}

func (ti *InMemTextIndex) InsertText(ctx context.Context, c index.CacheType, uuid string, text string) ([]*index.KeyValue, error) {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.TextMap[uuid] = text
	return nil, nil
}

func (ti *InMemTextIndex) DeleteText(ctx context.Context, c index.CacheType, uuid string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.TextMap, uuid)
	return nil
}

func (ti *InMemTextIndex) GetText(ctx context.Context, c index.CacheType, uuid string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap[uuid], nil
}

func (ti *InMemTextIndex) GetTextMap() map[string]string {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap
}

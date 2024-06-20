package in_mem

import (
	"context"
	"fmt"

	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"sync"
)

var (
	ErrVectorIndexAlreadyExists = fmt.Errorf("vector index already exists")
	ErrVectorIndexNotFound      = fmt.Errorf("vector index not found")
)

type InMemCollection struct {
	mu             sync.RWMutex
	TextMap        map[string]string
	VectorIndexMap map[string]*interfaces.VectorIndexWrapper
}

func NewCollection() *InMemCollection {
	return &InMemCollection{
		TextMap:        map[string]string{},
		VectorIndexMap: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (ti *InMemCollection) GetVectorIndexMap() map[string]*interfaces.VectorIndexWrapper {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.VectorIndexMap
}

func (ti *InMemCollection) GetVectorIndex(name string) (*interfaces.VectorIndexWrapper, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.VectorIndexMap[name]; !ok {
		return nil, ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemCollection) SetVectorIndex(name string, index *interfaces.VectorIndexWrapper) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.VectorIndexMap == nil {
		ti.VectorIndexMap = map[string]*interfaces.VectorIndexWrapper{}
	}
	if _, ok := ti.VectorIndexMap[name]; ok {
		return ErrVectorIndexAlreadyExists
	}
	ti.VectorIndexMap[name] = index
	return nil
}

func (ti *InMemCollection) DeleteVectorIndex(name string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.VectorIndexMap, name)
	return nil
}

func (ti *InMemCollection) InsertText(ctx context.Context, uuid string, text string) ([]*index.KeyValue, error) {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.TextMap[uuid] = text
	return nil, nil
}

func (ti *InMemCollection) DeleteText(ctx context.Context, uuid string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.TextMap, uuid)
	return nil
}

func (ti *InMemCollection) GetText(ctx context.Context, uuid string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap[uuid], nil
}

func (ti *InMemCollection) GetTextMap() map[string]string {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap
}

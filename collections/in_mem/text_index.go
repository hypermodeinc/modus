package in_mem

import (
	"context"

	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"sync"
)

type InMemCollection struct {
	mu             sync.RWMutex
	name           string
	TextMap        map[string]string
	VectorIndexMap map[string]*interfaces.VectorIndexWrapper
}

func NewCollection(name string) *InMemCollection {
	return &InMemCollection{
		TextMap:        map[string]string{},
		VectorIndexMap: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (ti *InMemCollection) GetName() string {
	return ti.name
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
		return nil, index.ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemCollection) SetVectorIndex(name string, vectorIndex *interfaces.VectorIndexWrapper) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.VectorIndexMap == nil {
		ti.VectorIndexMap = map[string]*interfaces.VectorIndexWrapper{}
	}
	if _, ok := ti.VectorIndexMap[name]; ok {
		return index.ErrVectorIndexAlreadyExists
	}
	ti.VectorIndexMap[name] = vectorIndex
	return nil
}

func (ti *InMemCollection) DeleteVectorIndex(name string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.VectorIndexMap, name)
	return nil
}

func (ti *InMemCollection) InsertText(ctx context.Context, uuid string, text string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.TextMap[uuid] = text
	return nil
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

func (ti *InMemCollection) GetTextMap(ctx context.Context) (map[string]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap, nil
}

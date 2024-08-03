package in_mem

import (
	"context"
	"fmt"

	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"hmruntime/db"
	"sync"
)

type InMemCollection struct {
	mu             sync.RWMutex
	collectionName string
	lastInsertedID int64
	TextMap        map[string]string // key: text
	LabelsMap      map[string][]string
	IdMap          map[string]int64                          // key: postgres id
	VectorIndexMap map[string]*interfaces.VectorIndexWrapper // searchMethod: vectorIndex
}

func NewCollection(name string) *InMemCollection {
	return &InMemCollection{
		collectionName: name,
		TextMap:        map[string]string{},
		LabelsMap:      map[string][]string{},
		IdMap:          map[string]int64{},
		VectorIndexMap: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (ti *InMemCollection) GetCollectionName() string {
	return ti.collectionName
}

func (ti *InMemCollection) GetVectorIndexMap() map[string]*interfaces.VectorIndexWrapper {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.VectorIndexMap
}

func (ti *InMemCollection) GetVectorIndex(ctx context.Context, searchMethod string) (*interfaces.VectorIndexWrapper, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.VectorIndexMap[searchMethod]; !ok {
		return nil, index.ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemCollection) SetVectorIndex(ctx context.Context, searchMethod string, vectorIndex *interfaces.VectorIndexWrapper) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.VectorIndexMap == nil {
		ti.VectorIndexMap = map[string]*interfaces.VectorIndexWrapper{}
	}
	if _, ok := ti.VectorIndexMap[searchMethod]; ok {
		return index.ErrVectorIndexAlreadyExists
	}
	ti.VectorIndexMap[searchMethod] = vectorIndex
	return nil
}

func (ti *InMemCollection) DeleteVectorIndex(ctx context.Context, searchMethod string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.VectorIndexMap, searchMethod)
	return nil
}

func (ti *InMemCollection) InsertTexts(ctx context.Context, keys []string, texts []string, labelsArr [][]string) error {
	if len(keys) != len(texts) {
		return fmt.Errorf("keys and texts must have the same length")
	}

	if len(labelsArr) != 0 && len(labelsArr) != len(keys) {
		return fmt.Errorf("labels must have the same length as keys or be empty")
	}

	ids, err := db.WriteCollectionTexts(ctx, ti.collectionName, keys, texts, labelsArr)
	if err != nil {
		return err
	}

	return ti.InsertTextsToMemory(ctx, ids, keys, texts, labelsArr)
}

func (ti *InMemCollection) InsertText(ctx context.Context, key string, text string, labels []string) error {
	id, err := db.WriteCollectionText(ctx, ti.collectionName, key, text, labels)
	if err != nil {
		return err
	}

	return ti.InsertTextToMemory(ctx, id, key, text, labels)
}

func (ti *InMemCollection) InsertTextsToMemory(ctx context.Context, ids []int64, keys []string, texts []string, labelsArr [][]string) error {

	if len(labelsArr) != 0 && len(labelsArr) != len(keys) {
		return fmt.Errorf("labels must have the same length as keys or be empty")
	}
	if len(ids) != len(keys) || len(ids) != len(texts) {
		return fmt.Errorf("ids, keys and texts must have the same length")
	}

	ti.mu.Lock()
	defer ti.mu.Unlock()
	for i, key := range keys {
		ti.TextMap[key] = texts[i]
		if len(labelsArr) != 0 {
			ti.LabelsMap[key] = labelsArr[i]
		}
		ti.IdMap[key] = ids[i]
		ti.lastInsertedID = ids[i]
	}
	return nil
}

func (ti *InMemCollection) InsertTextToMemory(ctx context.Context, id int64, key string, text string, labels []string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.TextMap[key] = text
	if len(labels) != 0 {
		ti.LabelsMap[key] = labels
	}
	ti.IdMap[key] = id
	ti.lastInsertedID = id
	return nil
}

func (ti *InMemCollection) DeleteText(ctx context.Context, key string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	err := db.DeleteCollectionText(ctx, ti.collectionName, key)
	if err != nil {
		return err
	}
	delete(ti.TextMap, key)
	return nil
}

func (ti *InMemCollection) GetText(ctx context.Context, key string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap[key], nil
}

func (ti *InMemCollection) GetTextMap(ctx context.Context) (map[string]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap, nil
}

func (ti *InMemCollection) GetLabels(ctx context.Context, key string) ([]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.LabelsMap[key], nil
}

func (ti *InMemCollection) GetLabelsMap(ctx context.Context) (map[string][]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.LabelsMap, nil
}

func (ti *InMemCollection) Len(ctx context.Context) (int, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return len(ti.TextMap), nil
}

func (ti *InMemCollection) GetExternalId(ctx context.Context, key string) (int64, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.IdMap[key], nil
}

func (ti *InMemCollection) GetCheckpointId(ctx context.Context) (int64, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.lastInsertedID, nil
}

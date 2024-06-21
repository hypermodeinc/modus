package collections

import (
	"context"
	"errors"
	"fmt"

	"hmruntime/collections/index/interfaces"
	"hmruntime/db"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/pluginmanager"
	"hmruntime/plugins"

	"sync"
	"time"
)

const collectionFactoryWriteInterval = 1

var (
	GlobalCollectionFactory *CollectionFactory
	ErrCollectionNotFound   = fmt.Errorf("collection not found")
)

func InitializeIndexFactory(ctx context.Context) {
	GlobalCollectionFactory = CreateFactory()
	manifestdata.RegisterManifestLoadedCallback(CleanAndProcessManifest)
	pluginmanager.RegisterPluginLoadedCallback(func(ctx context.Context, metadata plugins.PluginMetadata) error {
		CatchEmbedderReqs(ctx)
		return nil
	})
	go GlobalCollectionFactory.worker(ctx)
}

func CloseIndexFactory(ctx context.Context) {
	close(GlobalCollectionFactory.quit)
	<-GlobalCollectionFactory.done
}

type CollectionFactory struct {
	collectionMap map[string]interfaces.Collection
	mu            sync.RWMutex
	quit          chan struct{}
	done          chan struct{}
}

func (tif *CollectionFactory) worker(ctx context.Context) {
	defer close(tif.done)
	ticker := time.NewTicker(collectionFactoryWriteInterval * time.Minute)

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// read from postgres all collections & searchMethod after lastInsertedID
			tif.ReadFromPostgres(ctx)
		case <-tif.quit:
			return
		}
	}
}

func CreateFactory() *CollectionFactory {
	f := &CollectionFactory{
		collectionMap: map[string]interfaces.Collection{},
		quit:          make(chan struct{}),
		done:          make(chan struct{}),
	}
	return f
}

func (hf *CollectionFactory) isNameAvailableWithLock(name string) bool {
	_, nameUsed := hf.collectionMap[name]
	return !nameUsed
}

func (hf *CollectionFactory) Create(
	ctx context.Context,
	name string,
	index interfaces.Collection) (interfaces.Collection, error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.createWithLock(name, index)
}

func (hf *CollectionFactory) createWithLock(
	name string,
	index interfaces.Collection) (interfaces.Collection, error) {
	if !hf.isNameAvailableWithLock(name) {
		err := errors.New("index with name " + name + " already exists")
		return nil, err
	}
	retVal := index
	hf.collectionMap[name] = retVal
	return retVal, nil
}

func (hf *CollectionFactory) GetCollectionMap() map[string]interfaces.Collection {
	return hf.collectionMap
}

func (hf *CollectionFactory) Find(ctx context.Context, name string) (interfaces.Collection, error) {
	hf.mu.RLock()
	defer hf.mu.RUnlock()
	return hf.findWithLock(name)
}

func (hf *CollectionFactory) findWithLock(name string) (interfaces.Collection, error) {
	vecInd, ok := hf.collectionMap[name]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return vecInd, nil
}

func (hf *CollectionFactory) Remove(ctx context.Context, name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	err := db.DeleteCollectionTexts(ctx, name)
	if err != nil {
		return err
	}
	return hf.removeWithLock(name)
}

func (hf *CollectionFactory) removeWithLock(name string) error {
	delete(hf.collectionMap, name)
	return nil
}

func (hf *CollectionFactory) CreateOrReplace(
	ctx context.Context,
	name string,
	index interfaces.Collection) (interfaces.Collection, error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	vi, err := hf.findWithLock(name)
	if err != nil {
		return nil, err
	}
	if vi != nil {
		err = hf.removeWithLock(name)
		if err != nil {
			return nil, err
		}
	}
	return hf.createWithLock(name, index)
}

func (hf *CollectionFactory) ReadFromPostgres(ctx context.Context) {
	for _, collection := range hf.collectionMap {
		err := LoadTextsIntoCollection(ctx, collection)
		if err != nil {
			logger.Err(ctx, err).
				Str("collection_name", collection.GetCollectionName()).
				Msg("Failed to load texts into collection.")
		}

		for _, vectorIndex := range collection.GetVectorIndexMap() {
			err = LoadVectorsIntoVectorIndex(ctx, vectorIndex, collection)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collection.GetCollectionName()).
					Str("search_method", vectorIndex.GetSearchMethodName()).
					Msg("Failed to load vectors into vector index.")
			}
		}
	}
}

func LoadTextsIntoCollection(ctx context.Context, collection interfaces.Collection) error {
	// Get checkpoint id for collection
	textCheckpointId, err := collection.GetCheckpointId(ctx)
	if err != nil {
		return err
	}

	// Query all texts from checkpoint
	textIds, keys, texts, err := db.QueryCollectionTextsFromCheckpoint(ctx, collection.GetCollectionName(), textCheckpointId)
	if err != nil {
		return err
	}
	if len(textIds) != len(keys) || len(keys) != len(texts) {
		return errors.New("mismatch in keys and texts")
	}

	// Insert all texts into collection
	for i := range textIds {
		err = collection.InsertTextToMemory(ctx, textIds[i], keys[i], texts[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func LoadVectorsIntoVectorIndex(ctx context.Context, vectorIndex *interfaces.VectorIndexWrapper, collection interfaces.Collection) error {
	// Get checkpoint id for vector index
	vecCheckpointId, err := vectorIndex.GetCheckpointId(ctx)
	if err != nil {
		return err
	}

	// Query all vectors from checkpoint
	vectorIds, keys, vectors, err := db.QueryCollectionVectorsFromCheckpoint(ctx, collection.GetCollectionName(), vectorIndex.GetSearchMethodName(), vecCheckpointId)
	if err != nil {
		return err
	}
	if len(vectorIds) != len(vectors) || len(keys) != len(vectors) {
		return errors.New("mismatch in keys and vectors")
	}

	// Insert all vectors into vector index
	for i := range vectorIds {
		err = vectorIndex.InsertVectorToMemory(ctx, vectorIds[i], keys[i], vectors[i])
		if err != nil {
			return err
		}
	}
	return nil
}

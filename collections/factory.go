package collections

import (
	"context"
	"errors"
	"fmt"

	"hmruntime/collections/index/interfaces"
	"hmruntime/db"
	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/manifestdata"

	"sync"
	"time"
)

const collectionFactoryWriteInterval = 1

var (
	GlobalNamespaceManager *CollectionFactory
	ErrCollectionNotFound  = fmt.Errorf("collection not found")
	ErrNamespaceNotFound   = fmt.Errorf("namespace not found")
)

func InitializeIndexFactory(ctx context.Context) {
	GlobalNamespaceManager = CreateFactory()
	manifestdata.RegisterManifestLoadedCallback(CleanAndProcessManifest)
	functions.RegisterFunctionsLoadedCallback(func(ctx context.Context) {
		GlobalNamespaceManager.ReadFromPostgres(ctx)
	})

	go GlobalNamespaceManager.worker(ctx)
}

func CloseIndexFactory(ctx context.Context) {
	close(GlobalNamespaceManager.quit)
	<-GlobalNamespaceManager.done
}

type CollectionFactory struct {
	collectionMap map[string]*Collection
	mu            sync.RWMutex
	quit          chan struct{}
	done          chan struct{}
}

type Collection struct {
	collectionNamespaceMap map[string]interfaces.CollectionNamespace
	mu                     sync.RWMutex
}

func NewCollection(ctx context.Context, collectionName string) *Collection {
	return &Collection{
		collectionNamespaceMap: map[string]interfaces.CollectionNamespace{},
	}
}

func (cf *CollectionFactory) worker(ctx context.Context) {
	defer close(cf.done)
	timer := time.NewTimer(collectionFactoryWriteInterval * time.Minute)

	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// read from postgres all collections & searchMethod after lastInsertedID
			cf.ReadFromPostgres(ctx)
			timer.Reset(collectionFactoryWriteInterval * time.Minute)
		case <-cf.quit:
			return
		}
	}
}

func CreateFactory() *CollectionFactory {
	f := &CollectionFactory{
		collectionMap: map[string]*Collection{
			"": {
				collectionNamespaceMap: map[string]interfaces.CollectionNamespace{},
			},
		},
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
	return f
}

func (cf *CollectionFactory) isCollectionNameAvailableWithLock(name string) bool {
	_, nameUsed := cf.collectionMap[name]
	return !nameUsed
}

func (cf *CollectionFactory) CreateCollection(
	ctx context.Context,
	name string,
	coll *Collection) (*Collection, error) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.createCollectionWithLock(name, coll)
}

func (cf *CollectionFactory) createCollectionWithLock(
	name string,
	coll *Collection) (*Collection, error) {
	if !cf.isCollectionNameAvailableWithLock(name) {
		return nil, errors.New(fmt.Sprintf("collection with name %s already exists", name))
	}
	cf.collectionMap[name] = coll
	return coll, nil
}

func (cf *CollectionFactory) GetNamespaceCollectionFactoryMap() map[string]*Collection {
	return cf.collectionMap
}

func (cf *CollectionFactory) FindCollection(ctx context.Context, name string) (*Collection, error) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()
	return cf.findCollectionWithLock(name)
}

func (cf *CollectionFactory) findCollectionWithLock(name string) (*Collection, error) {
	coll, ok := cf.collectionMap[name]
	if !ok {
		return nil, ErrCollectionNotFound
	}
	return coll, nil
}

func (cf *CollectionFactory) RemoveCollection(ctx context.Context, name string) error {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	return cf.removeCollectionWithLock(name)
}

func (cf *CollectionFactory) removeCollectionWithLock(name string) error {
	delete(cf.collectionMap, name)
	return nil
}

func (cf *CollectionFactory) CreateOrReplaceCollection(
	ctx context.Context,
	name string,
	coll *Collection) (*Collection, error) {
	cf.mu.Lock()
	defer cf.mu.Unlock()
	vi, err := cf.findCollectionWithLock(name)
	if err != nil {
		return nil, err
	}
	if vi != nil {
		err = cf.removeCollectionWithLock(name)
		if err != nil {
			return nil, err
		}
	}
	return cf.createCollectionWithLock(name, coll)
}

func (c *Collection) isNamespaceAvailableWithLock(namespace string) bool {
	_, namespaceUsed := c.collectionNamespaceMap[namespace]
	return !namespaceUsed
}

func (c *Collection) CreateCollectionNamespace(
	ctx context.Context,
	namespace string,
	index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.createNamespaceWithLock(namespace, index)
}

func (c *Collection) createNamespaceWithLock(
	namespace string,
	index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	if !c.isNamespaceAvailableWithLock(namespace) {
		return nil, errors.New(fmt.Sprintf("namespace with name %s already exists", namespace))
	}
	c.collectionNamespaceMap[namespace] = index
	return index, nil
}

func (c *Collection) FindNamespace(ctx context.Context, namespace string) (interfaces.CollectionNamespace, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.findNamespaceWithLock(namespace)
}

func (c *Collection) findNamespaceWithLock(namespace string) (interfaces.CollectionNamespace, error) {
	ns, ok := c.collectionNamespaceMap[namespace]
	if !ok {
		return nil, ErrNamespaceNotFound
	}
	return ns, nil
}

func (c *Collection) RemoveNamespace(ctx context.Context, namespace string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.removeNamespaceWithLock(namespace)
}

func (c *Collection) removeNamespaceWithLock(namespace string) error {
	delete(c.collectionNamespaceMap, namespace)
	return nil
}

func (c *Collection) FindOrCreateNamespace(
	ctx context.Context,
	namespace string,
	index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	vi, err := c.findNamespaceWithLock(namespace)
	if err == ErrNamespaceNotFound {
		return c.createNamespaceWithLock(namespace, index)
	} else if err != nil {
		return nil, err
	}
	return vi, nil
}

func (c *Collection) CreateOrReplaceNamespace(
	ctx context.Context,
	namespace string,
	index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	vi, err := c.findNamespaceWithLock(namespace)
	if err != nil {
		return nil, err
	}
	if vi != nil {
		err = c.removeNamespaceWithLock(namespace)
		if err != nil {
			return nil, err
		}
	}
	return c.createNamespaceWithLock(namespace, index)
}

func (cf *CollectionFactory) ReadFromPostgres(ctx context.Context) {
	for _, namespaceCollectionFactory := range cf.collectionMap {
		for _, collection := range namespaceCollectionFactory.collectionNamespaceMap {
			err := LoadTextsIntoCollection(ctx, collection)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collection.GetCollectionName()).
					Msg("Failed to load texts into collection.")
				break
			}

			for _, vectorIndex := range collection.GetVectorIndexMap() {
				err = LoadVectorsIntoVectorIndex(ctx, vectorIndex, collection)
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", collection.GetCollectionName()).
						Str("search_method", vectorIndex.GetSearchMethodName()).
						Msg("Failed to load vectors into vector index.")
					break
				}

				// catch up on any texts that weren't embedded
				err := syncTextsWithVectorIndex(ctx, collection, vectorIndex)
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", collection.GetCollectionName()).
						Str("search_method", vectorIndex.GetSearchMethodName()).
						Msg("Failed to sync text with vector index.")
					break
				}

			}
		}
	}
}

func LoadTextsIntoCollection(ctx context.Context, collection interfaces.CollectionNamespace) error {
	// Get checkpoint id for collection
	textCheckpointId, err := collection.GetCheckpointId(ctx)
	if err != nil {
		return err
	}

	// Query all texts from checkpoint
	textIds, keys, texts, labels, err := db.QueryCollectionTextsFromCheckpoint(ctx, collection.GetCollectionName(), collection.GetNamespace(), textCheckpointId)
	if err != nil {
		return err
	}
	if len(textIds) != len(keys) || len(keys) != len(texts) {
		return errors.New("mismatch in keys and texts")
	}

	// Insert all texts into collection
	err = collection.InsertTextsToMemory(ctx, textIds, keys, texts, labels)
	if err != nil {
		return err
	}

	return nil
}

func LoadVectorsIntoVectorIndex(ctx context.Context, vectorIndex interfaces.VectorIndex, collection interfaces.CollectionNamespace) error {
	// Get checkpoint id for vector index
	vecCheckpointId, err := vectorIndex.GetCheckpointId(ctx)
	if err != nil {
		return err
	}

	// Query all vectors from checkpoint
	textIds, vectorIds, keys, vectors, err := db.QueryCollectionVectorsFromCheckpoint(ctx, collection.GetCollectionName(), vectorIndex.GetSearchMethodName(), vectorIndex.GetNamespace(), vecCheckpointId)
	if err != nil {
		return err
	}
	if len(vectorIds) != len(vectors) || len(keys) != len(vectors) {
		return errors.New("mismatch in keys and vectors")
	}

	// Insert all vectors into vector index
	err = batchInsertVectorsToMemory(ctx, vectorIndex, textIds, vectorIds, keys, vectors)
	if err != nil {
		return err
	}

	return nil
}

func syncTextsWithVectorIndex(ctx context.Context, collection interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex) error {
	// Query all texts from checkpoint
	lastIndexedTextId, err := vectorIndex.GetLastIndexedTextId(ctx)
	if err != nil {
		return err
	}
	textIds, keys, texts, _, err := db.QueryCollectionTextsFromCheckpoint(ctx, collection.GetCollectionName(), collection.GetNamespace(), lastIndexedTextId)
	if err != nil {
		return err
	}
	if len(textIds) != len(keys) || len(keys) != len(texts) {
		return errors.New("mismatch in keys and texts")
	}
	// Get last indexed text id
	err = ProcessTexts(ctx, collection, vectorIndex, keys, texts)
	if err != nil {
		return err
	}

	return nil
}

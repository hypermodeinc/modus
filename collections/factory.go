package collections

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"hmruntime/collections/in_mem"
	"hmruntime/collections/index/interfaces"
	"hmruntime/config"
	"hmruntime/logger"
	"hmruntime/storage"

	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const collectionFactoryWriteInterval = 1

var (
	GlobalCollectionFactory *CollectionFactory
	ErrCollectionNotFound   = fmt.Errorf("text index not found")
)

func InitializeIndexFactory(ctx context.Context) {
	GlobalCollectionFactory = CreateFactory()
	err := GlobalCollectionFactory.ReadFromBin()
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error reading index factory from bin")
	}
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
	var ticker *time.Ticker
	if config.GetEnvironmentName() == config.DevEnvironmentName {
		ticker = time.NewTicker(collectionFactoryWriteInterval * time.Minute)
	} else {
		ticker = time.NewTicker(collectionFactoryWriteInterval * time.Hour)
	}

	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			err := tif.WriteToBin()
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Error writing index factory to bin")
			}
		case <-tif.quit:
			err := tif.WriteToBin()
			if err != nil {
				logger.Error(ctx).Err(err).Msg("Error writing index factory to bin")
			}
			close(tif.done)
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

func (hf *CollectionFactory) Find(name string) (interfaces.Collection, error) {
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

func (hf *CollectionFactory) Remove(name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.removeWithLock(name)
}

func (hf *CollectionFactory) removeWithLock(name string) error {
	delete(hf.collectionMap, name)
	return nil
}

func (hf *CollectionFactory) CreateOrReplace(
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

func (hf *CollectionFactory) WriteToBin() error {
	operation := func() error {
		data, err := json.Marshal(hf.collectionMap)
		if err != nil {
			return fmt.Errorf("could not encode file content, %s", err)
		}

		if err := storage.WriteFile(context.Background(), "index_factory.bin", data); err != nil {
			return fmt.Errorf("could not write to file, %s", err)
		}
		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)
}

func (hf *CollectionFactory) ReadFromBin() error {
	operation := func() error {
		data, err := storage.GetFileContents(context.Background(), "index_factory.bin")
		if err != nil {
			return fmt.Errorf("could not get file content, %s", err)
		}

		newMap := make(map[string]*in_mem.InMemCollection)

		if err := json.Unmarshal(data, &newMap); err != nil {
			return fmt.Errorf("could not decode file content, %s", err)
		}

		hf.collectionMap = make(map[string]interfaces.Collection)
		for k, v := range newMap {
			hf.collectionMap[k] = v
		}

		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)
}

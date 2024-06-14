package vector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"hmruntime/logger"
	"hmruntime/storage"
	"hmruntime/vector/in_mem"
	"hmruntime/vector/index"
	"hmruntime/vector/index/interfaces"

	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

const textIndexFactoryWriteInterval = 1 * time.Hour

var (
	GlobalTextIndexFactory *TextIndexFactory
	ErrTextIndexNotFound   = fmt.Errorf("text index not found")
)

func InitializeIndexFactory(ctx context.Context) {
	GlobalTextIndexFactory = CreateFactory()
	err := GlobalTextIndexFactory.ReadFromBin()
	if err != nil {
		logger.Error(ctx).Err(err).Msg("Error reading index factory from bin")
	}
	go GlobalTextIndexFactory.worker(ctx)
}

func CloseIndexFactory(ctx context.Context) {
	close(GlobalTextIndexFactory.quit)
	<-GlobalTextIndexFactory.done
}

type TextIndexFactory struct {
	textIndexMap map[string]interfaces.TextIndex
	mu           sync.RWMutex
	quit         chan struct{}
	done         chan struct{}
}

func (tif *TextIndexFactory) worker(ctx context.Context) {
	ticker := time.NewTicker(textIndexFactoryWriteInterval)
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

func CreateFactory() *TextIndexFactory {
	f := &TextIndexFactory{
		textIndexMap: map[string]interfaces.TextIndex{},
		quit:         make(chan struct{}),
		done:         make(chan struct{}),
	}
	return f
}

func (hf *TextIndexFactory) isNameAvailableWithLock(name string) bool {
	_, nameUsed := hf.textIndexMap[name]
	return !nameUsed
}

func (hf *TextIndexFactory) Create(
	name string,
	index interfaces.TextIndex) (interfaces.TextIndex, error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.createWithLock(name, index)
}

func (hf *TextIndexFactory) createWithLock(
	name string,
	index interfaces.TextIndex) (interfaces.TextIndex, error) {
	if !hf.isNameAvailableWithLock(name) {
		err := errors.New("index with name " + name + " already exists")
		return nil, err
	}
	retVal := index
	hf.textIndexMap[name] = retVal
	return retVal, nil
}

func (hf *TextIndexFactory) GetTextIndexMap() map[string]interfaces.TextIndex {
	return hf.textIndexMap
}

func (hf *TextIndexFactory) Find(name string) (interfaces.TextIndex, error) {
	hf.mu.RLock()
	defer hf.mu.RUnlock()
	return hf.findWithLock(name)
}

func (hf *TextIndexFactory) findWithLock(name string) (interfaces.TextIndex, error) {
	vecInd, ok := hf.textIndexMap[name]
	if !ok {
		return nil, ErrTextIndexNotFound
	}
	return vecInd, nil
}

func (hf *TextIndexFactory) Remove(name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.removeWithLock(name)
}

func (hf *TextIndexFactory) removeWithLock(name string) error {
	delete(hf.textIndexMap, name)
	return nil
}

func (hf *TextIndexFactory) CreateOrReplace(
	name string,
	source index.VectorSource,
	index interfaces.TextIndex) (interfaces.TextIndex, error) {
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

func (hf *TextIndexFactory) WriteToBin() error {
	operation := func() error {
		data, err := json.Marshal(hf.textIndexMap)
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

func (hf *TextIndexFactory) ReadFromBin() error {
	operation := func() error {
		data, err := storage.GetFileContents(context.Background(), "index_factory.bin")
		if err != nil {
			return fmt.Errorf("could not get file content, %s", err)
		}

		newMap := make(map[string]*in_mem.InMemTextIndex)

		if err := json.Unmarshal(data, &newMap); err != nil {
			return fmt.Errorf("could not decode file content, %s", err)
		}

		hf.textIndexMap = make(map[string]interfaces.TextIndex)
		for k, v := range newMap {
			hf.textIndexMap[k] = v
		}

		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)
}

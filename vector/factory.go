package vector

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"hmruntime/storage"

	"hmruntime/vector/index"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var (
	GlobalTextIndexFactory *TextIndexFactory
	ErrTextIndexNotFound   = fmt.Errorf("text index not found")
)

func InitializeIndexFactory() {
	GlobalTextIndexFactory = CreateFactory()
	err := GlobalTextIndexFactory.ReadFromWAL()
	if err != nil {
		fmt.Println("Error reading from WAL, ", err)
	}
}

func CloseIndexFactory() {
	err := GlobalTextIndexFactory.WriteToWAL()
	if err != nil {
		fmt.Println("Error writing to WAL, ", err)
	}
}

type TextIndexFactory struct {
	textIndexMap map[string]index.TextIndex
	mu           sync.RWMutex
}

// CreateFactory creates an instance of the private struct IndexFactory.
// NOTE: if T and floatBits do not match in # of bits, there will be consequences.
func CreateFactory() *TextIndexFactory {
	f := &TextIndexFactory{
		textIndexMap: map[string]index.TextIndex{},
	}
	return f
}

func (hf *TextIndexFactory) isNameAvailableWithLock(name string) bool {
	_, nameUsed := hf.textIndexMap[name]
	return !nameUsed
}

// Create is an implementation of the IndexFactory interface function, invoked by an HNSWIndexFactory
// instance. It takes in a string name and a VectorSource implementation, and returns a VectorIndex and error
// flag. It creates an HNSW instance using the index name and populates other parts of the HNSW struct such as
// multFactor, maxLevels, efConstruction, maxNeighbors, and efSearch using struct parameters.
// It then populates the HNSW graphs using the InsertChunk function until there are no more items to populate.
// Finally, the function adds the name and hnsw object to the in memory map and returns the object.
func (hf *TextIndexFactory) Create(
	name string,
	index index.TextIndex) (index.TextIndex, error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.createWithLock(name, index)
}

func (hf *TextIndexFactory) createWithLock(
	name string,
	index index.TextIndex) (index.TextIndex, error) {
	if !hf.isNameAvailableWithLock(name) {
		err := errors.New("index with name " + name + " already exists")
		return nil, err
	}
	retVal := index
	hf.textIndexMap[name] = retVal
	return retVal, nil
}

func (hf *TextIndexFactory) GetTextIndexMap() map[string]index.TextIndex {
	return hf.textIndexMap
}

// Find is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It returns the VectorIndex corresponding with a string name using the in memory map.
func (hf *TextIndexFactory) Find(name string) (index.TextIndex, error) {
	hf.mu.RLock()
	defer hf.mu.RUnlock()
	return hf.findWithLock(name)
}

func (hf *TextIndexFactory) findWithLock(name string) (index.TextIndex, error) {
	vecInd, ok := hf.textIndexMap[name]
	if !ok {
		return nil, ErrTextIndexNotFound
	}
	return vecInd, nil
}

// Remove is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It removes the VectorIndex corresponding with a string name using the in memory map.
func (hf *TextIndexFactory) Remove(name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.removeWithLock(name)
}

func (hf *TextIndexFactory) removeWithLock(name string) error {
	delete(hf.textIndexMap, name)
	return nil
}

// CreateOrReplace is an implementation of the IndexFactory interface funciton,
// invoked by an persistentIndexFactory. It checks if a VectorIndex
// correpsonding with name exists. If it does, it removes it, and replaces it
// via the Create function using the passed VectorSource. If the VectorIndex
// does not exist, it creates that VectorIndex corresponding with the name using
// the VectorSource.
func (hf *TextIndexFactory) CreateOrReplace(
	name string,
	source index.VectorSource,
	index index.TextIndex) (index.TextIndex, error) {
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

func (hf *TextIndexFactory) WriteToWAL() error {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	operation := func() error {
		if err := encoder.Encode(hf.textIndexMap); err != nil {
			return fmt.Errorf("could not encode file content, %s", err)
		}

		if err := storage.WriteFile(context.Background(), "index.wal", buf.Bytes()); err != nil {
			return fmt.Errorf("could not write to file, %s", err)
		}
		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)
}

func (hf *TextIndexFactory) ReadFromWAL() error {

	operation := func() error {
		data, err := storage.GetFileContents(context.Background(), "index.wal")
		if err != nil {
			return fmt.Errorf("could not get file content, %s", err)
		}

		decoder := gob.NewDecoder(bytes.NewReader(data))

		if err := decoder.Decode(&hf.textIndexMap); err != nil {
			return fmt.Errorf("could not decode file content, %s", err)
		}

		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)

}

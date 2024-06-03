package vector

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"hmruntime/storage"
	c "hmruntime/vector/constraints"
	"hmruntime/vector/index"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"

	"hmruntime/vector/options"
)

var GlobalIndexFactory *IndexFactory[float64]

func InitializeIndexFactory() {
	GlobalIndexFactory = CreateFactory[float64]()
	err := GlobalIndexFactory.ReadFromWAL()
	if err != nil {
		fmt.Println("Error reading from WAL, ", err)
	}
}

func CloseIndexFactory() {
	err := GlobalIndexFactory.WriteToWAL()
	if err != nil {
		fmt.Println("Error writing to WAL, ", err)
	}
}

type IndexFactory[T c.Float] struct {
	indexMap map[string]index.VectorIndex[T]
	mu       sync.RWMutex
}

// CreateFactory creates an instance of the private struct IndexFactory.
// NOTE: if T and floatBits do not match in # of bits, there will be consequences.
func CreateFactory[T c.Float]() *IndexFactory[T] {
	f := &IndexFactory[T]{
		indexMap: map[string]index.VectorIndex[T]{},
	}
	return f
}

// Implements NamedFactory interface for use as a plugin.
func (hf *IndexFactory[T]) Name() string { return "in_mem" }

func (hf *IndexFactory[T]) isNameAvailableWithLock(name string) bool {
	_, nameUsed := hf.indexMap[name]
	return !nameUsed
}

// hf.AllowedOptions() allows persistentIndexFactory to implement the
// IndexFactory interface (see vector-indexer/index/index.go for details).
// We define here options for exponent, maxLevels, efSearch, efConstruction,
// and metric.
func (hf *IndexFactory[T]) AllowedOptions() options.AllowedOptions {
	retVal := options.NewAllowedOptions()

	return retVal
}

// Create is an implementation of the IndexFactory interface function, invoked by an HNSWIndexFactory
// instance. It takes in a string name and a VectorSource implementation, and returns a VectorIndex and error
// flag. It creates an HNSW instance using the index name and populates other parts of the HNSW struct such as
// multFactor, maxLevels, efConstruction, maxNeighbors, and efSearch using struct parameters.
// It then populates the HNSW graphs using the InsertChunk function until there are no more items to populate.
// Finally, the function adds the name and hnsw object to the in memory map and returns the object.
func (hf *IndexFactory[T]) Create(
	name string,
	o options.Options,
	source index.VectorSource[T],
	index index.VectorIndex[T]) (index.VectorIndex[T], error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.createWithLock(name, o, source, index)
}

func (hf *IndexFactory[T]) createWithLock(
	name string,
	o options.Options,
	source index.VectorSource[T],
	index index.VectorIndex[T]) (index.VectorIndex[T], error) {
	if !hf.isNameAvailableWithLock(name) {
		err := errors.New("index with name " + name + " already exists")
		return nil, err
	}
	retVal := index
	// while NextChunk() returns a chunk of vectors and their corresponding uids, insert them into the index graph
	for chunk, uids, more, err := source.NextChunk(); more; {
		if err != nil {
			return nil, err
		}
		for i, vector := range chunk {
			_, err := retVal.Insert(context.Background(), nil, uids[i], vector)
			if err != nil {
				return nil, err
			}
		}
	}
	err := retVal.ApplyOptions(o)
	if err != nil {
		return nil, err
	}
	hf.indexMap[name] = retVal
	return retVal, nil
}

// Find is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It returns the VectorIndex corresponding with a string name using the in memory map.
func (hf *IndexFactory[T]) Find(name string) (index.VectorIndex[T], error) {
	hf.mu.RLock()
	defer hf.mu.RUnlock()
	return hf.findWithLock(name)
}

func (hf *IndexFactory[T]) findWithLock(name string) (index.VectorIndex[T], error) {
	vecInd := hf.indexMap[name]
	return vecInd, nil
}

// Remove is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It removes the VectorIndex corresponding with a string name using the in memory map.
func (hf *IndexFactory[T]) Remove(name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.removeWithLock(name)
}

func (hf *IndexFactory[T]) removeWithLock(name string) error {
	delete(hf.indexMap, name)
	return nil
}

// CreateOrReplace is an implementation of the IndexFactory interface funciton,
// invoked by an persistentIndexFactory. It checks if a VectorIndex
// correpsonding with name exists. If it does, it removes it, and replaces it
// via the Create function using the passed VectorSource. If the VectorIndex
// does not exist, it creates that VectorIndex corresponding with the name using
// the VectorSource.
func (hf *IndexFactory[T]) CreateOrReplace(
	name string,
	o options.Options,
	source index.VectorSource[T],
	index index.VectorIndex[T]) (index.VectorIndex[T], error) {
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
	return hf.createWithLock(name, o, source, index)
}

func (hf *IndexFactory[T]) WriteToWAL() error {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	operation := func() error {
		if err := encoder.Encode(hf.indexMap); err != nil {
			return fmt.Errorf("could not encode file content, %s", err)
		}

		// write using storage.WriteFile
		if err := storage.WriteFile(context.Background(), "index.wal", buf.Bytes()); err != nil {
			return fmt.Errorf("could not write to file, %s", err)
		}
		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)
}

func (hf *IndexFactory[T]) ReadFromWAL() error {
	// read using storage.ReadFile

	operation := func() error {
		data, err := storage.GetFileContents(context.Background(), "index.wal")
		if err != nil {
			return fmt.Errorf("could not get file content, %s", err)
		}

		decoder := gob.NewDecoder(bytes.NewReader(data))

		if err := decoder.Decode(&hf.indexMap); err != nil {
			return fmt.Errorf("could not decode file content, %s", err)
		}

		return nil
	}

	exponentialBackoff := backoff.NewExponentialBackOff()
	exponentialBackoff.MaxElapsedTime = 10 * time.Second

	return backoff.Retry(operation, exponentialBackoff)

}

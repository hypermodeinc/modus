package in_mem

import (
	"errors"
	"hmruntime/vector/index"
	"sync"

	"hmruntime/vector/options"
)

type InMemIndexFactory struct {
	indexMap map[string]index.VectorIndex[float64]
	mu       sync.RWMutex
}

// CreateFactory creates an instance of the private struct InMemIndexFactory.
// NOTE: if T and floatBits do not match in # of bits, there will be consequences.
func CreateFactory() *InMemIndexFactory {
	f := &InMemIndexFactory{
		indexMap: map[string]index.VectorIndex[float64]{},
	}
	return f
}

// Implements NamedFactory interface for use as a plugin.
func (hf *InMemIndexFactory) Name() string { return "in_mem" }

func (hf *InMemIndexFactory) isNameAvailableWithLock(name string) bool {
	_, nameUsed := hf.indexMap[name]
	return !nameUsed
}

// hf.AllowedOptions() allows persistentIndexFactory to implement the
// IndexFactory interface (see vector-indexer/index/index.go for details).
// We define here options for exponent, maxLevels, efSearch, efConstruction,
// and metric.
func (hf *InMemIndexFactory) AllowedOptions() options.AllowedOptions {
	retVal := options.NewAllowedOptions()

	return retVal
}

// Create is an implementation of the IndexFactory interface function, invoked by an HNSWIndexFactory
// instance. It takes in a string name and a VectorSource implementation, and returns a VectorIndex and error
// flag. It creates an HNSW instance using the index name and populates other parts of the HNSW struct such as
// multFactor, maxLevels, efConstruction, maxNeighbors, and efSearch using struct parameters.
// It then populates the HNSW graphs using the InsertChunk function until there are no more items to populate.
// Finally, the function adds the name and hnsw object to the in memory map and returns the object.
func (hf *InMemIndexFactory) Create(
	name string,
	o options.Options,
	source index.VectorSource[float64],
	floatBits int) (index.VectorIndex[float64], error) {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.createWithLock(name, o, source, floatBits)
}

func (hf *InMemIndexFactory) createWithLock(
	name string,
	o options.Options,
	source index.VectorSource[float64],
	floatBits int) (index.VectorIndex[float64], error) {
	if !hf.isNameAvailableWithLock(name) {
		err := errors.New("index with name " + name + " already exists")
		return nil, err
	}
	retVal := &InMemBruteForceIndex{
		pred:        name,
		vectorNodes: map[uint64][]float64{},
	}
	// while NextChunk() returns a chunk of vectors and their corresponding uids, insert them into the index graph
	for chunk, uids, more, err := source.NextChunk(); more; {
		if err != nil {
			return nil, err
		}
		for i, vector := range chunk {
			retVal.vectorNodes[uids[i]] = vector
		}
	}
	err := retVal.applyOptions(o)
	if err != nil {
		return nil, err
	}
	hf.indexMap[name] = retVal
	return retVal, nil
}

// Find is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It returns the VectorIndex corresponding with a string name using the in memory map.
func (hf *InMemIndexFactory) Find(name string) (index.VectorIndex[float64], error) {
	hf.mu.RLock()
	defer hf.mu.RUnlock()
	return hf.findWithLock(name)
}

func (hf *InMemIndexFactory) findWithLock(name string) (index.VectorIndex[float64], error) {
	vecInd := hf.indexMap[name]
	return vecInd, nil
}

// Remove is an implementation of the IndexFactory interface function, invoked by an persistentIndexFactory
// instance. It removes the VectorIndex corresponding with a string name using the in memory map.
func (hf *InMemIndexFactory) Remove(name string) error {
	hf.mu.Lock()
	defer hf.mu.Unlock()
	return hf.removeWithLock(name)
}

func (hf *InMemIndexFactory) removeWithLock(name string) error {
	delete(hf.indexMap, name)
	return nil
}

// CreateOrReplace is an implementation of the IndexFactory interface funciton,
// invoked by an persistentIndexFactory. It checks if a VectorIndex
// correpsonding with name exists. If it does, it removes it, and replaces it
// via the Create function using the passed VectorSource. If the VectorIndex
// does not exist, it creates that VectorIndex corresponding with the name using
// the VectorSource.
func (hf *InMemIndexFactory) CreateOrReplace(
	name string,
	o options.Options,
	source index.VectorSource[float64],
	floatBits int) (index.VectorIndex[float64], error) {
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
	return hf.createWithLock(name, o, source, floatBits)
}

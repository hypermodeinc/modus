package sequential

import (
	"context"
	"fmt"
	"hmruntime/collections/utils"
	"sync"
	"testing"
)

func TestMultipleSequentialVectorIndexes(t *testing.T) {
	ctx := context.Background()

	// Define the base data to be inserted
	baseTextIds := []int64{1, 2, 3}
	baseKeys := []string{"key1", "key2", "key3"}
	baseVecs := [][]float32{
		{0.1, 0.2, 0.3},
		{0.4, 0.5, 0.6},
		{0.7, 0.8, 0.9},
	}

	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup

	// Define the number of indexes to create
	numIndexes := 20

	// Create and initialize the indexes
	for i := 11; i < 11+numIndexes; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			// Create a new SequentialVectorIndex
			index := NewSequentialVectorIndex("collection"+fmt.Sprint(i), "searchMethod"+fmt.Sprint(i), "embedder"+fmt.Sprint(i))

			// Generate unique data for this index
			textIds := make([]int64, len(baseTextIds))
			keys := make([]string, len(baseKeys))
			vecs := make([][]float32, len(baseVecs))
			for j := range baseTextIds {
				textIds[j] = baseTextIds[j] + int64(i*len(baseTextIds))
				keys[j] = baseKeys[j] + fmt.Sprint(i)
				vecs[j] = append([]float32{}, baseVecs[j]...)
				for k := range vecs[j] {
					vecs[j][k] += float32(i) / 10
				}
				var err error
				vecs[j], err = utils.Normalize(vecs[j])
				if err != nil {
					t.Errorf("Failed to normalize vector: %v", err)
				}
			}

			err := index.InsertVectorsToMemory(ctx, textIds, textIds, keys, vecs)
			if err != nil {
				t.Errorf("Failed to insert vectors into index: %v", err)
			}

			// Verify the vectors were inserted correctly
			for _, key := range keys {
				expectedVec, err := index.GetVector(ctx, key)
				if err != nil {
					t.Errorf("index %d: Failed to get expected vector from index: %v", i, err)
				}
				objs, err := index.SearchWithKey(ctx, key, 1, nil)
				if err != nil {
					t.Errorf("index %d: Failed to search vector in index: %v", i, err)
				}
				if len(objs) == 0 {
					t.Errorf("index %d: Expected obj with length 1, got %v", i, len(objs))
				}
				resVec, err := index.GetVector(ctx, objs[0].GetIndex())
				if err != nil {
					t.Errorf("index %d: Failed to get result vector from index: %v", i, err)
				}

				if !utils.EqualFloat32Slices(expectedVec, resVec) {
					t.Errorf("index %d: Expected vector %v, got %v", i, expectedVec, resVec)
				}

				checkpointId, err := index.GetCheckpointId(ctx)
				if err != nil {
					t.Errorf("Failed to get checkpoint ID: %v", err)
				}
				if checkpointId != textIds[len(textIds)-1] {
					t.Errorf("Expected checkpoint ID %v, got %v", textIds[len(textIds)-1], checkpointId)
				}

				lastIndexedTextID, err := index.GetLastIndexedTextId(ctx)
				if err != nil {
					t.Errorf("Failed to get last indexed text ID: %v", err)
				}
				if lastIndexedTextID != textIds[len(textIds)-1] {
					t.Errorf("Expected last indexed text ID %v, got %v", textIds[len(textIds)-1], lastIndexedTextID)
				}
			}

		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

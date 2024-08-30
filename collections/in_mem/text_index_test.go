/*
 * Copyright 2024 Hypermode, Inc.
 */

package in_mem

import (
	"context"
	"fmt"
	"sync"
	"testing"
)

func TestMultipleInMemCollections(t *testing.T) {
	ctx := context.Background()

	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup

	// Define the number of collections to create
	numCollections := 10

	// Create and initialize the collections
	for i := 0; i < numCollections; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			// Define the data to be inserted
			ids := []int64{int64(i*3 + 1), int64(i*3 + 2), int64(i*3 + 3)}
			keys := []string{fmt.Sprintf("key%d_1", i), fmt.Sprintf("key%d_2", i), fmt.Sprintf("key%d_3", i)}
			texts := []string{fmt.Sprintf("text%d_1", i), fmt.Sprintf("text%d_2", i), fmt.Sprintf("text%d_3", i)}
			labels := [][]string{
				{fmt.Sprintf("label%d_1", i)},
				{fmt.Sprintf("label%d_2", i)},
				{fmt.Sprintf("label%d_3", i)},
			}

			// Create a new InMemCollection
			col := NewCollectionNamespace("collection"+fmt.Sprint(i), "")

			// Insert the texts into the collection
			err := col.InsertTextsToMemory(ctx, ids, keys, texts, labels)
			if err != nil {
				t.Errorf("Failed to insert texts into collection: %v", err)
			}

			// Verify the texts were inserted correctly
			for j, key := range keys {
				text, err := col.GetText(ctx, key)
				if err != nil {
					t.Errorf("Failed to get text from collection: %v", err)
				}
				if text != texts[j] {
					t.Errorf("Expected text %s, got %s", texts[j], text)
				}

				// Verify the external ID
				extID, err := col.GetExternalId(ctx, key)
				if err != nil {
					t.Errorf("Failed to get external ID from collection: %v", err)
				}
				if extID != ids[j] {
					t.Errorf("Expected external ID %d, got %d", ids[j], extID)
				}
			}

			// Verify the checkpoint ID
			chkID, err := col.GetCheckpointId(ctx)
			if err != nil {
				t.Errorf("Failed to get checkpoint ID from collection: %v", err)
			}
			if chkID != ids[len(ids)-1] {
				t.Errorf("Expected checkpoint ID %d, got %d", ids[len(ids)-1], chkID)
			}
		}(i)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sequential

import (
	"container/heap"
	"context"
	"fmt"
	"sync"

	"github.com/hypermodeinc/modus/runtime/collections/index"
	"github.com/hypermodeinc/modus/runtime/collections/utils"
	"github.com/hypermodeinc/modus/runtime/db"
)

const (
	SequentialVectorIndexType = "SequentialVectorIndex"
)

type SequentialVectorIndex struct {
	mu                sync.RWMutex
	searchMethodName  string
	embedderName      string
	lastInsertedID    int64
	lastIndexedTextID int64
	VectorMap         map[string][]float32 // key: vector
}

func NewSequentialVectorIndex(searchMethod, embedder string) *SequentialVectorIndex {
	return &SequentialVectorIndex{
		searchMethodName: searchMethod,
		embedderName:     embedder,
		VectorMap:        make(map[string][]float32),
	}
}

func (ims *SequentialVectorIndex) GetSearchMethodName() string {
	return ims.searchMethodName
}

func (ims *SequentialVectorIndex) SetEmbedderName(embedderName string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.embedderName = embedderName
	return nil
}

func (ims *SequentialVectorIndex) GetEmbedderName() string {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.embedderName
}

func (ims *SequentialVectorIndex) GetVectorNodesMap() map[string][]float32 {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorMap
}

func (ims *SequentialVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MaxTupleHeap, error) {
	// calculate cosine similarity and return top maxResults results
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	if maxResults <= 0 {
		maxResults = 1
	}
	var results utils.MaxTupleHeap
	heap.Init(&results)
	for key, vector := range ims.VectorMap {
		if filter != nil && !filter(query, vector, key) {
			continue
		}
		similarity, err := utils.CosineDistance(query, vector)
		if err != nil {
			return nil, err
		}
		if results.Len() < maxResults {
			heap.Push(&results, utils.InitHeapElement(similarity, key, false))
		} else if utils.IsBetterScoreForDistance(similarity, results[0].GetValue()) {
			heap.Pop(&results)
			heap.Push(&results, utils.InitHeapElement(similarity, key, false))
		}
	}

	// Return top maxResults results
	var finalResults utils.MaxTupleHeap
	for results.Len() > 0 {
		finalResults = append(finalResults, heap.Pop(&results).(utils.MaxHeapElement))
	}
	// Reverse the finalResults to get the highest similarity first
	for i, j := 0, len(finalResults)-1; i < j; i, j = i+1, j-1 {
		finalResults[i], finalResults[j] = finalResults[j], finalResults[i]
	}
	return finalResults, nil
}

func (ims *SequentialVectorIndex) SearchWithKey(ctx context.Context, queryKey string, maxResults int, filter index.SearchFilter) (utils.MaxTupleHeap, error) {
	ims.mu.RLock()
	query := ims.VectorMap[queryKey]
	ims.mu.RUnlock()
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, query, maxResults, filter)
}

func (ims *SequentialVectorIndex) InsertVectors(ctx context.Context, textIds []int64, vecs [][]float32) error {
	if len(textIds) != len(vecs) {
		return fmt.Errorf("textIds and vecs must have the same length")
	}
	vectorIds, keys, err := db.WriteCollectionVectors(ctx, ims.searchMethodName, textIds, vecs)
	if err != nil {
		return err
	}

	return ims.InsertVectorsToMemory(ctx, textIds, vectorIds, keys, vecs)
}

func (ims *SequentialVectorIndex) InsertVector(ctx context.Context, textId int64, vec []float32) error {

	// Write vector to database, this textId is now the last inserted textId
	vectorId, key, err := db.WriteCollectionVector(ctx, ims.searchMethodName, textId, vec)
	if err != nil {
		return err
	}

	return ims.InsertVectorToMemory(ctx, textId, vectorId, key, vec)

}

func (ims *SequentialVectorIndex) InsertVectorsToMemory(ctx context.Context, textIds []int64, vectorIds []int64, keys []string, vecs [][]float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	for i, key := range keys {
		ims.VectorMap[key] = vecs[i]
		ims.lastInsertedID = vectorIds[i]
		ims.lastIndexedTextID = textIds[i]
	}
	return nil
}

func (ims *SequentialVectorIndex) InsertVectorToMemory(ctx context.Context, textId, vectorId int64, key string, vec []float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.VectorMap[key] = vec
	ims.lastInsertedID = vectorId
	ims.lastIndexedTextID = textId
	return nil
}

func (ims *SequentialVectorIndex) DeleteVector(ctx context.Context, textId int64, key string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	err := db.DeleteCollectionVector(ctx, ims.searchMethodName, textId)
	if err != nil {
		return err
	}
	delete(ims.VectorMap, key)
	return nil
}

func (ims *SequentialVectorIndex) GetVector(ctx context.Context, key string) ([]float32, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.VectorMap[key], nil
}

func (ims *SequentialVectorIndex) GetCheckpointId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastInsertedID, nil
}

func (ims *SequentialVectorIndex) GetLastIndexedTextId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastIndexedTextID, nil
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hnsw

import (
	"container/heap"
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/hypermodeinc/modus/runtime/collections/index"
	"github.com/hypermodeinc/modus/runtime/collections/utils"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/hnsw"
)

const (
	HnswVectorIndexType = "HnswVectorIndex"
)

type HnswVectorIndex struct {
	mu                sync.RWMutex
	searchMethodName  string
	embedderName      string
	lastInsertedID    int64
	lastIndexedTextID int64
	HnswIndex         *hnsw.Graph[string]
}

func NewHnswVectorIndex(searchMethod, embedder string) *HnswVectorIndex {
	return &HnswVectorIndex{
		searchMethodName: searchMethod,
		embedderName:     embedder,
		HnswIndex:        hnsw.NewGraph[string](),
	}
}

func (ims *HnswVectorIndex) GetSearchMethodName() string {
	return ims.searchMethodName
}

func (ims *HnswVectorIndex) SetEmbedderName(embedderName string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	ims.embedderName = embedderName
	return nil
}

func (ims *HnswVectorIndex) GetEmbedderName() string {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.embedderName
}

func (ims *HnswVectorIndex) GetVectorNodesMap() map[string][]float32 {
	return map[string][]float32{}
}

func (ims *HnswVectorIndex) Search(ctx context.Context, query []float32, maxResults int, filter index.SearchFilter) (utils.MaxTupleHeap, error) {
	// calculate cosine similarity and return top maxResults results
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	if ims.HnswIndex == nil {
		return nil, fmt.Errorf("vector index is not initialized")
	}
	neighbors, err := ims.HnswIndex.Search(query, maxResults)
	if err != nil {
		return nil, err
	}
	keys, distances := make([]string, len(neighbors)), make([]float64, len(neighbors))
	for i, neighbor := range neighbors {
		keys[i] = string(neighbor.Key)
		distances[i] = float64(neighbor.Distance)
	}
	var results utils.MaxTupleHeap
	heap.Init(&results)

	for i, key := range keys {
		heap.Push(&results, utils.InitHeapElement(distances[i], key, false))
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

func (ims *HnswVectorIndex) SearchWithKey(ctx context.Context, queryKey string, maxResults int, filter index.SearchFilter) (utils.MaxTupleHeap, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	query, found := ims.HnswIndex.Lookup(queryKey)
	if !found {
		return nil, fmt.Errorf("key not found")
	}
	if query == nil {
		return nil, nil
	}
	return ims.Search(ctx, query, maxResults, filter)
}

func (ims *HnswVectorIndex) InsertVectors(ctx context.Context, textIds []int64, vecs [][]float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	if len(textIds) != len(vecs) {
		return fmt.Errorf("textIds and vecs must have the same length")
	}
	vectorIds, keys, err := db.WriteCollectionVectors(ctx, ims.searchMethodName, textIds, vecs)
	if err != nil {
		return err
	}

	return ims.InsertVectorsToMemory(ctx, textIds, vectorIds, keys, vecs)
}

func (ims *HnswVectorIndex) InsertVector(ctx context.Context, textId int64, vec []float32) error {

	// Write vector to database, this textId is now the last inserted textId
	vectorId, key, err := db.WriteCollectionVector(ctx, ims.searchMethodName, textId, vec)
	if err != nil {
		return err
	}

	return ims.InsertVectorToMemory(ctx, textId, vectorId, key, vec)

}

func (ims *HnswVectorIndex) InsertVectorsToMemory(ctx context.Context, textIds []int64, vectorIds []int64, keys []string, vecs [][]float32) error {
	nodes, err := hnsw.MakeNodes(keys, vecs)
	if err != nil {
		return err
	}
	err = ims.HnswIndex.Add(nodes...)
	if err != nil {
		return err
	}
	ims.lastInsertedID = slices.Max(vectorIds)
	ims.lastIndexedTextID = slices.Max(textIds)
	return nil
}

func (ims *HnswVectorIndex) InsertVectorToMemory(ctx context.Context, textId, vectorId int64, key string, vec []float32) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	err := ims.HnswIndex.Add(hnsw.MakeNode(key, vec))
	if err != nil {
		return err
	}
	ims.lastInsertedID = vectorId
	ims.lastIndexedTextID = textId
	return nil
}

func (ims *HnswVectorIndex) DeleteVector(ctx context.Context, textId int64, key string) error {
	ims.mu.Lock()
	defer ims.mu.Unlock()
	err := db.DeleteCollectionVector(ctx, ims.searchMethodName, textId)
	if err != nil {
		return err
	}
	ims.HnswIndex.Delete(key)
	return nil
}

func (ims *HnswVectorIndex) GetVector(ctx context.Context, key string) ([]float32, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	vec, _ := ims.HnswIndex.Lookup(key)
	return vec, nil
}

func (ims *HnswVectorIndex) GetCheckpointId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastInsertedID, nil
}

func (ims *HnswVectorIndex) GetLastIndexedTextId(ctx context.Context) (int64, error) {
	ims.mu.RLock()
	defer ims.mu.RUnlock()
	return ims.lastIndexedTextID, nil
}

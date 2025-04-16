/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/puzpuzpuz/xsync/v4"
)

const collectionFactoryWriteInterval = 1

var (
	globalNamespaceManager *collectionFactory
	errCollectionNotFound  = fmt.Errorf("collection not found")
	errNamespaceNotFound   = fmt.Errorf("namespace not found")
)

type collectionFactory struct {
	collectionMap map[string]*collection
	mu            sync.RWMutex
	quit          chan struct{}
	done          chan struct{}
}

func newCollectionFactory() *collectionFactory {
	return &collectionFactory{
		collectionMap: map[string]*collection{
			"": {
				collectionNamespaceMap: xsync.NewMapOf[string, interfaces.CollectionNamespace](),
			},
		},
		quit: make(chan struct{}),
		done: make(chan struct{}),
	}
}

func (cf *collectionFactory) createCollection(name string, coll *collection) (*collection, error) {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	if _, found := cf.collectionMap[name]; found {
		return nil, fmt.Errorf("collection with name %s already exists", name)
	}
	cf.collectionMap[name] = coll
	return coll, nil
}

func (cf *collectionFactory) findCollection(name string) (*collection, error) {
	cf.mu.RLock()
	defer cf.mu.RUnlock()

	coll, ok := cf.collectionMap[name]
	if !ok {
		return nil, errCollectionNotFound
	}
	return coll, nil
}

func (cf *collectionFactory) getNamespaceCollectionFactoryMap() map[string]*collection {
	return cf.collectionMap
}

func (cf *collectionFactory) removeCollection(name string) error {
	cf.mu.Lock()
	defer cf.mu.Unlock()

	delete(cf.collectionMap, name)
	return nil
}

func (cf *collectionFactory) readFromPostgres(ctx context.Context) bool {
	resetTimerFaster := false
	var err error
	for _, namespaceCollectionFactory := range cf.collectionMap {
		for _, col := range namespaceCollectionFactory.getCollectionNamespaceMap() {
			resetTimerFaster, err = loadTextsIntoCollection(ctx, col)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", col.GetCollectionName()).
					Msg("Failed to load texts into collection.")
				break
			}

			for _, vectorIndex := range col.GetVectorIndexMap() {
				err = loadVectorsIntoVectorIndex(ctx, vectorIndex, col)
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", col.GetCollectionName()).
						Str("search_method", vectorIndex.GetSearchMethodName()).
						Msg("Failed to load vectors into vector index.")
					break
				}

				// catch up on any texts that weren't embedded
				err := syncTextsWithVectorIndex(ctx, col, vectorIndex)
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", col.GetCollectionName()).
						Str("search_method", vectorIndex.GetSearchMethodName()).
						Msg("Failed to sync text with vector index.")
					break
				}

			}
		}
	}
	return resetTimerFaster
}

func (cf *collectionFactory) worker(ctx context.Context) {
	defer close(cf.done)
	timer := time.NewTimer(collectionFactoryWriteInterval * time.Minute)

	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			// read from postgres all collections & searchMethod after lastInsertedID
			resetTimerFaster := cf.readFromPostgres(ctx)
			if resetTimerFaster {
				timer.Reset(10 * time.Second)
			} else {
				timer.Reset(collectionFactoryWriteInterval * time.Minute)
			}
		case <-cf.quit:
			return
		}
	}
}

func loadTextsIntoCollection(ctx context.Context, col interfaces.CollectionNamespace) (resetTimerFaster bool, err error) {
	// Get checkpoint id for collection
	textCheckpointId, err := col.GetCheckpointId(ctx)
	if err != nil {
		return false, err
	}

	// Query all texts from checkpoint
	textIds, keys, texts, labels, err := db.QueryCollectionTextsFromCheckpoint(ctx, col.GetCollectionName(), col.GetNamespace(), textCheckpointId)
	if err != nil {
		return false, err
	}
	if len(textIds) != len(keys) || len(keys) != len(texts) {
		return false, errors.New("mismatch in keys and texts")
	}

	if (textCheckpointId == 0) && (len(textIds) == 0) {
		return true, nil
	}

	// Insert all texts into collection
	err = col.InsertTextsToMemory(ctx, textIds, keys, texts, labels)
	if err != nil {
		return false, err
	}

	return false, nil
}

func loadVectorsIntoVectorIndex(ctx context.Context, vectorIndex interfaces.VectorIndex, col interfaces.CollectionNamespace) error {
	// Get checkpoint id for vector index
	vecCheckpointId, err := vectorIndex.GetCheckpointId(ctx)
	if err != nil {
		return err
	}

	// Query all vectors from checkpoint
	textIds, vectorIds, keys, vectors, err := db.QueryCollectionVectorsFromCheckpoint(ctx, col.GetCollectionName(), vectorIndex.GetSearchMethodName(), col.GetNamespace(), vecCheckpointId)
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

func syncTextsWithVectorIndex(ctx context.Context, col interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex) error {
	// Query all texts from checkpoint
	lastIndexedTextId, err := vectorIndex.GetLastIndexedTextId(ctx)
	if err != nil {
		return err
	}
	textIds, keys, texts, _, err := db.QueryCollectionTextsFromCheckpoint(ctx, col.GetCollectionName(), col.GetNamespace(), lastIndexedTextId)
	if err != nil {
		return err
	}
	if len(textIds) != len(keys) || len(keys) != len(texts) {
		return errors.New("mismatch in keys and texts")
	}
	// Get last indexed text id
	err = processTexts(ctx, col, vectorIndex, keys, texts)
	if err != nil {
		return err
	}

	return nil
}

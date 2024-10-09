/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package in_mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/hypermodeinc/modus/runtime/collections/index"
	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
	"github.com/hypermodeinc/modus/runtime/db"
)

const DefaultNamespace = ""

type InMemCollectionNamespace struct {
	mu             sync.RWMutex
	collectionName string
	namespace      string
	lastInsertedID int64
	TextMap        map[string]string // key: text
	LabelsMap      map[string][]string
	IdMap          map[string]int64                          // key: postgres id
	VectorIndexMap map[string]*interfaces.VectorIndexWrapper // searchMethod: vectorIndex
}

func NewCollectionNamespace(name, namespace string) *InMemCollectionNamespace {
	return &InMemCollectionNamespace{
		collectionName: name,
		namespace:      namespace,
		TextMap:        map[string]string{},
		LabelsMap:      map[string][]string{},
		IdMap:          map[string]int64{},
		VectorIndexMap: map[string]*interfaces.VectorIndexWrapper{},
	}
}

func (ti *InMemCollectionNamespace) GetCollectionName() string {
	return ti.collectionName
}

func (ti *InMemCollectionNamespace) GetNamespace() string {
	return ti.namespace
}

func (ti *InMemCollectionNamespace) GetVectorIndexMap() map[string]*interfaces.VectorIndexWrapper {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.VectorIndexMap
}

func (ti *InMemCollectionNamespace) GetVectorIndex(ctx context.Context, searchMethod string) (*interfaces.VectorIndexWrapper, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	if ind, ok := ti.VectorIndexMap[searchMethod]; !ok {
		return nil, index.ErrVectorIndexNotFound
	} else {
		return ind, nil
	}
}

func (ti *InMemCollectionNamespace) SetVectorIndex(ctx context.Context, searchMethod string, vectorIndex *interfaces.VectorIndexWrapper) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	if ti.VectorIndexMap == nil {
		ti.VectorIndexMap = map[string]*interfaces.VectorIndexWrapper{}
	}
	if _, ok := ti.VectorIndexMap[searchMethod]; ok {
		return index.ErrVectorIndexAlreadyExists
	}
	ti.VectorIndexMap[searchMethod] = vectorIndex
	return nil
}

func (ti *InMemCollectionNamespace) DeleteVectorIndex(ctx context.Context, searchMethod string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	delete(ti.VectorIndexMap, searchMethod)
	return nil
}

func (ti *InMemCollectionNamespace) InsertTexts(ctx context.Context, keys []string, texts []string, labelsArr [][]string) error {
	if len(keys) != len(texts) {
		return fmt.Errorf("keys and texts must have the same length")
	}

	if len(labelsArr) != 0 && len(labelsArr) != len(keys) {
		return fmt.Errorf("labels must have the same length as keys or be empty")
	}

	ids, err := db.WriteCollectionTexts(ctx, ti.collectionName, ti.namespace, keys, texts, labelsArr)
	if err != nil {
		return err
	}

	return ti.InsertTextsToMemory(ctx, ids, keys, texts, labelsArr)
}

func (ti *InMemCollectionNamespace) InsertText(ctx context.Context, key string, text string, labels []string) error {
	id, err := db.WriteCollectionText(ctx, ti.collectionName, ti.namespace, key, text, labels)
	if err != nil {
		return err
	}

	return ti.InsertTextToMemory(ctx, id, key, text, labels)
}

func (ti *InMemCollectionNamespace) InsertTextsToMemory(ctx context.Context, ids []int64, keys []string, texts []string, labelsArr [][]string) error {

	if len(labelsArr) != 0 && len(labelsArr) != len(keys) {
		return fmt.Errorf("labels must have the same length as keys or be empty")
	}
	if len(ids) != len(keys) || len(ids) != len(texts) {
		return fmt.Errorf("ids, keys and texts must have the same length")
	}

	ti.mu.Lock()
	defer ti.mu.Unlock()
	for i, key := range keys {
		ti.TextMap[key] = texts[i]
		if len(labelsArr) != 0 {
			ti.LabelsMap[key] = labelsArr[i]
		}
		ti.IdMap[key] = ids[i]
		ti.lastInsertedID = ids[i]
	}
	return nil
}

func (ti *InMemCollectionNamespace) InsertTextToMemory(ctx context.Context, id int64, key string, text string, labels []string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	ti.TextMap[key] = text
	if len(labels) != 0 {
		ti.LabelsMap[key] = labels
	}
	ti.IdMap[key] = id
	ti.lastInsertedID = id
	return nil
}

func (ti *InMemCollectionNamespace) DeleteText(ctx context.Context, key string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()
	err := db.DeleteCollectionText(ctx, ti.collectionName, ti.namespace, key)
	if err != nil {
		return err
	}
	delete(ti.TextMap, key)
	return nil
}

func (ti *InMemCollectionNamespace) GetText(ctx context.Context, key string) (string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap[key], nil
}

func (ti *InMemCollectionNamespace) GetTextMap(ctx context.Context) (map[string]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.TextMap, nil
}

func (ti *InMemCollectionNamespace) GetLabels(ctx context.Context, key string) ([]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.LabelsMap[key], nil
}

func (ti *InMemCollectionNamespace) GetLabelsMap(ctx context.Context) (map[string][]string, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.LabelsMap, nil
}

func (ti *InMemCollectionNamespace) Len(ctx context.Context) (int, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return len(ti.TextMap), nil
}

func (ti *InMemCollectionNamespace) GetExternalId(ctx context.Context, key string) (int64, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.IdMap[key], nil
}

func (ti *InMemCollectionNamespace) GetCheckpointId(ctx context.Context) (int64, error) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()
	return ti.lastInsertedID, nil
}

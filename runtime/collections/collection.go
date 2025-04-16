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
	"fmt"

	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
	"github.com/puzpuzpuz/xsync/v4"
)

type collection struct {
	collectionNamespaceMap *xsync.Map[string, interfaces.CollectionNamespace]
}

func newCollection() *collection {
	return &collection{
		collectionNamespaceMap: xsync.NewMap[string, interfaces.CollectionNamespace](),
	}
}

func (c *collection) getCollectionNamespaceMap() map[string]interfaces.CollectionNamespace {
	m := make(map[string]interfaces.CollectionNamespace, c.collectionNamespaceMap.Size())
	c.collectionNamespaceMap.Range(func(key string, value interfaces.CollectionNamespace) bool {
		m[key] = value
		return true
	})

	return m
}

func (c *collection) findNamespace(namespace string) (interfaces.CollectionNamespace, error) {
	ns, found := c.collectionNamespaceMap.Load(namespace)
	if !found {
		return nil, errNamespaceNotFound
	}
	return ns, nil
}

func (c *collection) findOrCreateNamespace(namespace string, index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	result, _ := c.collectionNamespaceMap.LoadOrStore(namespace, index)
	return result, nil // TODO: remove unused error
}

func (c *collection) createCollectionNamespace(namespace string, index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	_, found := c.collectionNamespaceMap.Load(namespace)
	if found {
		return nil, fmt.Errorf("namespace with name %s already exists", namespace)
	}

	c.collectionNamespaceMap.Store(namespace, index)
	return index, nil
}

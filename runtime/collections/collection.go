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
	"sync"

	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
)

type collection struct {
	collectionNamespaceMap map[string]interfaces.CollectionNamespace
	mu                     sync.RWMutex
}

func newCollection() *collection {
	return &collection{
		collectionNamespaceMap: map[string]interfaces.CollectionNamespace{},
	}
}

func (c *collection) getCollectionNamespaceMap() map[string]interfaces.CollectionNamespace {
	return c.collectionNamespaceMap
}

func (c *collection) findNamespace(namespace string) (interfaces.CollectionNamespace, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	ns, found := c.collectionNamespaceMap[namespace]
	if !found {
		return nil, errNamespaceNotFound
	}
	return ns, nil
}

func (c *collection) findOrCreateNamespace(namespace string, index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	c.mu.RLock()
	ns, found := c.collectionNamespaceMap[namespace]
	if found {
		defer c.mu.RUnlock()
		return ns, nil
	}

	c.mu.RUnlock()
	c.mu.Lock()
	defer c.mu.Unlock()

	ns, found = c.collectionNamespaceMap[namespace]
	if found {
		return ns, nil
	}

	c.collectionNamespaceMap[namespace] = index
	return index, nil
}

func (c *collection) createCollectionNamespace(namespace string, index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, found := c.collectionNamespaceMap[namespace]; found {
		return nil, fmt.Errorf("namespace with name %s already exists", namespace)
	}

	c.collectionNamespaceMap[namespace] = index
	return index, nil
}

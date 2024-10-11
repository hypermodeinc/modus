/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/collections"
)

func init() {
	const module_name = "modus_collections"

	registerHostFunction(module_name, "computeDistance", collections.ComputeDistance,
		withCancelledMessage("Cancelled computing distance."),
		withErrorMessage("Error computing distance."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction(module_name, "delete", collections.Delete,
		withCancelledMessage("Cancelled deleting from collection."),
		withErrorMessage("Error deleting from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction(module_name, "getNamespaces", collections.GetNamespaces,
		withCancelledMessage("Cancelled getting namespaces from collection."),
		withErrorMessage("Error getting namespaces from collection."),
		withMessageDetail(func(collectionName string) string {
			return fmt.Sprintf("Collection: %s", collectionName)
		}))

	registerHostFunction(module_name, "getText", collections.GetText,
		withCancelledMessage("Cancelled getting text from collection."),
		withErrorMessage("Error getting text from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction(module_name, "dumpTexts", collections.DumpTexts,
		withCancelledMessage("Cancelled dumping texts from collection."),
		withErrorMessage("Error dumping texts from collection."),
		withMessageDetail(func(collectionName, namespace string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s", collectionName, namespace)
		}))

	registerHostFunction(module_name, "getVector", collections.GetVector,
		withCancelledMessage("Cancelled getting vector from collection."),
		withErrorMessage("Error getting vector from collection."),
		withMessageDetail(func(collectionName, namespace, searchMethod, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s, Key: %s", collectionName, namespace, searchMethod, key)
		}))

	registerHostFunction(module_name, "getLabels", collections.GetLabels,
		withCancelledMessage("Cancelled getting labels from collection."),
		withErrorMessage("Error getting labels from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction(module_name, "classifyText", collections.ClassifyText,
		withCancelledMessage("Cancelled classifying text."),
		withErrorMessage("Error while classifying text."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction(module_name, "recomputeIndex", collections.RecomputeIndex,
		withStartingMessage("Starting recomputing index for collection."),
		withCompletedMessage("Completed recomputing index for collection."),
		withCancelledMessage("Cancelled recomputing index for collection."),
		withErrorMessage("Error recomputing index for collection."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction(module_name, "search", collections.Search,
		withCancelledMessage("Cancelled searching collection."),
		withErrorMessage("Error searching collection."),
		withMessageDetail(func(collectionName string, namespaces []string, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespaces: %v, Method: %s", collectionName, namespaces, searchMethod)
		}))

	registerHostFunction(module_name, "searchByVector", collections.SearchByVector,
		withCancelledMessage("Cancelled searching collection by vector."),
		withErrorMessage("Error searching collection by vector."),
		withMessageDetail(func(collectionName string, namespaces []string, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespaces: %v, Method: %s", collectionName, namespaces, searchMethod)
		}))

	registerHostFunction(module_name, "upsert", collections.Upsert,
		withCancelledMessage("Cancelled upserting to collection."),
		withErrorMessage("Error upserting to collection."),
		withMessageDetail(func(collectionName, namespace string, keys []string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Keys: %v", collectionName, namespace, keys)
		}))
}

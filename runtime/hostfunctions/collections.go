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

	registerHostFunction("hypermode", "computeDistance", collections.ComputeDistance,
		withCancelledMessage("Cancelled computing distance."),
		withErrorMessage("Error computing distance."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "deleteFromCollection", collections.DeleteFromCollection,
		withCancelledMessage("Cancelled deleting from collection."),
		withErrorMessage("Error deleting from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction("hypermode", "getNamespacesFromCollection", collections.GetNamespacesFromCollection,
		withCancelledMessage("Cancelled getting namespaces from collection."),
		withErrorMessage("Error getting namespaces from collection."),
		withMessageDetail(func(collectionName string) string {
			return fmt.Sprintf("Collection: %s", collectionName)
		}))

	registerHostFunction("hypermode", "getTextFromCollection", collections.GetTextFromCollection,
		withCancelledMessage("Cancelled getting text from collection."),
		withErrorMessage("Error getting text from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction("hypermode", "getTextsFromCollection", collections.GetTextsFromCollection,
		withCancelledMessage("Cancelled getting texts from collection."),
		withErrorMessage("Error getting texts from collection."),
		withMessageDetail(func(collectionName, namespace string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s", collectionName, namespace)
		}))

	registerHostFunction("hypermode", "getVector", collections.GetVector,
		withCancelledMessage("Cancelled getting vector from collection."),
		withErrorMessage("Error getting vector from collection."),
		withMessageDetail(func(collectionName, namespace, searchMethod, id string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s, ID: %s", collectionName, namespace, searchMethod, id)
		}))

	registerHostFunction("hypermode", "getLabels", collections.GetLabels,
		withCancelledMessage("Cancelled getting labels from collection."),
		withErrorMessage("Error getting labels from collection."),
		withMessageDetail(func(collectionName, namespace, id string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, ID: %s", collectionName, namespace, id)
		}))

	registerHostFunction("hypermode", "nnClassifyCollection", collections.NnClassify,
		withCancelledMessage("Cancelled classification."),
		withErrorMessage("Error during classification."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "recomputeSearchMethod", collections.RecomputeSearchMethod,
		withStartingMessage("Starting recomputing search method for collection."),
		withCompletedMessage("Completed recomputing search method for collection."),
		withCancelledMessage("Cancelled recomputing search method for collection."),
		withErrorMessage("Error recomputing search method for collection."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "searchCollection", collections.SearchCollection,
		withCancelledMessage("Cancelled searching collection."),
		withErrorMessage("Error searching collection."),
		withMessageDetail(func(collectionName string, namespaces []string, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespaces: %v, Method: %s", collectionName, namespaces, searchMethod)
		}))

	registerHostFunction("hypermode", "searchCollectionByVector", collections.SearchCollectionByVector,
		withCancelledMessage("Cancelled searching collection by vector."),
		withErrorMessage("Error searching collection by vector."),
		withMessageDetail(func(collectionName string, namespaces []string, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespaces: %v, Method: %s", collectionName, namespaces, searchMethod)
		}))

	registerHostFunction("hypermode", "upsertToCollection", collections.UpsertToCollection,
		withCancelledMessage("Cancelled collection upsert."),
		withErrorMessage("Error upserting to collection."),
		withMessageDetail(func(collectionName, namespace string, keys []string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Keys: %v", collectionName, namespace, keys)
		}))
}

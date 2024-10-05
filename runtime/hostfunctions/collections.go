/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"context"
	"fmt"

	"github.com/hypermodeinc/modus/runtime/collections"
)

func init() {
	registerCollectionsHostFunctions()
	registerLegacyCollectionsHostFunctions()
}

func registerCollectionsHostFunctions() {
	registerHostFunction("hypermode", "computeDistance_v2", collections.ComputeDistance,
		withCancelledMessage("Cancelled computing distance."),
		withErrorMessage("Error computing distance."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "deleteFromCollection_v2", collections.DeleteFromCollection,
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

	registerHostFunction("hypermode", "getTextFromCollection_v2", collections.GetTextFromCollection,
		withCancelledMessage("Cancelled getting text from collection."),
		withErrorMessage("Error getting text from collection."),
		withMessageDetail(func(collectionName, namespace, key string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Key: %s", collectionName, namespace, key)
		}))

	registerHostFunction("hypermode", "getTextsFromCollection_v2", collections.GetTextsFromCollection,
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

	registerHostFunction("hypermode", "nnClassifyCollection_v2", collections.NnClassify,
		withCancelledMessage("Cancelled classification."),
		withErrorMessage("Error during classification."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "recomputeSearchMethod_v2", collections.RecomputeSearchMethod,
		withStartingMessage("Starting recomputing search method for collection."),
		withCompletedMessage("Completed recomputing search method for collection."),
		withCancelledMessage("Cancelled recomputing search method for collection."),
		withErrorMessage("Error recomputing search method for collection."),
		withMessageDetail(func(collectionName, namespace, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Method: %s", collectionName, namespace, searchMethod)
		}))

	registerHostFunction("hypermode", "searchCollection_v2", collections.SearchCollection,
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

	registerHostFunction("hypermode", "upsertToCollection_v2", collections.UpsertToCollection,
		withCancelledMessage("Cancelled collection upsert."),
		withErrorMessage("Error upserting to collection."),
		withMessageDetail(func(collectionName, namespace string, keys []string) string {
			return fmt.Sprintf("Collection: %s, Namespace: %s, Keys: %v", collectionName, namespace, keys)
		}))
}

func registerLegacyCollectionsHostFunctions() {

	// Support functions from older SDK versions.
	// Each of these function wrappers must maintain the original signature.
	// We can remove these when we can be sure that nobody is using them.

	registerHostFunction("hypermode", "computeDistance",
		func(ctx context.Context, collectionName, searchMethod, id1, id2 string) (*collections.CollectionSearchResultObject, error) {
			return collections.ComputeDistance(ctx, collectionName, "", searchMethod, id1, id2)
		},
		withCancelledMessage("Cancelled computing distance."),
		withErrorMessage("Error computing distance."),
		withMessageDetail(func(collectionName, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod)
		}))

	registerHostFunction("hypermode", "computeSimilarity",
		func(ctx context.Context, collectionName, searchMethod, id1, id2 string) (*collections.CollectionSearchResultObject, error) {
			return collections.ComputeDistance(ctx, collectionName, "", searchMethod, id1, id2)
		},
		withCancelledMessage("Cancelled computing similarity."),
		withErrorMessage("Error computing similarity."),
		withMessageDetail(func(collectionName, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod)
		}))

	registerHostFunction("hypermode", "deleteFromCollection",
		func(ctx context.Context, collectionName, key string) (*collections.CollectionMutationResult, error) {
			return collections.DeleteFromCollection(ctx, collectionName, "", key)
		},
		withCancelledMessage("Cancelled deleting from collection."),
		withErrorMessage("Error deleting from collection."),
		withMessageDetail(func(collectionName, key string) string {
			return fmt.Sprintf("Collection: %s, Key: %s", collectionName, key)
		}))

	registerHostFunction("hypermode", "getTextFromCollection",
		func(ctx context.Context, collectionName, key string) (string, error) {
			return collections.GetTextFromCollection(ctx, collectionName, "", key)
		},
		withCancelledMessage("Cancelled getting text from collection."),
		withErrorMessage("Error getting text from collection."),
		withMessageDetail(func(collectionName, key string) string {
			return fmt.Sprintf("Collection: %s, Key: %s", collectionName, key)
		}))

	registerHostFunction("hypermode", "getTextsFromCollection",
		func(ctx context.Context, collectionName string) (map[string]string, error) {
			return collections.GetTextsFromCollection(ctx, collectionName, "")
		},
		withCancelledMessage("Cancelled getting texts from collection."),
		withErrorMessage("Error getting texts from collection."),
		withMessageDetail(func(collectionName string) string {
			return fmt.Sprintf("Collection: %s", collectionName)
		}))

	registerHostFunction("hypermode", "nnClassifyCollection",
		func(ctx context.Context, collectionName, searchMethod, text string) (*collections.CollectionClassificationResult, error) {
			return collections.NnClassify(ctx, collectionName, "", searchMethod, text)
		},
		withCancelledMessage("Cancelled classification."),
		withErrorMessage("Error during classification."),
		withMessageDetail(func(collectionName, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod)
		}))

	registerHostFunction("hypermode", "recomputeSearchMethod",
		func(ctx context.Context, collectionName, searchMethod string) (*collections.SearchMethodMutationResult, error) {
			return collections.RecomputeSearchMethod(ctx, collectionName, "", searchMethod)
		},
		withStartingMessage("Starting recomputing search method for collection."),
		withCompletedMessage("Completed recomputing search method for collection."),
		withCancelledMessage("Cancelled recomputing search method for collection."),
		withErrorMessage("Error recomputing search method for collection."),
		withMessageDetail(func(collectionName, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod)
		}))

	registerHostFunction("hypermode", "searchCollection",
		func(ctx context.Context, collectionName, searchMethod, text string, limit int32, returnText bool) (*collections.CollectionSearchResult, error) {
			return collections.SearchCollection(ctx, collectionName, nil, searchMethod, text, limit, returnText)
		},
		withCancelledMessage("Cancelled searching collection."),
		withErrorMessage("Error searching collection."),
		withMessageDetail(func(collectionName, searchMethod string) string {
			return fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod)
		}))

	registerHostFunction("hypermode", "upsertToCollection",
		func(ctx context.Context, collectionName string, keys, texts []string) (*collections.CollectionMutationResult, error) {
			return collections.UpsertToCollection(ctx, collectionName, "", keys, texts, nil)
		},
		withCancelledMessage("Cancelled collection upsert."),
		withErrorMessage("Error upserting to collection."),
		withMessageDetail(func(collectionName string, keys []string) string {
			return fmt.Sprintf("Collection: %s, Keys: %v", collectionName, keys)
		}))
}

//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package collections

import "github.com/hypermodeAI/functions-go/pkg/testutils"

var UpsertCallStack = testutils.NewCallStack()
var DeleteCallStack = testutils.NewCallStack()
var SearchCallStack = testutils.NewCallStack()
var NnClassifyCallStack = testutils.NewCallStack()
var RecomputeSearchMethodCallStack = testutils.NewCallStack()
var ComputeDistanceCallStack = testutils.NewCallStack()
var GetTextCallStack = testutils.NewCallStack()
var GetTextsCallStack = testutils.NewCallStack()
var GetNamespacesCallStack = testutils.NewCallStack()
var GetVectorCallStack = testutils.NewCallStack()
var GetLabelsCallStack = testutils.NewCallStack()
var SearchByVectorCallStack = testutils.NewCallStack()

func hostUpsertToCollection(collection, namespace *string, keys, texts *[]string, labels *[][]string) *CollectionMutationResult {
	UpsertCallStack.Push(collection, namespace, keys, texts, labels)

	return &CollectionMutationResult{
		Collection: *collection,
		Status:     "success",
	}
}

func hostDeleteFromCollection(collection, namespace, key *string) *CollectionMutationResult {
	DeleteCallStack.Push(collection, namespace, key)

	return &CollectionMutationResult{
		Collection: *collection,
		Status:     "success",
	}
}

func hostSearchCollection(collection *string, namespaces *[]string, searchMethod, text *string, limit int32, returnText bool) *CollectionSearchResult {
	SearchCallStack.Push(collection, namespaces, searchMethod, text, limit, returnText)

	return &CollectionSearchResult{
		Collection: *collection,
		Status:     "success",
	}
}

func hostNnClassifyCollection(collection, namespace, searchMethod, text *string) *CollectionClassificationResult {
	NnClassifyCallStack.Push(collection, namespace, searchMethod, text)

	return &CollectionClassificationResult{
		Collection: *collection,
		Status:     "success",
	}
}

func hostRecomputeSearchMethod(collection, namespace, searchMethod *string) *SearchMethodMutationResult {
	RecomputeSearchMethodCallStack.Push(collection, namespace, searchMethod)

	return &SearchMethodMutationResult{
		Collection: *collection,
		Status:     "success",
	}
}

func hostComputeDistance(collection, namespace, searchMethod, key1, key2 *string) *CollectionSearchResultObject {
	ComputeDistanceCallStack.Push(collection, namespace, searchMethod, key1, key2)

	return &CollectionSearchResultObject{
		Distance: 0.0,
		Score:    1.0,
	}
}

func hostGetTextFromCollection(collection, namespace, key *string) *string {
	GetTextCallStack.Push(collection, namespace, key)

	helloWorld := "Hello, World!"
	return &helloWorld
}

func hostGetTextsFromCollection(collection, namespace *string) *map[string]string {
	GetTextsCallStack.Push(collection, namespace)

	ret := map[string]string{
		"key1": "Hello, World!",
		"key2": "Hello, World!",
	}

	return &ret
}

func hostGetNamespacesFromCollection(collection *string) *[]string {
	GetNamespacesCallStack.Push(collection)

	ret := []string{"namespace1", "namespace2"}

	return &ret
}

func hostGetVector(collection, namespace, searchMethod, key *string) *[]float32 {
	GetVectorCallStack.Push(collection, namespace, searchMethod, key)

	ret := []float32{0.1, 0.2, 0.3}

	return &ret
}

func hostGetLabels(collection, namespace, key *string) *[]string {
	GetLabelsCallStack.Push(collection, namespace, key)

	ret := []string{"label1", "label2"}

	return &ret
}

func hostSearchCollectionByVector(collection *string, namespaces *[]string, searchMethod *string, vector *[]float32, limit int32, returnText bool) *CollectionSearchResult {
	SearchByVectorCallStack.Push(collection, namespaces, searchMethod, vector, limit, returnText)

	return &CollectionSearchResult{
		Collection: *collection,
		Status:     "success",
	}
}

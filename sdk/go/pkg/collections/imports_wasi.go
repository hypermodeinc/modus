//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import "unsafe"

//go:noescape
//go:wasmimport hypermode upsertToCollection_v2
func _hostUpsertToCollection(collection, namespace *string, keys, texts, labels unsafe.Pointer) unsafe.Pointer

//hypermode:import hypermode upsertToCollection_v2
func hostUpsertToCollection(collection, namespace *string, keys, texts *[]string, labels *[][]string) *CollectionMutationResult {
	keysPointer := unsafe.Pointer(keys)
	textsPointer := unsafe.Pointer(texts)
	labelsPointer := unsafe.Pointer(labels)
	response := _hostUpsertToCollection(collection, namespace, keysPointer, textsPointer, labelsPointer)
	if response == nil {
		return nil
	}
	return (*CollectionMutationResult)(response)
}

//go:noescape
//go:wasmimport hypermode deleteFromCollection_v2
func _hostDeleteFromCollection(collection, namespace, key *string) unsafe.Pointer

//hypermode:import hypermode deleteFromCollection_v2
func hostDeleteFromCollection(collection, namespace, key *string) *CollectionMutationResult {
	response := _hostDeleteFromCollection(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*CollectionMutationResult)(response)
}

//go:noescape
//go:wasmimport hypermode searchCollection_v2
func _hostSearchCollection(collection *string, namespaces unsafe.Pointer, searchMethod, text *string, limit int32, returnText bool) unsafe.Pointer

//hypermode:import hypermode searchCollection_v2
func hostSearchCollection(collection *string, namespaces *[]string, searchMethod, text *string, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	response := _hostSearchCollection(collection, namespacesPtr, searchMethod, text, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

//go:noescape
//go:wasmimport hypermode nnClassifyCollection_v2
func _hostNnClassifyCollection(collection, namespace, searchMethod, text *string) unsafe.Pointer

//hypermode:import hypermode nnClassifyCollection_v2
func hostNnClassifyCollection(collection, namespace, searchMethod, text *string) *CollectionClassificationResult {
	response := _hostNnClassifyCollection(collection, namespace, searchMethod, text)
	if response == nil {
		return nil
	}
	return (*CollectionClassificationResult)(response)
}

//go:noescape
//go:wasmimport hypermode recomputeSearchMethod_v2
func _hostRecomputeSearchMethod(collection, namespace, searchMethod *string) unsafe.Pointer

//hypermode:import hypermode recomputeSearchMethod_v2
func hostRecomputeSearchMethod(collection, namespace, searchMethod *string) *SearchMethodMutationResult {
	response := _hostRecomputeSearchMethod(collection, namespace, searchMethod)
	if response == nil {
		return nil
	}
	return (*SearchMethodMutationResult)(response)
}

//go:noescape
//go:wasmimport hypermode computeDistance_v2
func _hostComputeDistance(collection, namespace, searchMethod, key1, key2 *string) unsafe.Pointer

//hypermode:import hypermode computeDistance_v2
func hostComputeDistance(collection, namespace, searchMethod, key1, key2 *string) *CollectionSearchResultObject {
	response := _hostComputeDistance(collection, namespace, searchMethod, key1, key2)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResultObject)(response)
}

//go:noescape
//go:wasmimport hypermode getTextFromCollection_v2
func _hostGetTextFromCollection(collection, namespace, key *string) unsafe.Pointer

//hypermode:import hypermode getTextFromCollection_v2
func hostGetTextFromCollection(collection, namespace, key *string) *string {
	response := _hostGetTextFromCollection(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*string)(response)
}

//go:noescape
//go:wasmimport hypermode getTextsFromCollection_v2
func _hostGetTextsFromCollection(collection, namespace *string) unsafe.Pointer

//hypermode:import hypermode getTextsFromCollection_v2
func hostGetTextsFromCollection(collection, namespace *string) *map[string]string {
	response := _hostGetTextsFromCollection(collection, namespace)
	if response == nil {
		return nil
	}
	return (*map[string]string)(response)
}

//go:noescape
//go:wasmimport hypermode getNamespacesFromCollection
func _hostGetNamespacesFromCollection(collection *string) unsafe.Pointer

//hypermode:import hypermode getNamespacesFromCollection
func hostGetNamespacesFromCollection(collection *string) *[]string {
	response := _hostGetNamespacesFromCollection(collection)
	if response == nil {
		return nil
	}
	return (*[]string)(response)
}

//go:noescape
//go:wasmimport hypermode getVector
func _hostGetVector(collection, namespace, searchMethod, key *string) unsafe.Pointer

//hypermode:import hypermode getVector
func hostGetVector(collection, namespace, searchMethod, key *string) *[]float32 {
	response := _hostGetVector(collection, namespace, searchMethod, key)
	if response == nil {
		return nil
	}
	return (*[]float32)(response)
}

//go:noescape
//go:wasmimport hypermode getLabels
func _hostGetLabels(collection, namespace, key *string) unsafe.Pointer

//hypermode:import hypermode getLabels
func hostGetLabels(collection, namespace, key *string) *[]string {
	response := _hostGetLabels(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*[]string)(response)
}

//go:noescape
//go:wasmimport hypermode searchCollectionByVector
func _hostSearchCollectionByVector(collection *string, namespaces unsafe.Pointer, searchMethod *string, vector unsafe.Pointer, limit int32, returnText bool) unsafe.Pointer

//hypermode:import hypermode searchCollectionByVector
func hostSearchCollectionByVector(collection *string, namespaces *[]string, searchMethod *string, vector *[]float32, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	vectorPtr := unsafe.Pointer(vector)
	response := _hostSearchCollectionByVector(collection, namespacesPtr, searchMethod, vectorPtr, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

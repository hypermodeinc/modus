//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import "unsafe"

//go:noescape
//go:wasmimport modus_collections upsert
func _hostUpsertToCollection(collection, namespace *string, keys, texts, labels unsafe.Pointer) unsafe.Pointer

//modus:import modus_collections upsert
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
//go:wasmimport modus_collections delete
func _hostDeleteFromCollection(collection, namespace, key *string) unsafe.Pointer

//modus:import modus_collections delete
func hostDeleteFromCollection(collection, namespace, key *string) *CollectionMutationResult {
	response := _hostDeleteFromCollection(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*CollectionMutationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections search
func _hostSearchCollection(collection *string, namespaces unsafe.Pointer, searchMethod, text *string, limit int32, returnText bool) unsafe.Pointer

//modus:import modus_collections search
func hostSearchCollection(collection *string, namespaces *[]string, searchMethod, text *string, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	response := _hostSearchCollection(collection, namespacesPtr, searchMethod, text, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

//go:noescape
//go:wasmimport modus_collections classifyText
func _hostNnClassifyCollection(collection, namespace, searchMethod, text *string) unsafe.Pointer

//modus:import modus_collections classifyText
func hostNnClassifyCollection(collection, namespace, searchMethod, text *string) *CollectionClassificationResult {
	response := _hostNnClassifyCollection(collection, namespace, searchMethod, text)
	if response == nil {
		return nil
	}
	return (*CollectionClassificationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections recomputeIndex
func _hostRecomputeSearchMethod(collection, namespace, searchMethod *string) unsafe.Pointer

//modus:import modus_collections recomputeIndex
func hostRecomputeSearchMethod(collection, namespace, searchMethod *string) *SearchMethodMutationResult {
	response := _hostRecomputeSearchMethod(collection, namespace, searchMethod)
	if response == nil {
		return nil
	}
	return (*SearchMethodMutationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections computeDistance
func _hostComputeDistance(collection, namespace, searchMethod, key1, key2 *string) unsafe.Pointer

//modus:import modus_collections computeDistance
func hostComputeDistance(collection, namespace, searchMethod, key1, key2 *string) *CollectionSearchResultObject {
	response := _hostComputeDistance(collection, namespace, searchMethod, key1, key2)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResultObject)(response)
}

//go:noescape
//go:wasmimport modus_collections getText
func _hostGetTextFromCollection(collection, namespace, key *string) unsafe.Pointer

//modus:import modus_collections getText
func hostGetTextFromCollection(collection, namespace, key *string) *string {
	response := _hostGetTextFromCollection(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*string)(response)
}

//go:noescape
//go:wasmimport modus_collections dumpTexts
func _hostGetTextsFromCollection(collection, namespace *string) unsafe.Pointer

//modus:import modus_collections dumpTexts
func hostGetTextsFromCollection(collection, namespace *string) *map[string]string {
	response := _hostGetTextsFromCollection(collection, namespace)
	if response == nil {
		return nil
	}
	return (*map[string]string)(response)
}

//go:noescape
//go:wasmimport modus_collections getNamespaces
func _hostGetNamespacesFromCollection(collection *string) unsafe.Pointer

//modus:import modus_collections getNamespaces
func hostGetNamespacesFromCollection(collection *string) *[]string {
	response := _hostGetNamespacesFromCollection(collection)
	if response == nil {
		return nil
	}
	return (*[]string)(response)
}

//go:noescape
//go:wasmimport modus_collections getVector
func _hostGetVector(collection, namespace, searchMethod, key *string) unsafe.Pointer

//modus:import modus_collections getVector
func hostGetVector(collection, namespace, searchMethod, key *string) *[]float32 {
	response := _hostGetVector(collection, namespace, searchMethod, key)
	if response == nil {
		return nil
	}
	return (*[]float32)(response)
}

//go:noescape
//go:wasmimport modus_collections getLabels
func _hostGetLabels(collection, namespace, key *string) unsafe.Pointer

//modus:import modus_collections getLabels
func hostGetLabels(collection, namespace, key *string) *[]string {
	response := _hostGetLabels(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*[]string)(response)
}

//go:noescape
//go:wasmimport modus_collections searchByVector
func _hostSearchCollectionByVector(collection *string, namespaces unsafe.Pointer, searchMethod *string, vector unsafe.Pointer, limit int32, returnText bool) unsafe.Pointer

//modus:import modus_collections searchByVector
func hostSearchCollectionByVector(collection *string, namespaces *[]string, searchMethod *string, vector *[]float32, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	vectorPtr := unsafe.Pointer(vector)
	response := _hostSearchCollectionByVector(collection, namespacesPtr, searchMethod, vectorPtr, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

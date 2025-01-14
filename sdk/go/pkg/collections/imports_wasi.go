//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import "unsafe"

//go:noescape
//go:wasmimport modus_collections upsert
func _hostUpsert(collection, namespace *string, keys, texts, labels unsafe.Pointer) unsafe.Pointer

//modus:import modus_collections upsert
func hostUpsert(collection, namespace *string, keys, texts *[]string, labels *[][]string) *CollectionMutationResult {
	keysPointer := unsafe.Pointer(keys)
	textsPointer := unsafe.Pointer(texts)
	labelsPointer := unsafe.Pointer(labels)
	response := _hostUpsert(collection, namespace, keysPointer, textsPointer, labelsPointer)
	if response == nil {
		return nil
	}
	return (*CollectionMutationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections delete
func _hostDelete(collection, namespace, key *string) unsafe.Pointer

//modus:import modus_collections delete
func hostDelete(collection, namespace, key *string) *CollectionMutationResult {
	response := _hostDelete(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*CollectionMutationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections search
func _hostSearch(collection *string, namespaces unsafe.Pointer, searchMethod, text *string, limit int32, returnText bool) unsafe.Pointer

//modus:import modus_collections search
func hostSearch(collection *string, namespaces *[]string, searchMethod, text *string, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	response := _hostSearch(collection, namespacesPtr, searchMethod, text, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

//go:noescape
//go:wasmimport modus_collections classifyText
func _hostClassifyText(collection, namespace, searchMethod, text *string) unsafe.Pointer

//modus:import modus_collections classifyText
func hostClassifyText(collection, namespace, searchMethod, text *string) *CollectionClassificationResult {
	response := _hostClassifyText(collection, namespace, searchMethod, text)
	if response == nil {
		return nil
	}
	return (*CollectionClassificationResult)(response)
}

//go:noescape
//go:wasmimport modus_collections recomputeIndex
func _hostRecomputeIndex(collection, namespace, searchMethod *string) unsafe.Pointer

//modus:import modus_collections recomputeIndex
func hostRecomputeIndex(collection, namespace, searchMethod *string) *SearchMethodMutationResult {
	response := _hostRecomputeIndex(collection, namespace, searchMethod)
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
func _hostGetText(collection, namespace, key *string) unsafe.Pointer

//modus:import modus_collections getText
func hostGetText(collection, namespace, key *string) *string {
	response := _hostGetText(collection, namespace, key)
	if response == nil {
		return nil
	}
	return (*string)(response)
}

//go:noescape
//go:wasmimport modus_collections dumpTexts
func _hostDumpTexts(collection, namespace *string) unsafe.Pointer

//modus:import modus_collections dumpTexts
func hostDumpTexts(collection, namespace *string) *map[string]string {
	response := _hostDumpTexts(collection, namespace)
	if response == nil {
		return nil
	}
	return (*map[string]string)(response)
}

//go:noescape
//go:wasmimport modus_collections getNamespaces
func _hostGetNamespaces(collection *string) unsafe.Pointer

//modus:import modus_collections getNamespaces
func hostGetNamespaces(collection *string) *[]string {
	response := _hostGetNamespaces(collection)
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
func _hostSearchByVector(collection *string, namespaces unsafe.Pointer, searchMethod *string, vector unsafe.Pointer, limit int32, returnText bool) unsafe.Pointer

//modus:import modus_collections searchByVector
func hostSearchByVector(collection *string, namespaces *[]string, searchMethod *string, vector *[]float32, limit int32, returnText bool) *CollectionSearchResult {
	namespacesPtr := unsafe.Pointer(namespaces)
	vectorPtr := unsafe.Pointer(vector)
	response := _hostSearchByVector(collection, namespacesPtr, searchMethod, vectorPtr, limit, returnText)
	if response == nil {
		return nil
	}
	return (*CollectionSearchResult)(response)
}

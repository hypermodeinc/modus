/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"

	"github.com/hypermodeinc/modus/runtime/collections/in_mem"
	"github.com/hypermodeinc/modus/runtime/collections/index"
	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
	collection_utils "github.com/hypermodeinc/modus/runtime/collections/utils"
	"github.com/hypermodeinc/modus/runtime/functions"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/utils"
	"github.com/hypermodeinc/modus/runtime/wasmhost"
)

var errInvalidEmbedderSignature = errors.New("invalid embedder function signature")

func Initialize(ctx context.Context) {
	globalNamespaceManager = newCollectionFactory()
	manifestdata.RegisterManifestLoadedCallback(cleanAndProcessManifest)
	functions.RegisterFunctionsLoadedCallback(func(ctx context.Context) {
		globalNamespaceManager.readFromPostgres(ctx)
	})

	go globalNamespaceManager.worker(ctx)
}

func Shutdown(ctx context.Context) {
	close(globalNamespaceManager.quit)
	<-globalNamespaceManager.done
}

func UpsertToCollection(ctx context.Context, collectionName, namespace string, keys, texts []string, labels [][]string) (*CollectionMutationResult, error) {

	// Get the collectionName data from the manifest
	collectionData := manifestdata.GetManifest().Collections[collectionName]

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := func(namespace string, index interfaces.CollectionNamespace) (interfaces.CollectionNamespace, error) {
		return col.findOrCreateNamespace(namespace, index)
	}(namespace, in_mem.NewCollectionNamespace(collectionName, namespace))
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		keys = make([]string, len(texts))
		for i := range keys {
			keys[i] = utils.GenerateUUIDv7()
		}
	}

	if len(labels) != 0 && len(labels) != len(texts) {
		return nil, fmt.Errorf("mismatch in number of labels and texts: %d != %d", len(labels), len(texts))
	}

	err = collNs.InsertTexts(ctx, keys, texts, labels)
	if err != nil {
		return nil, err
	}

	// compute embeddings for each search method, and insert into vector index
	for searchMethodName, searchMethod := range collectionData.SearchMethods {
		vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethodName)
		if err == index.ErrVectorIndexNotFound {
			vectorIndex, err = createIndexObject(searchMethod, searchMethodName)
			if err != nil {
				return nil, err
			}
			err = collNs.SetVectorIndex(ctx, searchMethodName, vectorIndex)
			if err != nil {
				return nil, err
			}
		} else if err != nil {
			return nil, err
		}

		embedder := searchMethod.Embedder
		if err := validateEmbedder(ctx, embedder); err != nil {
			return nil, err
		}

		executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
		if err != nil {
			return nil, err
		}

		result := executionInfo.Result()

		textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
		if err != nil {
			return nil, err
		}

		if len(textVecs) != len(texts) {
			return nil, fmt.Errorf("mismatch in number of embeddings generated by embedder %s", embedder)
		}

		ids := make([]int64, len(keys))
		for i := range textVecs {
			key := keys[i]

			id, err := collNs.GetExternalId(ctx, key)
			if err != nil {
				return nil, err
			}
			ids[i] = id
		}

		err = vectorIndex.InsertVectors(ctx, ids, textVecs)
		if err != nil {
			return nil, err
		}
	}

	return NewCollectionMutationResult(collectionName, "upsert", "success", keys, ""), nil
}

func DeleteFromCollection(ctx context.Context, collectionName, namespace, key string) (*CollectionMutationResult, error) {
	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	textId, err := collNs.GetExternalId(ctx, key)
	if err != nil {
		return nil, err
	}
	for _, vectorIndex := range collNs.GetVectorIndexMap() {
		err = vectorIndex.DeleteVector(ctx, textId, key)
		if err != nil {
			return nil, err
		}
	}
	err = collNs.DeleteText(ctx, key)
	if err != nil {
		return nil, err
	}

	keys := []string{key}

	return NewCollectionMutationResult(collectionName, "delete", "success", keys, ""), nil
}

func SearchCollection(ctx context.Context, collectionName string, namespaces []string, searchMethod, text string, limit int32, returnText bool) (*CollectionSearchResult, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if len(namespaces) == 0 {
		namespaces = []string{in_mem.DefaultNamespace}
	}

	embedder, err := getEmbedder(ctx, collectionName, searchMethod)
	if err != nil {
		return nil, err
	}

	texts := []string{text}

	executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
	if err != nil {
		return nil, err
	}

	result := executionInfo.Result()

	textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		return nil, err
	}

	if len(textVecs) == 0 {
		return nil, fmt.Errorf("no embeddings generated by embedder %s", embedder)
	}

	// merge all objects
	mergedObjects := make([]*CollectionSearchResultObject, 0, len(namespaces)*int(limit))
	for _, ns := range namespaces {
		collNs, err := col.findNamespace(ns)
		if err != nil {
			return nil, err
		}

		vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
		if err != nil {
			return nil, err
		}

		objects, err := vectorIndex.Search(ctx, textVecs[0], int(limit), nil)
		if err != nil {
			return nil, err
		}

		for _, object := range objects {
			text, err := collNs.GetText(ctx, object.GetIndex())
			if err != nil {
				return nil, err
			}
			labels, err := collNs.GetLabels(ctx, object.GetIndex())
			if err != nil {
				return nil, err
			}
			mergedObjects = append(mergedObjects, NewCollectionSearchResultObject(ns, object.GetIndex(), text, labels, object.GetValue(), 1-object.GetValue()))
		}
	}

	// sort by score
	sort.Slice(mergedObjects, func(i, j int) bool {
		return mergedObjects[i].Distance < mergedObjects[j].Distance
	})

	if len(mergedObjects) > int(limit) {
		mergedObjects = mergedObjects[:int(limit)]
	}

	return NewCollectionSearchResult(collectionName, searchMethod, "success", mergedObjects, ""), nil
}

func SearchCollectionByVector(ctx context.Context, collectionName string, namespaces []string, searchMethod string, vector []float32, limit int32, returnText bool) (*CollectionSearchResult, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if len(namespaces) == 0 {
		namespaces = []string{in_mem.DefaultNamespace}
	}

	// merge all objects
	mergedObjects := make([]*CollectionSearchResultObject, 0, len(namespaces)*int(limit))
	for _, ns := range namespaces {
		collNs, err := col.findNamespace(ns)
		if err != nil {
			return nil, err
		}

		vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
		if err != nil {
			return nil, err
		}

		objects, err := vectorIndex.Search(ctx, vector, int(limit), nil)
		if err != nil {
			return nil, err
		}

		for _, object := range objects {
			text, err := collNs.GetText(ctx, object.GetIndex())
			if err != nil {
				return nil, err
			}
			labels, err := collNs.GetLabels(ctx, object.GetIndex())
			if err != nil {
				return nil, err
			}
			mergedObjects = append(mergedObjects, NewCollectionSearchResultObject(ns, object.GetIndex(), text, labels, object.GetValue(), 1-object.GetValue()))
		}
	}

	// sort by score
	sort.Slice(mergedObjects, func(i, j int) bool {
		return mergedObjects[i].Distance < mergedObjects[j].Distance
	})

	if len(mergedObjects) > int(limit) {
		mergedObjects = mergedObjects[:int(limit)]
	}

	return NewCollectionSearchResult(collectionName, searchMethod, "success", mergedObjects, ""), nil
}

func NnClassify(ctx context.Context, collectionName, namespace, searchMethod, text string) (*CollectionClassificationResult, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	embedder, err := getEmbedder(ctx, collectionName, searchMethod)
	if err != nil {
		return nil, err
	}

	texts := []string{text}

	executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
	if err != nil {
		return nil, err
	}

	result := executionInfo.Result()

	textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		return nil, err
	}

	if len(textVecs) == 0 {
		return nil, fmt.Errorf("no embeddings generated by embedder %s", embedder)
	}

	lenTexts, err := collNs.Len(ctx)
	if err != nil {
		return nil, err
	}

	nns, err := vectorIndex.Search(ctx, textVecs[0], int(math.Log10(float64(lenTexts)))*int(math.Log10(float64(lenTexts))), nil)
	if err != nil {
		return nil, err
	}

	// remove elements with score out of first standard deviation

	// calculate mean
	var sum float64
	for _, nn := range nns {
		sum += nn.GetValue()
	}
	mean := sum / float64(len(nns))

	// calculate standard deviation
	var variance float64
	for _, nn := range nns {
		variance += math.Pow(float64(nn.GetValue())-mean, 2)
	}
	stdDev := math.Sqrt(variance / float64(len(nns)))

	// remove elements with score out of first standard deviation and return the most frequent label
	labelCounts := make(map[string]int)

	res := NewCollectionClassificationResult(collectionName, searchMethod, "success", []*CollectionClassificationLabelObject{}, []*CollectionClassificationResultObject{}, "")

	totalLabels := 0

	for _, nn := range nns {
		if math.Abs(nn.GetValue()-mean) <= 2*stdDev {
			labels, err := collNs.GetLabels(ctx, nn.GetIndex())
			if err != nil {
				return nil, err
			}
			for _, label := range labels {
				labelCounts[label]++
				totalLabels++
			}

			res.Cluster = append(res.Cluster, NewCollectionClassificationResultObject(nn.GetIndex(), labels, nn.GetValue(), 1-nn.GetValue()))
		}
	}

	// Create a slice of pairs
	labelsResult := make([]*CollectionClassificationLabelObject, 0, len(labelCounts))
	for label, count := range labelCounts {
		labelsResult = append(labelsResult, NewCollectionClassificationLabelObject(label, float64(count)/float64(totalLabels)))
	}

	// Sort the pairs by count in descending order
	sort.Slice(labelsResult, func(i, j int) bool {
		return labelsResult[i].Confidence > labelsResult[j].Confidence
	})

	res.LabelsResult = labelsResult

	return res, nil
}

func GetVector(ctx context.Context, collectionName, namespace, searchMethod, id string) ([]float32, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	vec, err := vectorIndex.GetVector(ctx, id)
	if err != nil {
		return nil, err
	}

	return vec, nil
}

func GetLabels(ctx context.Context, collectionName, namespace, key string) ([]string, error) {
	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	labels, err := collNs.GetLabels(ctx, key)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func ComputeDistance(ctx context.Context, collectionName, namespace, searchMethod, id1, id2 string) (*CollectionSearchResultObject, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	vec1, err := vectorIndex.GetVector(ctx, id1)
	if err != nil {
		return nil, err
	}

	if vec1 == nil {
		return nil, fmt.Errorf("vector for id %s not found", id1)
	}

	vec2, err := vectorIndex.GetVector(ctx, id2)
	if err != nil {
		return nil, err
	}

	if vec2 == nil {
		return nil, fmt.Errorf("vector for id %s not found", id2)
	}

	distance, err := collection_utils.CosineDistance(vec1, vec2)
	if err != nil {
		return nil, err
	}

	return NewCollectionSearchResultObject(namespace, "", "", []string{}, distance, 1-distance), nil
}

func RecomputeSearchMethod(ctx context.Context, collectionName, namespace, searchMethod string) (*SearchMethodMutationResult, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	vectorIndex, err := collNs.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		return nil, err
	}

	err = processTextMap(ctx, collNs, interfaces.VectorIndex(vectorIndex))
	if err != nil {
		return nil, err
	}

	return NewSearchMethodMutationResult(collectionName, searchMethod, "recompute", "success", ""), nil
}

func GetTextFromCollection(ctx context.Context, collectionName, namespace, key string) (string, error) {
	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return "", err
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return "", err
	}

	text, err := collNs.GetText(ctx, key)
	if err != nil {
		return "", err
	}

	return text, nil
}

func GetTextsFromCollection(ctx context.Context, collectionName, namespace string) (map[string]string, error) {

	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	if namespace == "" {
		namespace = in_mem.DefaultNamespace
	}

	collNs, err := col.findNamespace(namespace)
	if err != nil {
		return nil, err
	}

	textMap, err := collNs.GetTextMap(ctx)
	if err != nil {
		return nil, err
	}
	if textMap == nil {
		return make(map[string]string), nil
	}

	return textMap, nil
}

func GetNamespacesFromCollection(ctx context.Context, collectionName string) ([]string, error) {
	col, err := globalNamespaceManager.findCollection(collectionName)
	if err != nil {
		return nil, err
	}

	namespaceMap := col.getCollectionNamespaceMap()

	namespaces := make([]string, 0, len(namespaceMap))
	for namespace := range namespaceMap {
		namespaces = append(namespaces, namespace)
	}

	return namespaces, nil
}

func getEmbedder(ctx context.Context, collectionName string, searchMethod string) (string, error) {
	manifestColl, ok := manifestdata.GetManifest().Collections[collectionName]
	if !ok {
		return "", fmt.Errorf("collection %s not found in manifest", collectionName)
	}

	manifestSearchMethod, ok := manifestColl.SearchMethods[searchMethod]
	if !ok {
		return "", fmt.Errorf("search method %s not found in collection %s", searchMethod, collectionName)
	}

	embedder := manifestSearchMethod.Embedder
	if embedder == "" {
		return "", fmt.Errorf("embedder not found in search method %s of collection %s", searchMethod, collectionName)
	}

	if err := validateEmbedder(ctx, embedder); err != nil {
		return "", err
	}

	return embedder, nil
}

func validateEmbedder(ctx context.Context, embedder string) error {

	info, err := wasmhost.GetWasmHost(ctx).GetFunctionInfo(embedder)
	if err != nil {
		return err
	}
	fn := info.Metadata()

	// Embedder functions must take a single string[] parameter and return a single f32[][] or f64[][] array.
	// The types are language-specific, so we use the plugin language's type info.

	if len(fn.Parameters) != 1 || len(fn.Results) != 1 {
		return errInvalidEmbedderSignature
	}

	lti := info.Plugin().Language.TypeInfo()

	p := fn.Parameters[0]
	if !lti.IsListType(p.Type) || !lti.IsStringType(lti.GetListSubtype(p.Type)) {
		return errInvalidEmbedderSignature
	}

	r := fn.Results[0]
	if !lti.IsListType(r.Type) {
		return errInvalidEmbedderSignature
	}

	a := lti.GetListSubtype(r.Type)
	if !lti.IsListType(a) || !lti.IsFloatType(lti.GetListSubtype(a)) {
		return errInvalidEmbedderSignature
	}

	return nil
}

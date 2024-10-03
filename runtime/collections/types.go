/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

func NewCollectionMutationResult(collection, operation, status string, keys []string, err string) *CollectionMutationResult {
	if keys == nil {
		keys = []string{}
	}
	return &CollectionMutationResult{
		Collection: collection,
		Operation:  operation,
		Status:     status,
		Keys:       keys,
		Error:      err,
	}
}

type CollectionMutationResult struct {
	Collection string
	Operation  string
	Status     string
	Keys       []string
	Error      string
}

func NewSearchMethodMutationResult(collection, searchMethod, operation, status, err string) *SearchMethodMutationResult {
	return &SearchMethodMutationResult{
		Collection:   collection,
		SearchMethod: searchMethod,
		Operation:    operation,
		Status:       status,
		Error:        err,
	}
}

type SearchMethodMutationResult struct {
	Collection   string
	SearchMethod string
	Operation    string
	Status       string
	Error        string
}

func NewCollectionSearchResult(collection, searchMethod, status string, objects []*CollectionSearchResultObject, err string) *CollectionSearchResult {
	if objects == nil {
		objects = []*CollectionSearchResultObject{}
	}
	return &CollectionSearchResult{
		Collection:   collection,
		SearchMethod: searchMethod,
		Status:       status,
		Objects:      objects,
		Error:        err,
	}
}

type CollectionSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []*CollectionSearchResultObject
	Error        string
}

func NewCollectionSearchResultObject(namespace, key, text string, labels []string, distance, score float64) *CollectionSearchResultObject {
	if labels == nil {
		labels = []string{}
	}
	return &CollectionSearchResultObject{
		Namespace: namespace,
		Key:       key,
		Text:      text,
		Labels:    labels,
		Distance:  distance,
		Score:     score,
	}
}

type CollectionSearchResultObject struct {
	Namespace string
	Key       string
	Text      string
	Labels    []string
	Distance  float64
	Score     float64
}

func NewCollectionClassificationResult(collection, searchMethod, status string, labelsResult []*CollectionClassificationLabelObject, cluster []*CollectionClassificationResultObject, err string) *CollectionClassificationResult {
	if labelsResult == nil {
		labelsResult = []*CollectionClassificationLabelObject{}
	}
	if cluster == nil {
		cluster = []*CollectionClassificationResultObject{}
	}
	return &CollectionClassificationResult{
		Collection:   collection,
		SearchMethod: searchMethod,
		Status:       status,
		LabelsResult: labelsResult,
		Cluster:      cluster,
		Error:        err,
	}
}

type CollectionClassificationResult struct {
	Collection   string
	SearchMethod string
	Status       string
	LabelsResult []*CollectionClassificationLabelObject
	Cluster      []*CollectionClassificationResultObject
	Error        string
}

func NewCollectionClassificationLabelObject(label string, confidence float64) *CollectionClassificationLabelObject {
	return &CollectionClassificationLabelObject{
		Label:      label,
		Confidence: confidence,
	}
}

type CollectionClassificationLabelObject struct {
	Label      string
	Confidence float64
}

func NewCollectionClassificationResultObject(key string, labels []string, distance, score float64) *CollectionClassificationResultObject {
	if labels == nil {
		labels = []string{}
	}
	return &CollectionClassificationResultObject{
		Key:      key,
		Labels:   labels,
		Distance: distance,
		Score:    score,
	}
}

type CollectionClassificationResultObject struct {
	Key      string
	Labels   []string
	Distance float64
	Score    float64
}

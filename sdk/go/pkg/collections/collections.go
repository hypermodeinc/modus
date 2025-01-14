/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package collections

import (
	"fmt"
)

type CollectionStatus = string

const (
	Success CollectionStatus = "success"
	Error   CollectionStatus = "error"
)

type CollectionMutationResult struct {
	Collection string
	Status     string
	Error      string
	Operation  string
	Keys       []string
}

type SearchMethodMutationResult struct {
	Collection   string
	Status       string
	Error        string
	Operation    string
	SearchMethod string
}

type CollectionSearchResult struct {
	Collection   string
	Status       string
	Error        string
	SearchMethod string
	Objects      []*CollectionSearchResultObject
}

type CollectionSearchResultObject struct {
	Namespace string
	Key       string
	Text      string
	Labels    []string
	Distance  float64
	Score     float64
}

type CollectionClassificationResult struct {
	Collection   string
	Status       string
	Error        string
	SearchMethod string
	LabelsResult []*CollectionClassificationLabelObject
	Cluster      []*CollectionClassificationResultObject
}

type CollectionClassificationLabelObject struct {
	Label      string
	Confidence float64
}

type CollectionClassificationResultObject struct {
	Key      string
	Labels   []string
	Distance float64
	Score    float64
}

type NamespaceOption func(*NamespaceOptions)

type NamespaceOptions struct {
	namespace string
}

func WithNamespace(namespace string) NamespaceOption {
	return func(o *NamespaceOptions) {
		o.namespace = namespace
	}
}

func UpsertBatch(collection string, keys []string, texts []string, labelsArr [][]string, opts ...NamespaceOption) (*CollectionMutationResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if len(texts) == 0 {
		return nil, fmt.Errorf("Texts is empty")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	if keys == nil {
		keys = []string{}
	}

	if labelsArr == nil {
		labelsArr = [][]string{}
	}

	result := hostUpsert(&collection, &nsOpts.namespace, &keys, &texts, &labelsArr)

	if result == nil {
		return nil, fmt.Errorf("Failed to upsert")
	}

	return result, nil
}

func Upsert(collection string, key *string, text string, labels []string, opts ...NamespaceOption) (*CollectionMutationResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	keyArr := []string{}

	if key != nil {
		keyArr = []string{*key}
	}

	labelsArr := [][]string{}

	if labels != nil {
		labelsArr = [][]string{labels}
	}

	if text == "" {
		return nil, fmt.Errorf("Text is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostUpsert(&collection, &nsOpts.namespace, &keyArr, &[]string{text}, &labelsArr)

	if result == nil {
		return nil, fmt.Errorf("Failed to upsert")
	}

	return result, nil
}

func Remove(collection, key string, opts ...NamespaceOption) (*CollectionMutationResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if key == "" {
		return nil, fmt.Errorf("Key is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostDelete(&collection, &nsOpts.namespace, &key)

	if result == nil {
		return nil, fmt.Errorf("Failed to delete")
	}

	return result, nil
}

type SearchOption func(*SearchOptions)

type SearchOptions struct {
	namespaces []string
	limit      int
	returnText bool
}

func WithNamespaces(namespaces []string) SearchOption {
	return func(o *SearchOptions) {
		o.namespaces = namespaces
	}
}

func WithLimit(limit int) SearchOption {
	return func(o *SearchOptions) {
		o.limit = limit
	}
}

func WithReturnText(returnText bool) SearchOption {
	return func(o *SearchOptions) {
		o.returnText = returnText
	}
}

func Search(collection, searchMethod, text string, opts ...SearchOption) (*CollectionSearchResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	if text == "" {
		return nil, fmt.Errorf("Text is required")
	}

	sOpts := &SearchOptions{
		namespaces: []string{},
		limit:      10,
		returnText: false,
	}

	for _, opt := range opts {
		opt(sOpts)
	}

	result := hostSearch(&collection, &sOpts.namespaces, &searchMethod, &text, int32(sOpts.limit), sOpts.returnText)

	if result == nil {
		return nil, fmt.Errorf("Failed to search")
	}

	return result, nil
}

func SearchByVector(collection, searchMethod string, vector []float32, opts ...SearchOption) (*CollectionSearchResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	if len(vector) == 0 {
		return &CollectionSearchResult{
			Collection:   collection,
			Status:       Error,
			Error:        "Vector is required",
			SearchMethod: "",
			Objects:      nil,
		}, fmt.Errorf("Vector is required")
	}

	sOpts := &SearchOptions{
		namespaces: []string{},
		limit:      10,
		returnText: false,
	}

	for _, opt := range opts {
		opt(sOpts)
	}

	result := hostSearchByVector(&collection, &sOpts.namespaces, &searchMethod, &vector, int32(sOpts.limit), sOpts.returnText)

	if result == nil {
		return nil, fmt.Errorf("Failed to search")
	}

	return result, nil
}

func NnClassify(collection, searchMethod, text string, opts ...NamespaceOption) (*CollectionClassificationResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	if text == "" {
		return nil, fmt.Errorf("Text is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostClassifyText(&collection, &nsOpts.namespace, &searchMethod, &text)

	if result == nil {
		return nil, fmt.Errorf("Failed to classify")
	}

	return result, nil
}

func RecomputeSearchMethod(collection, searchMethod string, opts ...NamespaceOption) (*SearchMethodMutationResult, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostRecomputeIndex(&collection, &nsOpts.namespace, &searchMethod)

	if result == nil {
		return nil, fmt.Errorf("Failed to recompute")
	}

	return result, nil
}

func ComputeDistance(collection, searchMethod, key1, key2 string, opts ...NamespaceOption) (*CollectionSearchResultObject, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	if key1 == "" {
		return nil, fmt.Errorf("Key1 is required")
	}

	if key2 == "" {
		return nil, fmt.Errorf("Key2 is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostComputeDistance(&collection, &nsOpts.namespace, &searchMethod, &key1, &key2)

	if result == nil {
		return nil, fmt.Errorf("Failed to compute distance")
	}

	return result, nil
}

func GetText(collection, key string, opts ...NamespaceOption) (string, error) {
	if collection == "" {
		return "", fmt.Errorf("Collection name is required")
	}

	if key == "" {
		return "", fmt.Errorf("Key is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostGetText(&collection, &nsOpts.namespace, &key)

	if result == nil {
		return "", fmt.Errorf("Failed to get text for key")
	}

	return *result, nil
}

func GetTexts(collection string, opts ...NamespaceOption) (map[string]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostDumpTexts(&collection, &nsOpts.namespace)

	if result == nil {
		return nil, fmt.Errorf("Failed to get texts")
	}

	return *result, nil
}

func GetNamespaces(collection string) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	result := hostGetNamespaces(&collection)

	if result == nil {
		return nil, fmt.Errorf("Failed to get namespaces")
	}

	return *result, nil
}

func GetVector(collection, searchMethod, key string, opts ...NamespaceOption) ([]float32, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if searchMethod == "" {
		return nil, fmt.Errorf("Search method is required")
	}

	if key == "" {
		return nil, fmt.Errorf("Key is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostGetVector(&collection, &nsOpts.namespace, &searchMethod, &key)

	if result == nil {
		return nil, fmt.Errorf("Failed to get vector for key")
	}

	return *result, nil
}

func GetLabels(collection, key string, opts ...NamespaceOption) ([]string, error) {
	if collection == "" {
		return nil, fmt.Errorf("Collection name is required")
	}

	if key == "" {
		return nil, fmt.Errorf("Key is required")
	}

	nsOpts := &NamespaceOptions{
		namespace: "",
	}

	for _, opt := range opts {
		opt(nsOpts)
	}

	result := hostGetLabels(&collection, &nsOpts.namespace, &key)

	if result == nil {
		return []string{}, nil
	}

	return *result, nil
}

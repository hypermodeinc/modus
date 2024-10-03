/*
 * Copyright 2024 Hypermode, Inc.
 */

package collections_test

import (
	"reflect"
	"testing"

	"github.com/hypermodeAI/functions-go/pkg/collections"
)

var (
	collection   = "collection"
	searchMethod = "searchMethod"
	namespace    = "namespace"
	key          = "key"
	key1         = "key1"
	key2         = "key2"
	keyArr       = []string{"key"}
	text         = "text"
	textArr      = []string{"text"}
	labels       = []string{"label"}
	labelsArr    = [][]string{{"label"}}
)

func TestHostUpsertBatchToCollection(t *testing.T) {
	result, err := collections.UpsertBatch(collection, keyArr, textArr, labelsArr, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionMutationResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.UpsertCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&keyArr, values[2]) {
			t.Errorf("Expected keys: %v, but received: %v", &keyArr, values[2])
		}
		if !reflect.DeepEqual(&textArr, values[3]) {
			t.Errorf("Expected texts: %v, but received: %v", &textArr, values[3])
		}
		if !reflect.DeepEqual(&labelsArr, values[4]) {
			t.Errorf("Expected labels: %v, but received: %v", &labelsArr, values[4])
		}
	}
}

func TestHostUpsertToCollection(t *testing.T) {
	result, err := collections.Upsert(collection, nil, text, labels, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionMutationResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.UpsertCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&[]string{}, values[2]) {
			t.Errorf("Expected keys: %v, but received: %v", &keyArr, values[2])
		}
		if !reflect.DeepEqual(&textArr, values[3]) {
			t.Errorf("Expected texts: %v, but received: %v", &textArr, values[3])
		}
		if !reflect.DeepEqual(&labelsArr, values[4]) {
			t.Errorf("Expected labels: %v, but received: %v", &labelsArr, values[4])
		}
	}
}

func TestHostRemoveFromCollection(t *testing.T) {
	result, err := collections.Remove(collection, key, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionMutationResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.DeleteCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&key, values[2]) {
			t.Errorf("Expected key: %v, but received: %v", &key, values[2])
		}
	}
}

func TestHostSearchCollection(t *testing.T) {
	result, err := collections.Search(collection, searchMethod, text, collections.WithNamespaces([]string{namespace}), collections.WithLimit(1), collections.WithReturnText(true))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionSearchResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.SearchCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&[]string{namespace}, values[1]) {
			t.Errorf("Expected namespaces: %v, but received: %v", &[]string{namespace}, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
		if !reflect.DeepEqual(&text, values[3]) {
			t.Errorf("Expected text: %v, but received: %v", &text, values[3])
		}
		if !reflect.DeepEqual(int32(1), values[4]) {
			t.Errorf("Expected limit: %v, but received: %v", int32(1), values[4])
		}
		if !reflect.DeepEqual(true, values[5]) {
			t.Errorf("Expected returnText: %v, but received: %v", true, values[5])
		}
	}
}

func TestHostNnClassifyCollection(t *testing.T) {
	result, err := collections.NnClassify(collection, searchMethod, text, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionClassificationResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.NnClassifyCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
		if !reflect.DeepEqual(&text, values[3]) {
			t.Errorf("Expected text: %v, but received: %v", &text, values[3])
		}
	}
}

func TestHostRecomputeSearchMethod(t *testing.T) {
	result, err := collections.RecomputeSearchMethod(collection, searchMethod, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.SearchMethodMutationResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.RecomputeSearchMethodCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
	}
}

func TestHostComputeDistance(t *testing.T) {
	result, err := collections.ComputeDistance(collection, searchMethod, key1, key2, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionSearchResultObject{
		Distance: 0.0,
		Score:    1.0,
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.ComputeDistanceCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
		if !reflect.DeepEqual(&key1, values[3]) {
			t.Errorf("Expected key1: %v, but received: %v", &key1, values[3])
		}
		if !reflect.DeepEqual(&key2, values[4]) {
			t.Errorf("Expected key2: %v, but received: %v", &key2, values[4])
		}
	}
}

func TestHostGetTextFromCollection(t *testing.T) {
	result, err := collections.GetText(collection, key, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', but received: %s", result)
	}

	values := collections.GetTextCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&key, values[2]) {
			t.Errorf("Expected key: %v, but received: %v", &key, values[2])
		}
	}
}

func TestHostGetTextsFromCollection(t *testing.T) {
	result, err := collections.GetTexts(collection, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := map[string]string{
		"key1": "Hello, World!",
		"key2": "Hello, World!",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.GetTextsCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
	}
}

func TestHostGetTextsFromCollectionWithNamespace(t *testing.T) {
	result, err := collections.GetTexts(collection, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := map[string]string{
		"key1": "Hello, World!",
		"key2": "Hello, World!",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.GetTextsCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
	}
}

func TestHostGetNamespacesFromCollection(t *testing.T) {
	result, err := collections.GetNamespaces(collection)
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := []string{"namespace1", "namespace2"}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.GetNamespacesCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
	}
}

func TestHostGetVectorFromCollection(t *testing.T) {
	result, err := collections.GetVector(collection, searchMethod, key, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := []float32{0.1, 0.2, 0.3}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.GetVectorCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
		if !reflect.DeepEqual(&key, values[3]) {
			t.Errorf("Expected key: %v, but received: %v", &key, values[3])
		}
	}
}

func TestHostGetLabelsFromCollection(t *testing.T) {
	result, err := collections.GetLabels(collection, key, collections.WithNamespace(namespace))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := []string{"label1", "label2"}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.GetLabelsCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&namespace, values[1]) {
			t.Errorf("Expected namespace: %v, but received: %v", &namespace, values[1])
		}
		if !reflect.DeepEqual(&key, values[2]) {
			t.Errorf("Expected key: %v, but received: %v", &key, values[2])
		}
	}
}

func TestHostSearchByVector(t *testing.T) {
	result, err := collections.SearchByVector(collection, searchMethod, []float32{0.1, 0.2, 0.3}, collections.WithNamespaces([]string{namespace}), collections.WithLimit(1), collections.WithReturnText(true))
	if err != nil {
		t.Fatal(err.Error())
	}
	if result == nil {
		t.Fatal("Expected a result, but none was found.")
	}
	expected := &collections.CollectionSearchResult{
		Collection: "collection",
		Status:     "success",
	}
	if !reflect.DeepEqual(expected, result) {
		t.Errorf("Expected result: %v, but received: %v", expected, result)
	}

	values := collections.SearchByVectorCallStack.Pop()
	if values == nil {
		t.Error("Expected a request, but none was found.")
	} else {
		if !reflect.DeepEqual(&collection, values[0]) {
			t.Errorf("Expected collection: %v, but received: %v", &collection, values[0])
		}
		if !reflect.DeepEqual(&[]string{namespace}, values[1]) {
			t.Errorf("Expected namespaces: %v, but received: %v", &[]string{namespace}, values[1])
		}
		if !reflect.DeepEqual(&searchMethod, values[2]) {
			t.Errorf("Expected searchMethod: %v, but received: %v", &searchMethod, values[2])
		}
		if !reflect.DeepEqual(&[]float32{0.1, 0.2, 0.3}, values[3]) {
			t.Errorf("Expected vector: %v, but received: %v", &[]float32{0.1, 0.2, 0.3}, values[3])
		}
		if !reflect.DeepEqual(int32(1), values[4]) {
			t.Errorf("Expected limit: %v, but received: %v", int32(1), values[4])
		}
		if !reflect.DeepEqual(true, values[5]) {
			t.Errorf("Expected returnText: %v, but received: %v", true, values[5])
		}
	}
}

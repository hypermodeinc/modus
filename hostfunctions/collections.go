package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/collections"
	"hmruntime/logger"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32, pText uint32) uint32 {
	return hostUpsertToCollectionV2(ctx, mod, pCollectionName, 0, pKey, pText, 0)
}

func hostUpsertToCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pKeys uint32, pTexts uint32, pLabels uint32) uint32 {
	var collectionName string
	var namespace string
	var keys []string
	var texts []string
	var labels [][]string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pKeys, &keys}, param{pTexts, &texts}, param{pLabels, &labels})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled collection upsert.",
		Error:     "Error upserting to collection.",
		Detail:    fmt.Sprintf("Collection: %s, Keys: %v", collectionName, keys),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting collection upsert."
		msgs.Completed = "Completed collection upsert."
	}

	// Prepare the host function
	var mutationRes *collections.CollectionMutationResult
	fn := func() (err error) {
		mutationRes, err = collections.UpsertToCollection(ctx, collectionName, namespace, keys, texts, labels)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	return hostDeleteFromCollectionV2(ctx, mod, pCollectionName, 0, pKey)
}

func hostDeleteFromCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pKey uint32) uint32 {
	var collectionName string
	var namespace string
	var key string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pKey, &key})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled deleting from collection.",
		Error:     "Error deleting from collection.",
		Detail:    fmt.Sprintf("Collection: %s, Key: %s", collectionName, key),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting deleting from collection."
		msgs.Completed = "Completed deleting from collection."
	}

	// Prepare the host function
	var mutationRes *collections.CollectionMutationResult
	fn := func() (err error) {
		mutationRes, err = collections.DeleteFromCollection(ctx, collectionName, namespace, key)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostSearchCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32,
	pText uint32, pLimit uint32, pReturnText uint32) uint32 {
	return hostSearchCollectionV2(ctx, mod, pCollectionName, 0, pSearchMethod, pText, pLimit, pReturnText)
}

func hostSearchCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pSearchMethod uint32,
	pText uint32, pLimit uint32, pReturnText uint32) uint32 {
	var collectionName string
	var namespace string
	var searchMethod string
	var text string
	var limit int32
	var returnText bool

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pSearchMethod, &searchMethod}, param{pText, &text}, param{pLimit, &limit}, param{pReturnText, &returnText})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled searching collection.",
		Error:     "Error searching collection.",
		Detail:    fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting searching collection."
		msgs.Completed = "Completed searching collection."
	}

	// Prepare the host function
	var searchRes *collections.CollectionSearchResult
	fn := func() (err error) {
		searchRes, err = collections.SearchCollection(ctx, collectionName, namespace, searchMethod, text, limit, returnText)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, searchRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostNnClassifyCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32, pText uint32) uint32 {
	return hostNnClassifyCollectionV2(ctx, mod, pCollectionName, 0, pSearchMethod, pText)
}

func hostNnClassifyCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pSearchMethod uint32, pText uint32) uint32 {
	var collectionName string
	var namespace string
	var searchMethod string
	var text string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pSearchMethod, &searchMethod}, param{pText, &text})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled classification.",
		Error:     "Error during classification.",
		Detail:    fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting classification."
		msgs.Completed = "Completed classification."
	}

	// Prepare the host function
	var classification *collections.CollectionClassificationResult
	fn := func() (err error) {
		classification, err = collections.NnClassify(ctx, collectionName, namespace, searchMethod, text)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, classification)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostComputeDistance(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32, pId1 uint32, pId2 uint32) uint32 {
	return hostComputeDistanceV2(ctx, mod, pCollectionName, 0, pSearchMethod, pId1, pId2)
}

func hostComputeDistanceV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pSearchMethod uint32, pId1 uint32, pId2 uint32) uint32 {
	var collectionName string
	var namespace string
	var searchMethod string
	var id1 string
	var id2 string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pSearchMethod, &searchMethod}, param{pId1, &id1}, param{pId2, &id2})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled computing distance.",
		Error:     "Error computing distance.",
		Detail:    fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting computing distance."
		msgs.Completed = "Completed computing distance."
	}

	// Prepare the host function
	var resObj *collections.CollectionSearchResultObject
	fn := func() (err error) {
		resObj, err = collections.ComputeDistance(ctx, collectionName, namespace, searchMethod, id1, id2)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, resObj)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostRecomputeSearchMethod(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32) uint32 {
	return hostRecomputeSearchMethodV2(ctx, mod, pCollectionName, 0, pSearchMethod)
}

func hostRecomputeSearchMethodV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pSearchMethod uint32) uint32 {
	var collectionName string
	var namespace string
	var searchMethod string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pSearchMethod, &searchMethod})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Starting:  "Starting recomputing search method.",
		Completed: "Completed recomputing search method.",
		Cancelled: "Cancelled recomputing search method for collection.",
		Error:     "Error recomputing search method for collection.",
		Detail:    fmt.Sprintf("Collection: %s, Method: %s", collectionName, searchMethod),
	}

	// Prepare the host function
	var mutationRes *collections.SearchMethodMutationResult
	fn := func() (err error) {
		mutationRes, err = collections.RecomputeSearchMethod(ctx, mod, collectionName, namespace, searchMethod)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetTextFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	return hostGetTextFromCollectionV2(ctx, mod, pCollectionName, 0, pKey)
}

func hostGetTextFromCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32, pKey uint32) uint32 {
	var collectionName string
	var namespace string
	var key string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace}, param{pKey, &key})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled getting text from collection.",
		Error:     "Error getting text from collection.",
		Detail:    fmt.Sprintf("Collection: %s, Key: %s", collectionName, key),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting getting text from collection."
		msgs.Completed = "Completed getting text from collection."
	}

	// Prepare the host function
	var text string
	fn := func() (err error) {
		text, err = collections.GetTextFromCollection(ctx, collectionName, namespace, key)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetTextsFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32) uint32 {
	return hostGetTextsFromCollectionV2(ctx, mod, pCollectionName, 0)
}

func hostGetTextsFromCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pNamespace uint32) uint32 {
	var collectionName string
	var namespace string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pNamespace, &namespace})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled getting texts from collection.",
		Error:     "Error getting texts from collection.",
		Detail:    fmt.Sprintf("Collection: %s", collectionName),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting getting texts from collection."
		msgs.Completed = "Completed getting texts from collection."
	}

	// Prepare the host function
	var texts map[string]string
	fn := func() (err error) {
		texts, err = collections.GetTextsFromCollection(ctx, collectionName, namespace)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetNamespacesFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32) uint32 {
	var collectionName string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	// Prepare log messages
	msgs := &hostFunctionMessages{
		Cancelled: "Cancelled getting namespaces from collection.",
		Error:     "Error getting namespaces from collection.",
		Detail:    fmt.Sprintf("Collection: %s", collectionName),
	}

	// Track start/complete messages only if debug is enabled, to reduce log noise
	if utils.HypermodeDebugEnabled() {
		msgs.Starting = "Starting getting namespaces from collection."
		msgs.Completed = "Completed getting namespaces from collection."
	}

	// Prepare the host function
	var namespaces []string
	fn := func() (err error) {
		namespaces, err = collections.GetNamespacesFromCollection(ctx, collectionName)
		return err
	}

	// Call the host function
	if ok := callHostFunction(ctx, fn, msgs); !ok {
		return 0
	}

	// Write the results
	offset, err := writeResult(ctx, mod, namespaces)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

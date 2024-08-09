package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/collections"
	"hmruntime/logger"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func init() {
	// Current host functions for collections
	addHostFunction("upsertToCollection_v2", hostUpsertToCollection, withI32Params(5), withI32Result())
	addHostFunction("deleteFromCollection_v2", hostDeleteFromCollection, withI32Params(3), withI32Result())
	addHostFunction("searchCollection_v2", hostSearchCollection, withI32Params(6), withI32Result())
	addHostFunction("nnClassifyCollection_v2", hostNnClassifyCollection, withI32Params(4), withI32Result())
	addHostFunction("recomputeSearchMethod_v2", hostRecomputeSearchMethod, withI32Params(3), withI32Result())
	addHostFunction("computeDistance_v2", hostComputeDistance, withI32Params(5), withI32Result())
	addHostFunction("getTextFromCollection_v2", hostGetTextFromCollection, withI32Params(3), withI32Result())
	addHostFunction("getTextsFromCollection_v2", hostGetTextsFromCollection, withI32Params(2), withI32Result())
	addHostFunction("getNamespacesFromCollection", hostGetNamespacesFromCollection, withI32Param(), withI32Result())

	// Support functions from older SDK versions
	addHostFunction("upsertToCollection", hostUpsertToCollection, withI32Params(3), withI32Result())
	addHostFunction("deleteFromCollection", hostDeleteFromCollection, withI32Params(2), withI32Result())
	addHostFunction("searchCollection", hostSearchCollection, withI32Params(5), withI32Result())
	addHostFunction("nnClassifyCollection", hostNnClassifyCollection, withI32Params(3), withI32Result())
	addHostFunction("recomputeSearchMethod", hostRecomputeSearchMethod, withI32Params(2), withI32Result())
	addHostFunction("computeSimilarity", hostComputeDistance, withI32Params(4), withI32Result())
	addHostFunction("computeDistance", hostComputeDistance, withI32Params(4), withI32Result())
	addHostFunction("getTextFromCollection", hostGetTextFromCollection, withI32Params(2), withI32Result())
	addHostFunction("getTextsFromCollection", hostGetTextsFromCollection, withI32Param(), withI32Result())
}

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace string
	var keys, texts []string
	var labels [][]string
	if len(stack) == 3 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &keys, &texts); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace and labels)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &keys, &texts, &labels); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, mutationRes); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, key string
	if len(stack) == 2 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &key); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &key); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, mutationRes); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostSearchCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, searchMethod, text string
	var limit int32
	var returnText bool
	if len(stack) == 5 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &searchMethod, &text, &limit, &returnText); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &searchMethod, &text, &limit, &returnText); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, searchRes); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostNnClassifyCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, searchMethod, text string
	if len(stack) == 3 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &searchMethod, &text); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &searchMethod, &text); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, classification); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostComputeDistance(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, searchMethod, id1, id2 string
	if len(stack) == 4 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &searchMethod, &id1, &id2); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &searchMethod, &id1, &id2); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, resObj); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostRecomputeSearchMethod(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, searchMethod string
	if len(stack) == 2 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &searchMethod); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &searchMethod); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, mutationRes); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostGetTextFromCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace, key string
	if len(stack) == 2 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName, &key); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace, &key); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, text); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostGetTextsFromCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName, namespace string
	if len(stack) == 1 {
		// v1
		if err := readParams(ctx, mod, stack, &collectionName); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
	} else {
		// v2 (with namespace)
		if err := readParams(ctx, mod, stack, &collectionName, &namespace); err != nil {
			logger.Err(ctx, err).Msg("Error reading input parameters.")
			return
		}
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, texts); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

func hostGetNamespacesFromCollection(ctx context.Context, mod wasm.Module, stack []uint64) {

	// Read input parameters
	var collectionName string
	if err := readParams(ctx, mod, stack, &collectionName); err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return
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
		return
	}

	// Write the results
	if err := writeResults(ctx, mod, stack, namespaces); err != nil {
		logger.Err(ctx, err).Msg("Error writing results to wasm memory.")
	}
}

package hostfunctions

import (
	"context"
	"fmt"
	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/utils"
	"hmruntime/vector"

	wasm "github.com/tetratelabs/wazero/api"
)

type VectorIndexOperationResult struct {
	mutation VectorIndexMutationResult
	query    VectorIndexSearchResult
}

type VectorIndexMutationResult struct {
	status    string
	operation string
	id        string
}

type VectorIndexSearchResult struct {
	status  string
	objects []VectorIndexSearchResultObject
}

type VectorIndexSearchResultObject struct {
	id    string
	score float64
}

func getEmbedderFunctionInfo(embedder string) (functions.FunctionInfo, error) {
	info, ok := functions.Functions[embedder]
	if !ok {
		return functions.FunctionInfo{}, fmt.Errorf("embedder function not found: %s", embedder)
	}
	if len(info.Function.Parameters) > 1 {
		return functions.FunctionInfo{}, fmt.Errorf("embedder function must have only one parameter: %s", embedder)
	}
	if info.Function.Parameters[0].Type.Name != "string" {
		return functions.FunctionInfo{}, fmt.Errorf("embedder function must take a string parameter: %s", embedder)
	}
	if info.Function.ReturnType.Name != "f64[]" {
		return functions.FunctionInfo{}, fmt.Errorf("embedder function must return a float64 array: %s", embedder)
	}

	return info, nil
}

func hostUpsertToTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32, pText uint32) uint32 {
	var collection string
	var id string
	var text string

	err := readParams3(ctx, mod, pCollection, pId, pText, &collection, &id, &text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	if id == "" {
		id = utils.GenerateUUIDV7()
	}

	// Get the collection data from the manifest
	collectionData := manifestdata.Manifest.Collections[collection]

	// insert text into text index
	textIndex, err := vector.GlobalTextIndexFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding text index.")
		return 0
	}

	// compute embeddings for each search method, and insert into vector index
	for searchMethodName, searchMethod := range collectionData.SearchMethods {
		vectorIndex, err := textIndex.GetVectorIndex(searchMethodName)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting vector index.")
			return 0
		}

		embedder := searchMethod.Embedder
		info, err := getEmbedderFunctionInfo(embedder)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting embedder function info.")
			return 0
		}

		parameters := make(map[string]interface{})
		parameters[info.Function.Parameters[0].Name] = text

		result, err := functions.CallFunction(ctx, mod, info, parameters)
		if err != nil {
			logger.Err(ctx, err).Msg("Error calling function.")
			return 0
		}

		textVec := result.([]float64)
		_, err = vectorIndex.InsertVector(ctx, nil, id, textVec)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting into vector index.")
			return 0
		}
	}
	_, err = textIndex.InsertText(ctx, nil, id, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error inserting into text index.")
		return 0
	}

	output := &VectorIndexOperationResult{
		mutation: VectorIndexMutationResult{
			status:    "success",
			operation: "insert",
			id:        id,
		},
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostDeleteFromTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32) uint32 {
	var collection string
	var id string

	err := readParams2(ctx, mod, pCollection, pId, &collection, &id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	textIndex, err := vector.GlobalTextIndexFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding text index.")
		return 0
	}
	for _, vectorIndex := range textIndex.GetVectorIndexMap() {
		err = vectorIndex.DeleteVector(ctx, nil, id)
		if err != nil {
			logger.Err(ctx, err).Msg("Error deleting from index.")
			return 0
		}
	}
	err = textIndex.DeleteText(ctx, nil, id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error deleting from index.")
		return 0
	}

	output := &VectorIndexOperationResult{
		mutation: VectorIndexMutationResult{
			status:    "success",
			operation: "delete",
			id:        id,
		},
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostSearchTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pSearchMethod uint32,
	pText uint32, pLimit uint32) uint32 {
	var collection string
	var searchMethod string
	var text string
	var limit int

	err := readParams4(ctx, mod, pCollection, pSearchMethod, pText, pLimit,
		&collection, &searchMethod, &text, &limit)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	textIndex, err := vector.GlobalTextIndexFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding text index.")
		return 0
	}

	vectorIndex, err := textIndex.GetVectorIndex(searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting vector index.")
		return 0
	}

	embedder := manifestdata.Manifest.Collections[collection].SearchMethods[searchMethod].Embedder
	info, err := getEmbedderFunctionInfo(embedder)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting embedder function info.")
		return 0
	}

	parameters := make(map[string]interface{})
	parameters[info.Function.Parameters[0].Name] = text

	result, err := functions.CallFunction(ctx, mod, info, parameters)
	if err != nil {
		logger.Err(ctx, err).Msg("Error calling function.")
		return 0
	}

	textVec := result.([]float64)
	objects, err := vectorIndex.Search(ctx, nil, textVec, limit, nil)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching vector index.")
		return 0
	}

	output := &VectorIndexOperationResult{
		query: VectorIndexSearchResult{
			status:  "success",
			objects: make([]VectorIndexSearchResultObject, len(objects)),
		},
	}

	for i, object := range objects {
		output.query.objects[i] = VectorIndexSearchResultObject{
			id:    object.GetIndex(),
			score: object.GetValue(),
		}
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostRecomputeTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pSearchMethod uint32) uint32 {
	var collection string
	var searchMethod string

	err := readParams2(ctx, mod, pCollection, pSearchMethod, &collection, &searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	textIndex, err := vector.GlobalTextIndexFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding text index.")
		return 0
	}

	vectorIndex, err := textIndex.GetVectorIndex(searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting vector index.")
		return 0
	}

	embedder := manifestdata.Manifest.Collections[collection].SearchMethods[searchMethod].Embedder
	info, err := getEmbedderFunctionInfo(embedder)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting embedder function info.")
		return 0
	}

	for uuid, text := range textIndex.GetTextMap() {
		parameters := make(map[string]interface{})
		parameters[info.Function.Parameters[0].Name] = text

		result, err := functions.CallFunction(ctx, mod, info, parameters)
		if err != nil {
			logger.Err(ctx, err).Msg("Error calling function.")
			return 0
		}

		textVec := result.([]float64)
		_, err = vectorIndex.InsertVector(ctx, nil, uuid, textVec)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting into vector index.")
			return 0
		}
	}

	output := &VectorIndexOperationResult{
		mutation: VectorIndexMutationResult{
			status:    "success",
			operation: "recompute",
		},
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

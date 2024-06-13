package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/plugins"
	"hmruntime/utils"
	"hmruntime/vector"
	"hmruntime/wasmhost/module"

	wasm "github.com/tetratelabs/wazero/api"
)

type textIndexMutationResult struct {
	Collection string
	Operation  string
	Status     string
	ID         string
}

func (r *textIndexMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "TextIndexMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/TextIndexMutationResult",
	}
}

type textIndexSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []textIndexSearchResultObject
}

func (r *textIndexSearchResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "TextIndexSearchResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/TextIndexSearchResult",
	}
}

type textIndexSearchResultObject struct {
	ID    string
	Text  string
	Score float64
}

func (r *textIndexSearchResultObject) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "TextIndexSearchResultObject",
		Path: "~lib/@hypermode/functions-as/assembly/collections/TextIndexSearchResultObject",
	}
}

func hostUpsertToTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32, pText uint32) uint32 {
	var collection string
	var id string
	var text string

	var err error
	if pId != 0 {
		err = readParams3(ctx, mod, pCollection, pId, pText, &collection, &id, &text)
	} else {
		err = readParams2(ctx, mod, pCollection, pText, &collection, &text)
	}
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
		err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
		if err != nil {
			logger.Err(ctx, err).Msg("Error verifying function signature.")
			return 0
		}

		result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			logger.Err(ctx, err).Msg("Error calling function.")
			return 0
		}

		resultArr := result.([]interface{})

		textVec := make([]float64, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(float64)
		}

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

	output := textIndexMutationResult{
		Collection: collection,
		Status:     "success",
		Operation:  "insert",
		ID:         id,
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

	output := textIndexMutationResult{
		Collection: collection,
		Status:     "success",
		Operation:  "delete",
		ID:         id,
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostSearchTextIndex(ctx context.Context, mod wasm.Module, pCollection uint32, pSearchMethod uint32,
	pText uint32, pLimit uint32, pReturnText uint32) uint32 {
	var collection string
	var searchMethod string
	var text string
	var limit int32
	var returnText bool

	err := readParams5(ctx, mod, pCollection, pSearchMethod, pText, pLimit, pReturnText,
		&collection, &searchMethod, &text, &limit, &returnText)
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
	err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")
		return 0
	}

	result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error calling function.")
		return 0
	}

	resultArr := result.([]interface{})
	textVec := make([]float64, len(resultArr))
	for i, val := range resultArr {
		textVec[i] = val.(float64)
	}

	objects, err := vectorIndex.Search(ctx, nil, textVec, int(limit), nil)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching vector index.")
		return 0
	}

	output := textIndexSearchResult{
		Collection: collection,
		Status:     "success",
		Objects:    make([]textIndexSearchResultObject, len(objects)),
	}

	for i, object := range objects {
		if returnText {
			text, err := textIndex.GetText(ctx, nil, object.GetIndex())
			if err != nil {
				logger.Err(ctx, err).Msg("Error getting text.")
				return 0
			}
			output.Objects[i] = textIndexSearchResultObject{
				ID:    object.GetIndex(),
				Text:  text,
				Score: object.GetValue(),
			}
		} else {
			output.Objects[i] = textIndexSearchResultObject{
				ID:    object.GetIndex(),
				Score: object.GetValue(),
			}
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
	err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")
		return 0
	}

	err = vector.ProcessTextMapWithModule(ctx, mod, textIndex, embedder, vectorIndex)
	if err != nil {
		logger.Err(ctx, err).Msg("Error processing text map.")
		return 0
	}

	output := textIndexMutationResult{
		Collection: collection,
		Status:     "success",
		Operation:  "recompute",
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetText(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32) uint32 {
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

	text, err := textIndex.GetText(ctx, nil, id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting text.")
		return 0
	}

	offset, err := writeResult(ctx, mod, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

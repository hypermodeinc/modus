package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/collections"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/plugins"
	"hmruntime/utils"
	"hmruntime/wasmhost/module"

	wasm "github.com/tetratelabs/wazero/api"
)

type textIndexMutationResult struct {
	Collection string
	Operation  string
	Status     string
	ID         string
	Error      string
}

func (r *textIndexMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	}
}

type textIndexSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []textIndexSearchResultObject
	Error        string
}

func (r *textIndexSearchResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	}
}

type textIndexSearchResultObject struct {
	ID    string
	Text  string
	Score float64
}

func (r *textIndexSearchResultObject) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResultObject",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	}
}

func WriteCollectionMutationResultOffset(ctx context.Context, mod wasm.Module, collection, operation, status, id, error string) (uint32, error) {
	output := textIndexMutationResult{
		Collection: collection,
		Operation:  operation,
		Status:     status,
		ID:         id,
		Error:      error,
	}

	return writeResult(ctx, mod, output)
}

func WriteCollectionSearchResultOffset(ctx context.Context, mod wasm.Module, collection, searchMethod, status string,
	objects []textIndexSearchResultObject, error string) (uint32, error) {

	output := textIndexSearchResult{
		Collection:   collection,
		SearchMethod: searchMethod,
		Status:       status,
		Objects:      objects,
		Error:        error,
	}

	return writeResult(ctx, mod, output)
}

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32, pText uint32) uint32 {
	var collection string
	var id string
	var text string

	err := readParams3(ctx, mod, pCollection, pId, pText, &collection, &id, &text)

	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	if id == "" {
		id = utils.GenerateUUIDV7()
	}

	// Get the collection data from the manifest
	collectionData := manifestdata.Manifest.Collections[collection]

	// insert text into text index
	textIndex, err := collections.GlobalCollectionFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collection.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error finding collection: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	// compute embeddings for each search method, and insert into vector index
	for searchMethodName, searchMethod := range collectionData.SearchMethods {
		vectorIndex, err := textIndex.GetVectorIndex(searchMethodName)
		if err != nil {
			logger.Err(ctx, err).Msg("Error finding search method.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error finding search method: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		embedder := searchMethod.Embedder
		err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
		if err != nil {
			logger.Err(ctx, err).Msg("Error verifying function signature.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error verifying function signature: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			logger.Err(ctx, err).Msg("Error calling function.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error calling function: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		resultArr := result.([]interface{})

		textVec := make([]float64, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(float64)
		}

		_, err = vectorIndex.InsertVector(ctx, id, textVec)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting into vector index.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error inserting into vector index: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}
	}
	_, err = textIndex.InsertText(ctx, id, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error inserting into text index.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "error", "", fmt.Sprintf("Error inserting into text index: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "upsert", "success", id, "")
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, pCollection uint32, pId uint32) uint32 {
	var collection string
	var id string

	err := readParams2(ctx, mod, pCollection, pId, &collection, &id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "delete", "error", "", fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	textIndex, err := collections.GlobalCollectionFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collection.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "delete", "error", "", fmt.Sprintf("Error finding collection: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}
	for _, vectorIndex := range textIndex.GetVectorIndexMap() {
		err = vectorIndex.DeleteVector(ctx, id)
		if err != nil {
			logger.Err(ctx, err).Msg("Error deleting from index.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "delete", "error", "", fmt.Sprintf("Error deleting from index: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}
	}
	err = textIndex.DeleteText(ctx, id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error deleting from index.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "delete", "error", "", fmt.Sprintf("Error deleting from index: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "delete", "success", id, "")
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostSearchCollection(ctx context.Context, mod wasm.Module, pCollection uint32, pSearchMethod uint32,
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

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	textIndex, err := collections.GlobalCollectionFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collection.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error finding collection: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	vectorIndex, err := textIndex.GetVectorIndex(searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding search method.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error finding search method: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	embedder := manifestdata.Manifest.Collections[collection].SearchMethods[searchMethod].Embedder
	err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error verifying function signature: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
	if err != nil {
		logger.Err(ctx, err).Msg("Error calling function.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error calling function: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	resultArr := result.([]interface{})
	textVec := make([]float64, len(resultArr))
	for i, val := range resultArr {
		textVec[i] = val.(float64)
	}

	objects, err := vectorIndex.Search(ctx, textVec, int(limit), nil)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching vector index.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error searching vector index: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	output := textIndexSearchResult{
		Collection:   collection,
		SearchMethod: searchMethod,
		Status:       "success",
		Objects:      make([]textIndexSearchResultObject, len(objects)),
	}

	for i, object := range objects {
		if returnText {
			text, err := textIndex.GetText(ctx, object.GetIndex())
			if err != nil {
				logger.Err(ctx, err).Msg("Error getting text.")

				offset, err := WriteCollectionSearchResultOffset(ctx, mod, collection, searchMethod, "error", nil, fmt.Sprintf("Error getting text: %s", err.Error()))
				if err != nil {
					logger.Err(ctx, err).Msg("Error writing result.")
				}
				return offset
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
	}

	return offset
}

func hostRecomputeSearchMethod(ctx context.Context, mod wasm.Module, pCollection uint32, pSearchMethod uint32) uint32 {
	var collection string
	var searchMethod string

	err := readParams2(ctx, mod, pCollection, pSearchMethod, &collection, &searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "error", "", fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	textIndex, err := collections.GlobalCollectionFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collection.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "error", "", fmt.Sprintf("Error finding collection: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	vectorIndex, err := textIndex.GetVectorIndex(searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding search method.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "error", "", fmt.Sprintf("Error finding search method: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	embedder := manifestdata.Manifest.Collections[collection].SearchMethods[searchMethod].Embedder
	err = module.VerifyFunctionSignature(embedder, "string", "f64[]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "error", "", fmt.Sprintf("Error verifying function signature: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	err = collections.ProcessTextMapWithModule(ctx, mod, textIndex, embedder, vectorIndex)
	if err != nil {
		logger.Err(ctx, err).Msg("Error processing text map.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "error", "", fmt.Sprintf("Error processing text map: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	offset, err := WriteCollectionMutationResultOffset(ctx, mod, collection, "recompute", "success", "", "")
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

	textIndex, err := collections.GlobalCollectionFactory.Find(collection)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collection.")
		return 0
	}

	text, err := textIndex.GetText(ctx, id)
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

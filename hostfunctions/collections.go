package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/collections"
	collection_utils "hmruntime/collections/utils"
	"hmruntime/functions"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/plugins"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	wasm "github.com/tetratelabs/wazero/api"
)

type collectionMutationResult struct {
	Collection string
	Operation  string
	Status     string
	Keys       []string
	Error      string
}

func (r *collectionMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionMutationResult",
	}
}

type searchMethodMutationResult struct {
	Collection   string
	SearchMethod string
	Operation    string
	Status       string
	Error        string
}

func (r *searchMethodMutationResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "SearchMethodMutationResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/SearchMethodMutationResult",
	}
}

type collectionSearchResult struct {
	Collection   string
	SearchMethod string
	Status       string
	Objects      []collectionSearchResultObject
	Error        string
}

func (r *collectionSearchResult) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResult",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResult",
	}
}

type collectionSearchResultObject struct {
	Key      string
	Text     string
	Distance float64
	Score    float64
}

func (r *collectionSearchResultObject) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "CollectionSearchResultObject",
		Path: "~lib/@hypermode/functions-as/assembly/collections/CollectionSearchResultObject",
	}
}

func WriteCollectionMutationResultOffset(ctx context.Context, mod wasm.Module, collectionName, operation, status string, keys []string, error string) (uint32, error) {
	if keys == nil {
		keys = []string{}
	}
	output := collectionMutationResult{
		Collection: collectionName,
		Operation:  operation,
		Status:     status,
		Keys:       keys,
		Error:      error,
	}

	return writeResult(ctx, mod, output)
}

func WriteSearchMethodMutationResultOffset(ctx context.Context, mod wasm.Module, collectionName, searchMethod, operation, status, error string) (uint32, error) {
	output := searchMethodMutationResult{
		Collection:   collectionName,
		SearchMethod: searchMethod,
		Operation:    operation,
		Status:       status,
		Error:        error,
	}

	return writeResult(ctx, mod, output)
}

func WriteCollectionSearchResultOffset(ctx context.Context, mod wasm.Module, collectionName, searchMethod, status string,
	objects []collectionSearchResultObject, error string) (uint32, error) {

	output := collectionSearchResult{
		Collection:   collectionName,
		SearchMethod: searchMethod,
		Status:       status,
		Objects:      objects,
		Error:        error,
	}

	return writeResult(ctx, mod, output)
}

func WriteCollectionSearchResultObjectOffset(ctx context.Context, mod wasm.Module, key, text string, distance, score float64) (uint32, error) {
	output := collectionSearchResultObject{
		Key:      key,
		Text:     text,
		Distance: distance,
		Score:    score,
	}

	return writeResult(ctx, mod, output)
}

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKeys uint32, pTexts uint32) uint32 {
	var collectionName string
	var keys []string
	var texts []string

	err := readParams3(ctx, mod, pCollectionName, pKeys, pTexts, &collectionName, &keys, &texts)

	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset

	}

	// Get the collectionName data from the manifest
	collectionData := manifestdata.Manifest.Collections[collectionName]

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error finding collectionName: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	if len(keys) == 0 {
		keys = make([]string, len(texts))
		for i := range keys {
			keys[i] = utils.GenerateUUIDV7()
		}
	}

	err = collection.InsertTexts(ctx, keys, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error inserting into collection.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error inserting into collection: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	// compute embeddings for each search method, and insert into vector index
	for searchMethodName, searchMethod := range collectionData.SearchMethods {
		vectorIndex, err := collection.GetVectorIndex(ctx, searchMethodName)
		if err != nil {
			logger.Err(ctx, err).Msg("Error finding search method.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error finding search method: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		embedder := searchMethod.Embedder

		info, err := functions.GetFunctionInfo(embedder)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting function info.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error getting function info: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
		if err != nil {
			logger.Err(ctx, err).Msg("Error verifying function signature.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error verifying function signature: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
		if err != nil {
			logger.Err(ctx, err).Msg("Error calling function.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error calling function: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		result := executionInfo.Result

		textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
		if err != nil {
			logger.Err(ctx, err).Msg("Error converting to float32.")
		}

		if len(textVecs) != len(texts) {
			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, "length of vectors does not match length of texts.")
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}

		ids := make([]int64, len(keys))
		for i := range textVecs {
			key := keys[i]

			id, err := collection.GetExternalId(ctx, key)
			if err != nil {
				logger.Err(ctx, err).Msg("Error getting external id.")
			}
			ids[i] = id
		}

		err = vectorIndex.InsertVectors(ctx, ids, textVecs)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting into vector index.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "error", nil, fmt.Sprintf("Error inserting into vector index: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset

		}
	}

	offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "upsert", "success", keys, "")
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	var collectionName string
	var key string

	err := readParams2(ctx, mod, pCollectionName, pKey, &collectionName, &key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "error", nil, fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "error", nil, fmt.Sprintf("Error finding collectionName: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}
	textId, err := collection.GetExternalId(ctx, key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting external id.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "error", nil, fmt.Sprintf("Error getting external id: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}
	for _, vectorIndex := range collection.GetVectorIndexMap() {
		err = vectorIndex.DeleteVector(ctx, textId, key)
		if err != nil {
			logger.Err(ctx, err).Msg("Error deleting from index.")

			offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "error", nil, fmt.Sprintf("Error deleting from index: %s", err.Error()))
			if err != nil {
				logger.Err(ctx, err).Msg("Error writing result.")
			}
			return offset
		}
	}
	err = collection.DeleteText(ctx, key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error deleting from index.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "error", nil, fmt.Sprintf("Error deleting from index: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	keys := []string{key}

	offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "delete", "success", keys, "")
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostSearchCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32,
	pText uint32, pLimit uint32, pReturnText uint32) uint32 {
	var collectionName string
	var searchMethod string
	var text string
	var limit int32
	var returnText bool

	err := readParams5(ctx, mod, pCollectionName, pSearchMethod, pText, pLimit, pReturnText,
		&collectionName, &searchMethod, &text, &limit, &returnText)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error finding collectionName: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding search method.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error finding search method: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	embedder := manifestdata.Manifest.Collections[collectionName].SearchMethods[searchMethod].Embedder

	info, err := functions.GetFunctionInfo(embedder)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting function info.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error getting function info: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error verifying function signature: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	texts := []string{text}

	executionInfo, err := wasmhost.CallFunction(ctx, embedder, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error calling function.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error calling function: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	result := executionInfo.Result

	textVecs, err := collection_utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error converting to float32.")
	}

	if len(textVecs) == 0 {
		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, "no vectors returned from embedder.")
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	objects, err := vectorIndex.Search(ctx, textVecs[0], int(limit), nil)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching vector index.")

		offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error searching vector index: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	output := collectionSearchResult{
		Collection:   collectionName,
		SearchMethod: searchMethod,
		Status:       "success",
		Objects:      make([]collectionSearchResultObject, len(objects)),
	}

	for i, object := range objects {
		if returnText {
			text, err := collection.GetText(ctx, object.GetIndex())
			if err != nil {
				logger.Err(ctx, err).Msg("Error getting text.")

				offset, err := WriteCollectionSearchResultOffset(ctx, mod, collectionName, searchMethod, "error", nil, fmt.Sprintf("Error getting text: %s", err.Error()))
				if err != nil {
					logger.Err(ctx, err).Msg("Error writing result.")
				}
				return offset
			}
			output.Objects[i] = collectionSearchResultObject{
				Key:      object.GetIndex(),
				Text:     text,
				Distance: object.GetValue(),
				Score:    1 - object.GetValue(),
			}
		} else {
			output.Objects[i] = collectionSearchResultObject{
				Key:      object.GetIndex(),
				Distance: object.GetValue(),
				Score:    1 - object.GetValue(),
			}
		}
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostComputeDistance(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32, pId1 uint32, pId2 uint32) uint32 {
	var collectionName string
	var searchMethod string
	var id1 string
	var id2 string

	err := readParams4(ctx, mod, pCollectionName, pSearchMethod, pId1, pId2, &collectionName, &searchMethod, &id1, &id2)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")
		return 0
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding search method.")
		return 0
	}

	vec1, err := vectorIndex.GetVector(ctx, id1)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting vector.")
		return 0
	}

	vec2, err := vectorIndex.GetVector(ctx, id2)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting vector.")
		return 0
	}

	similarity, err := collection_utils.CosineDistance(vec1, vec2)
	if err != nil {
		logger.Err(ctx, err).Msg("Error computing similarity.")
		return 0
	}

	output, nil := WriteCollectionSearchResultObjectOffset(ctx, mod, "", "", similarity, 1-similarity)
	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
	}

	return offset
}

func hostRecomputeSearchMethod(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32) uint32 {
	var collectionName string
	var searchMethod string

	err := readParams2(ctx, mod, pCollectionName, pSearchMethod, &collectionName, &searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error reading input parameters: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error finding collectionName: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	vectorIndex, err := collection.GetVectorIndex(ctx, searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding search method.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error finding search method: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	embedder := manifestdata.Manifest.Collections[collectionName].SearchMethods[searchMethod].Embedder

	info, err := functions.GetFunctionInfo(embedder)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting function info.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error getting function info: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	err = functions.VerifyFunctionSignature(info, "string[]", "f32[][]")
	if err != nil {
		logger.Err(ctx, err).Msg("Error verifying function signature.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error verifying function signature: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	err = collections.ProcessTextMapWithModule(ctx, mod, collection, embedder, vectorIndex)
	if err != nil {
		logger.Err(ctx, err).Msg("Error processing text map.")

		offset, err := WriteCollectionMutationResultOffset(ctx, mod, collectionName, "recompute", "error", nil, fmt.Sprintf("Error processing text map: %s", err.Error()))
		if err != nil {
			logger.Err(ctx, err).Msg("Error writing result.")
		}
		return offset
	}

	offset, err := WriteSearchMethodMutationResultOffset(ctx, mod, collectionName, searchMethod, "recompute", "success", "")
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetTextFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	var collectionName string
	var key string

	err := readParams2(ctx, mod, pCollectionName, pKey, &collectionName, &key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")
		return 0
	}

	text, err := collection.GetText(ctx, key)
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

func hostGetTextsFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32) uint32 {
	var collectionName string

	err := readParam(ctx, mod, pCollectionName, &collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	collection, err := collections.GlobalCollectionFactory.Find(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding collectionName.")
		return 0
	}

	textMap, err := collection.GetTextMap(ctx)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting text map.")
		return 0
	}

	offset, err := writeResult(ctx, mod, textMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

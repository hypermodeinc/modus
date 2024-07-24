package hostfunctions

import (
	"context"

	"hmruntime/collections"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKeys uint32, pTexts uint32) uint32 {
	var collectionName string
	var keys []string
	var texts []string

	err := readParams3(ctx, mod, pCollectionName, pKeys, pTexts, &collectionName, &keys, &texts)

	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")

		return 0

	}

	mutationRes, err := collections.UpsertToCollection(ctx, collectionName, keys, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error upserting to collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset

}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	var collectionName string
	var key string

	err := readParams2(ctx, mod, pCollectionName, pKey, &collectionName, &key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	mutationRes, err := collections.DeleteFromCollection(ctx, collectionName, key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error deleting from collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
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
		return 0
	}

	searchRes, err := collections.SearchCollection(ctx, collectionName, searchMethod, text, limit, returnText)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *searchRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
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

	resObj, err := collections.ComputeDistance(ctx, collectionName, searchMethod, id1, id2)
	if err != nil {
		logger.Err(ctx, err).Msg("Error computing distance.")
		return 0
	}

	offset, err := writeResult(ctx, mod, resObj)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset

}

func hostRecomputeSearchMethod(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32) uint32 {
	var collectionName string
	var searchMethod string

	err := readParams2(ctx, mod, pCollectionName, pSearchMethod, &collectionName, &searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	mutationRes, err := collections.RecomputeSearchMethod(ctx, mod, collectionName, searchMethod)
	if err != nil {
		logger.Err(ctx, err).Msg("Error recompute search method.")
		return 0
	}

	offset, err := writeResult(ctx, mod, *mutationRes)
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

	text, err := collections.GetTextFromCollection(ctx, collectionName, key)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting text from collection.")
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

	texts, err := collections.GetTextsFromCollection(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting texts from collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

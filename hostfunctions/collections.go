package hostfunctions

import (
	"context"

	"hmruntime/collections"
	"hmruntime/logger"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostUpsertToCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32, pText uint32) uint32 {
	return hostUpsertToCollectionV2(ctx, mod, pCollectionName, pKey, pText, 0)
}

func hostUpsertToCollectionV2(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKeys uint32, pTexts uint32, pLabels uint32) uint32 {
	var collectionName string
	var keys []string
	var texts []string
	var labels [][]string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pKeys, &keys}, param{pTexts, &texts}, param{pLabels, &labels})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	mutationRes, err := collections.UpsertToCollection(ctx, collectionName, keys, texts, labels)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error upserting to collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset

}

func hostDeleteFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	var collectionName string
	var key string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pKey, &key})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	mutationRes, err := collections.DeleteFromCollection(ctx, collectionName, key)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error deleting from collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, mutationRes)
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

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pSearchMethod, &searchMethod}, param{pText, &text}, param{pLimit, &limit}, param{pReturnText, &returnText})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	searchRes, err := collections.SearchCollection(ctx, collectionName, searchMethod, text, limit, returnText)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error searching collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, searchRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostNnClassifyCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pSearchMethod uint32, pText uint32) uint32 {
	var collectionName string
	var searchMethod string
	var text string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pSearchMethod, &searchMethod}, param{pText, &text})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	classification, err := collections.NnClassify(ctx, collectionName, searchMethod, text)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error classifying.")
		return 0
	}

	offset, err := writeResult(ctx, mod, classification)
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

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pSearchMethod, &searchMethod}, param{pId1, &id1}, param{pId2, &id2})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	resObj, err := collections.ComputeDistance(ctx, collectionName, searchMethod, id1, id2)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error computing distance.")
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

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pSearchMethod, &searchMethod})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	mutationRes, err := collections.RecomputeSearchMethod(ctx, mod, collectionName, searchMethod)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error recompute search method.")
		return 0
	}

	offset, err := writeResult(ctx, mod, mutationRes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostGetTextFromCollection(ctx context.Context, mod wasm.Module, pCollectionName uint32, pKey uint32) uint32 {
	var collectionName string
	var key string

	err := readParams(ctx, mod, param{pCollectionName, &collectionName}, param{pKey, &key})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	text, err := collections.GetTextFromCollection(ctx, collectionName, key)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error getting text from collection.")
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

	err := readParams(ctx, mod, param{pCollectionName, &collectionName})
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	texts, err := collections.GetTextsFromCollection(ctx, collectionName)
	if err != nil {
		logger.Err(ctx, err).
			Bool("user_visible", true).
			Msg("Error getting texts from collection.")
		return 0
	}

	offset, err := writeResult(ctx, mod, texts)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

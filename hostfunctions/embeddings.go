/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/utils"
	"hmruntime/vector"
	"hmruntime/vector/in_mem"

	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	err := readParams2(ctx, mod, pModelName, pSentenceMap, &modelName, &sentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName, manifest.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	result, err := models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Error(ctx).Msg("Empty result returned from model.")
		return 0
	}

	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}

func hostEmbedAndIndex(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32, pCollectionName uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	var collectionName string
	err := readParams3(ctx, mod, pModelName, pSentenceMap, pCollectionName, &modelName, &sentenceMap, &collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName, manifest.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	result, err := models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	index, err := vector.GlobalIndexFactory.Find(collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding index.")
		return 0
	}
	if index == nil {
		// create index
		index, err = vector.GlobalIndexFactory.Create(collectionName, nil, nil, &in_mem.InMemBruteForceIndex{})
		if err != nil {
			logger.Err(ctx, err).Msg("Error creating index.")
			return 0
		}
	}
	// // convert index to InMemIndex
	// bfIndex := index.(*in_mem.InMemBruteForceIndex)

	for k, v := range result {
		// generate random uint64
		uid := utils.NextUint64()
		_, err := index.Insert(ctx, nil, uid, v)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting into index.")
			return 0
		}
		err = index.InsertDataNode(ctx, nil, uid, k)
		if err != nil {
			logger.Err(ctx, err).Msg("Error inserting data node into index.")
			return 0
		}
	}

	offset, err := writeResult(ctx, mod, result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset

}

func hostEmbedAndSearchIndex(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32, pCollectionName uint32, pNumResults uint32) uint32 {

	var modelName string
	var sentenceMap map[string]string
	var collectionName string
	var numResults int
	err := readParams4(ctx, mod, pModelName, pSentenceMap, pCollectionName, pNumResults, &modelName, &sentenceMap, &collectionName, &numResults)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName, manifest.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	result, err := models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	index, err := vector.GlobalIndexFactory.Find(collectionName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding index.")
		return 0
	}

	if index == nil {
		logger.Err(ctx, err).Msg("Index not found.")
		return 0
	}

	// convert index to InMemIndex
	bfIndex := index.(*in_mem.InMemBruteForceIndex)

	// search the index
	nns := make(map[string][]uint64)
	for k, v := range result {
		uids, err := bfIndex.Search(ctx, nil, v, numResults, nil)
		if err != nil {
			logger.Err(ctx, err).Msg("Error searching index.")
			return 0
		}
		nns[k] = uids
	}

	// get the data nodes
	dataNodes := make(map[string][]string)
	for k, v := range nns {
		for _, uid := range v {
			dataNodes[k] = append(dataNodes[k], bfIndex.GetDataNode(ctx, nil, uid))
		}
	}

	offset, err := writeResult(ctx, mod, dataNodes)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing classification result.")
		return 0
	}

	return offset
}

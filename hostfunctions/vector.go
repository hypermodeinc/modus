package hostfunctions

import (
	"context"
	"hmruntime/logger"
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
	vector    []float64
}

type VectorIndexSearchResult struct {
	status  string
	objects []VectorIndexSearchResultObject
}

type VectorIndexSearchResultObject struct {
	id    string
	score float64
}

func hostInsertToVectorIndex(ctx context.Context, mod wasm.Module, pCollectionName uint32, pVectorIndexName uint32, pId uint32, pVector uint32) uint32 {
	var collectionName string
	var vectorIndexName string
	var id string
	var vec []float64

	err := readParams4(ctx, mod, pCollectionName, pVectorIndexName, pId, pVector, &collectionName, &vectorIndexName, &id, &vec)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}
	indexName := collectionName + "." + vectorIndexName
	index, err := vector.GlobalIndexFactory.Find(indexName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding index.")
		return 0
	}
	if index == nil {
		logger.Err(ctx, err).Msg("Index not found.")
		return 0
	}
	_, err = index.Insert(ctx, nil, id, vec)
	if err != nil {
		logger.Err(ctx, err).Msg("Error inserting into index.")
		return 0
	}

	output := &VectorIndexOperationResult{
		mutation: VectorIndexMutationResult{
			status:    "success",
			operation: "insert",
			id:        id,
			vector:    vec,
		},
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostSearchVectorIndex(ctx context.Context, mod wasm.Module, pCollectionName uint32, pVectorIndexName uint32, pVector uint32, pLimit uint32) uint32 {
	var collectionName string
	var vectorIndexName string
	var vec []float64
	var limit int

	err := readParams4(ctx, mod, pCollectionName, pVectorIndexName, pVector, pLimit, &collectionName, &vectorIndexName, &vec, &limit)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}
	indexName := collectionName + "." + vectorIndexName
	index, err := vector.GlobalIndexFactory.Find(indexName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding index.")
		return 0
	}
	if index == nil {
		logger.Err(ctx, err).Msg("Index not found.")
		return 0
	}
	uids, err := index.Search(ctx, nil, vec, limit, nil)
	if err != nil {
		logger.Err(ctx, err).Msg("Error searching index.")
		return 0
	}
	objects := make([]VectorIndexSearchResultObject, 0)
	for _, uid := range uids {
		objects = append(objects, VectorIndexSearchResultObject{
			id:    uid,
			score: 0.0,
		})
	}
	output := &VectorIndexOperationResult{
		query: VectorIndexSearchResult{
			status:  "success",
			objects: objects,
		},
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result.")
		return 0
	}

	return offset
}

func hostDeleteFromVectorIndex(ctx context.Context, mod wasm.Module, pCollectionName uint32, pVectorIndexName uint32, pId uint32) uint32 {
	var collectionName string
	var vectorIndexName string
	var id string

	err := readParams3(ctx, mod, pCollectionName, pVectorIndexName, pId, &collectionName, &vectorIndexName, &id)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	indexName := collectionName + "." + vectorIndexName
	index, err := vector.GlobalIndexFactory.Find(indexName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error finding index.")
		return 0
	}
	if index == nil {
		logger.Err(ctx, err).Msg("Index not found.")
		return 0
	}
	err = index.Delete(ctx, nil, id)
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

/*
 * Copyright 2024 Hypermode, Inc.
 */

package collections

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/hypermodeinc/modus/runtime/collections/in_mem"
	"github.com/hypermodeinc/modus/runtime/collections/in_mem/sequential"
	"github.com/hypermodeinc/modus/runtime/collections/index"
	"github.com/hypermodeinc/modus/runtime/collections/index/interfaces"
	"github.com/hypermodeinc/modus/runtime/collections/utils"
	"github.com/hypermodeinc/modus/runtime/db"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/manifestdata"
	"github.com/hypermodeinc/modus/runtime/wasmhost"

	"github.com/hypermodeAI/manifest"
)

const batchSize = 25

func batchInsertVectorsToMemory(ctx context.Context, vectorIndex interfaces.VectorIndex, textIds, vectorIds []int64, keys []string, vecs [][]float32) error {
	if len(vectorIds) != len(vecs) || len(keys) != len(vecs) || len(textIds) != len(vecs) {
		return errors.New("mismatch in vectors, keys, and textIds")
	}
	for i := 0; i < len(textIds); i += batchSize {
		end := i + batchSize
		if end > len(textIds) {
			end = len(textIds)
		}
		textIdsBatch := textIds[i:end]
		vectorIdsBatch := vectorIds[i:end]
		keysBatch := keys[i:end]
		vecsBatch := vecs[i:end]

		err := vectorIndex.InsertVectorsToMemory(ctx, textIdsBatch, vectorIdsBatch, keysBatch, vecsBatch)
		if err != nil {
			return err
		}
	}
	return nil
}

func cleanAndProcessManifest(ctx context.Context) error {
	deleteIndexesNotInManifest(ctx, manifestdata.GetManifest())
	processManifestCollections(ctx, manifestdata.GetManifest())
	return nil
}

func createIndexObject(searchMethod manifest.SearchMethodInfo, searchMethodName string) (*interfaces.VectorIndexWrapper, error) {
	vectorIndex := &interfaces.VectorIndexWrapper{}
	switch searchMethod.Index.Type {
	case interfaces.SequentialManifestType:
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, searchMethod.Embedder)
	case interfaces.HnswManifestType:
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, searchMethod.Embedder)
		// // TODO: hnsw currently broken, auto-sync is not working, it keeps embedding forever, even though it has correctly indexed. for now, set it to sequential. fix in future PR
		// vectorIndex.Type = hnsw.HnswVectorIndexType
		// vectorIndex.VectorIndex = hnsw.NewHnswVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
	case "":
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, searchMethod.Embedder)
	default:
		return nil, fmt.Errorf("Unknown index type: %s", searchMethod.Index.Type)
	}

	return vectorIndex, nil
}

func deleteIndexesNotInManifest(ctx context.Context, man *manifest.HypermodeManifest) {
	for collectionName := range globalNamespaceManager.getNamespaceCollectionFactoryMap() {
		if _, ok := man.Collections[collectionName]; !ok {
			err := globalNamespaceManager.removeCollection(collectionName)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to remove collection.")
			}
			continue
		}
		col, err := globalNamespaceManager.findCollection(collectionName)
		if err != nil {
			logger.Err(ctx, err).
				Str("collection_name", collectionName).
				Msg("Failed to find collection.")
			continue
		}
		for _, collNs := range col.collectionNamespaceMap {
			vectorIndexMap := collNs.GetVectorIndexMap()
			if vectorIndexMap == nil {
				continue
			}
			for searchMethodName := range vectorIndexMap {
				err := deleteVectorIndexesNotInManifest(ctx, man, collNs, collectionName, searchMethodName)
				if err != nil {
					logger.Err(ctx, err).
						Str("index_name", searchMethodName).
						Msg("Failed to delete vector index.")
					continue
				}
			}
		}
	}
}

func deleteVectorIndexesNotInManifest(ctx context.Context, man *manifest.HypermodeManifest, col interfaces.CollectionNamespace, collectionName, searchMethodName string) error {
	_, ok := man.Collections[collectionName].SearchMethods[searchMethodName]
	if !ok {
		err := col.DeleteVectorIndex(ctx, searchMethodName)
		if err != nil {
			return err
		}
	}
	return nil
}

func processManifestCollections(ctx context.Context, man *manifest.HypermodeManifest) {
	for collectionName, collectionInfo := range man.Collections {
		col, err := globalNamespaceManager.findCollection(collectionName)
		if err == errCollectionNotFound {
			col, err = globalNamespaceManager.createCollection(collectionName, newCollection())
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to create collection.")
			}
			// forces all users to use in-memory index for now
			// TODO implement other types of indexes based on manifest info
			// fetch all tenants and create a collection for each tenant
			namespaces, err := db.GetUniqueNamespaces(ctx, collectionName)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to get unique namespaces.")
			}
			if !slices.Contains(namespaces, in_mem.DefaultNamespace) {
				namespaces = append(namespaces, in_mem.DefaultNamespace)
			}

			for _, namespace := range namespaces {
				_, err := col.createCollectionNamespace(namespace, in_mem.NewCollectionNamespace(collectionName, namespace))
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", collectionName).
						Str("namespace", namespace).
						Msg("Failed to create collection namespace.")
				}
			}
		}
		for _, collNs := range col.collectionNamespaceMap {
			for searchMethodName, searchMethod := range collectionInfo.SearchMethods {
				vi, err := collNs.GetVectorIndex(ctx, searchMethodName)

				// if the index does not exist, create it
				if err != nil {
					if err != index.ErrVectorIndexNotFound {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to get vector index.")
					} else {
						err := setIndex(ctx, collNs, searchMethod, searchMethodName)
						if err != nil {
							logger.Err(ctx, err).
								Str("index_name", searchMethodName).
								Msg("Failed to set vector index.")
						}
					}
				} else if vi != nil && vi.Type != searchMethod.Index.Type {
					if err := collNs.DeleteVectorIndex(ctx, searchMethodName); err != nil {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to delete vector index.")
					} else {
						err := setIndex(ctx, collNs, searchMethod, searchMethodName)
						if err != nil {
							logger.Err(ctx, err).
								Str("index_name", searchMethodName).
								Msg("Failed to set vector index.")
						}
					}
				} else if vi.GetEmbedderName() != searchMethod.Embedder {
					//TODO: figure out what to actually do if the embedder is different, for now just updating the name
					// what if the user changes the internals? -> they want us to reindex, model might have diff dimensions
					// but what if they just changed the name? -> they want us to just update the name. but we cant know that
					// imo we should just update the name and let the user reindex if they want to
					if err := vi.SetEmbedderName(searchMethod.Embedder); err != nil {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to update vector index.")
					}
				}
			}
		}
	}
}

func processTexts(ctx context.Context, col interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex, keys []string, texts []string) error {
	if len(keys) != len(texts) {
		return fmt.Errorf("mismatch in keys and texts")
	}
	for i := 0; i < len(keys); i += batchSize {
		end := i + batchSize
		if end > len(keys) {
			end = len(keys)
		}
		keysBatch := keys[i:end]
		textsBatch := texts[i:end]

		executionInfo, err := wasmhost.CallFunction(ctx, vectorIndex.GetEmbedderName(), textsBatch)
		if err != nil {
			return err
		}

		result := executionInfo.Result()

		textVecs, err := utils.ConvertToFloat32_2DArray(result)
		if err != nil {
			return err
		}

		if len(textVecs) == 0 {
			return fmt.Errorf("no vectors returned for texts: %v", textsBatch)
		}

		textIds := make([]int64, len(keysBatch))

		for i, key := range keysBatch {
			textId, err := col.GetExternalId(ctx, key)
			if err != nil {
				return err
			}
			textIds[i] = textId
		}

		err = vectorIndex.InsertVectors(ctx, textIds, textVecs)
		if err != nil {
			return err
		}
	}
	return nil
}

func processTextMap(ctx context.Context, col interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex) error {

	textMap, err := col.GetTextMap(ctx)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(textMap))
	texts := make([]string, 0, len(textMap))
	for key, text := range textMap {
		keys = append(keys, key)
		texts = append(texts, text)
	}

	return processTexts(ctx, col, vectorIndex, keys, texts)
}

func setIndex(ctx context.Context, collNs interfaces.CollectionNamespace, searchMethod manifest.SearchMethodInfo, searchMethodName string) error {
	vectorIndex, err := createIndexObject(searchMethod, searchMethodName)
	if err != nil {
		return err
	}

	err = collNs.SetVectorIndex(ctx, searchMethodName, vectorIndex)
	if err != nil {
		return err
	}
	return nil
}

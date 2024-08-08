package collections

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"hmruntime/collections/in_mem"
	"hmruntime/collections/in_mem/sequential"
	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"hmruntime/collections/utils"
	"hmruntime/db"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/wasmhost"

	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

const (
	batchSize = 25
)

func ProcessTexts(ctx context.Context, collection interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex, keys []string, texts []string) error {
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

		result := executionInfo.Result

		textVecs, err := utils.ConvertToFloat32_2DArray(result)
		if err != nil {
			return err
		}

		if len(textVecs) == 0 {
			return fmt.Errorf("no vectors returned for texts: %v", textsBatch)
		}

		textIds := make([]int64, len(keysBatch))

		for i, key := range keysBatch {
			textId, err := collection.GetExternalId(ctx, key)
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

func ProcessText(ctx context.Context, collection interfaces.CollectionNamespace, vectorIndex interfaces.VectorIndex, key, text string) error {
	texts := []string{text}
	executionInfo, err := wasmhost.CallFunction(ctx, vectorIndex.GetEmbedderName(), texts)
	if err != nil {
		return err
	}

	result := executionInfo.Result

	textVecs, err := utils.ConvertToFloat32_2DArray(result)
	if err != nil {
		return err
	}

	if len(textVecs) == 0 {
		return fmt.Errorf("no vectors returned for text: %s", text)
	}

	textId, err := collection.GetExternalId(ctx, key)
	if err != nil {
		return err
	}
	err = vectorIndex.InsertVector(ctx, textId, textVecs[0])
	if err != nil {
		return err
	}
	return nil
}

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, collection interfaces.CollectionNamespace, embedder string, vectorIndex interfaces.VectorIndex) error {

	textMap, err := collection.GetTextMap(ctx)
	if err != nil {
		return err
	}

	keys := make([]string, 0, len(textMap))
	texts := make([]string, 0, len(textMap))
	for key, text := range textMap {
		keys = append(keys, key)
		texts = append(texts, text)
	}

	return ProcessTexts(ctx, collection, vectorIndex, keys, texts)
}

func CleanAndProcessManifest(ctx context.Context) error {
	deleteIndexesNotInManifest(ctx, manifestdata.GetManifest())
	processManifestCollections(ctx, manifestdata.GetManifest())
	return nil
}

func processManifestCollections(ctx context.Context, man *manifest.HypermodeManifest) {
	for collectionName, collectionInfo := range man.Collections {
		collection, err := GlobalNamespaceManager.FindCollection(ctx, collectionName)
		if err == ErrCollectionNotFound {
			collection, err = GlobalNamespaceManager.CreateCollection(ctx, collectionName, NewCollection(ctx, collectionName))
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to create collection.")
			}
			// forces all users to use in-memory index for now
			// TODO implement other types of indexes based on manifest info
			// fetch all tenants and create a collection for each tenant
			namespaces, err := db.GetUniqueNamespaces(ctx)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to get unique namespaces.")
			}
			if !slices.Contains(namespaces, in_mem.DefaultNamespace) {
				namespaces = append(namespaces, in_mem.DefaultNamespace)
			}

			for _, namespace := range namespaces {
				_, err := collection.CreateCollectionNamespace(ctx, namespace, in_mem.NewCollectionNamespace(collectionName, namespace))
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", collectionName).
						Str("namespace", namespace).
						Msg("Failed to create collection namespace.")
				}
			}
		}
		for _, collNs := range collection.collectionNamespaceMap {
			for searchMethodName, searchMethod := range collectionInfo.SearchMethods {
				vi, err := collNs.GetVectorIndex(ctx, searchMethodName)

				// if the index does not exist, create it
				if err != nil {
					if err != index.ErrVectorIndexNotFound {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to get vector index.")
					} else {
						setIndex(ctx, collNs, searchMethod, searchMethodName)
					}
				} else if vi != nil && vi.Type != searchMethod.Index.Type {
					if err := collNs.DeleteVectorIndex(ctx, searchMethodName); err != nil {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to delete vector index.")
					} else {
						setIndex(ctx, collNs, searchMethod, searchMethodName)
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

func createIndexObject(searchMethod manifest.SearchMethodInfo, namespace, searchMethodName string) (*interfaces.VectorIndexWrapper, error) {
	vectorIndex := &interfaces.VectorIndexWrapper{}
	switch searchMethod.Index.Type {
	case interfaces.SequentialManifestType:
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, namespace, searchMethod.Embedder)
	case interfaces.HnswManifestType:
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, namespace, searchMethod.Embedder)
		// // TODO: hnsw currently broken ,autosync is not working, it keeps embedding forever, even though it has correctly indexed. for now, set it to sequential. fix in future PR
		// vectorIndex.Type = hnsw.HnswVectorIndexType
		// vectorIndex.VectorIndex = hnsw.NewHnswVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
	case "":
		vectorIndex.Type = sequential.SequentialVectorIndexType
		vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(searchMethodName, namespace, searchMethod.Embedder)
	default:
		return nil, fmt.Errorf("Unknown index type: %s", searchMethod.Index.Type)
	}

	return vectorIndex, nil
}

func setIndex(ctx context.Context, collNs interfaces.CollectionNamespace, searchMethod manifest.SearchMethodInfo, searchMethodName string) {
	vectorIndex, err := createIndexObject(searchMethod, collNs.GetNamespace(), searchMethodName)

	err = collNs.SetVectorIndex(ctx, searchMethodName, vectorIndex)
	if err != nil {
		logger.Err(ctx, err).
			Str("index_name", searchMethodName).
			Msg("Failed to create vector index.")
	}
}

func deleteIndexesNotInManifest(ctx context.Context, man *manifest.HypermodeManifest) {
	for collectionName := range GlobalNamespaceManager.GetNamespaceCollectionFactoryMap() {
		if _, ok := man.Collections[collectionName]; !ok {
			err := GlobalNamespaceManager.RemoveCollection(ctx, collectionName)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to remove collection.")
			}
			continue
		}
		collection, err := GlobalNamespaceManager.FindCollection(ctx, collectionName)
		if err != nil {
			logger.Err(ctx, err).
				Str("collection_name", collectionName).
				Msg("Failed to find collection.")
			continue
		}
		for _, collNs := range collection.collectionNamespaceMap {
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

func deleteVectorIndexesNotInManifest(ctx context.Context, man *manifest.HypermodeManifest, collection interfaces.CollectionNamespace, collectionName, searchMethodName string) error {
	_, ok := man.Collections[collectionName].SearchMethods[searchMethodName]
	if !ok {
		err := collection.DeleteVectorIndex(ctx, searchMethodName)
		if err != nil {
			return err
		}
	}
	return nil
}

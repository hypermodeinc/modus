package collections

import (
	"context"
	"fmt"

	"hmruntime/collections/in_mem"
	"hmruntime/collections/in_mem/hnsw"
	"hmruntime/collections/in_mem/sequential"
	"hmruntime/collections/index"
	"hmruntime/collections/index/interfaces"
	"hmruntime/collections/utils"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/wasmhost"

	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

const (
	batchSize = 25
)

func ProcessTexts(ctx context.Context, collection interfaces.Collection, vectorIndex interfaces.VectorIndex, keys []string, texts []string) error {
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

func ProcessText(ctx context.Context, collection interfaces.Collection, vectorIndex interfaces.VectorIndex, key, text string) error {
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

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	textMap, err := collection.GetTextMap(ctx)
	if err != nil {
		logger.Err(ctx, err).
			Str("collection_name", collection.GetCollectionName()).
			Msg("Failed to get text map.")
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
	deleteIndexesNotInManifest(ctx, manifestdata.Manifest)
	processManifestCollections(ctx, manifestdata.Manifest)
	return nil
}

func processManifestCollections(ctx context.Context, Manifest manifest.HypermodeManifest) {
	for collectionName, collectionInfo := range Manifest.Collections {
		collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
		if err == ErrCollectionNotFound {
			// forces all users to use in-memory index for now
			// TODO implement other types of indexes based on manifest info
			collection, err = GlobalCollectionFactory.Create(ctx, collectionName, in_mem.NewCollection(collectionName))
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to create vector index.")
			}
		}
		for searchMethodName, searchMethod := range collectionInfo.SearchMethods {
			vi, err := collection.GetVectorIndex(ctx, searchMethodName)

			// if the index does not exist, create it
			if err == index.ErrVectorIndexNotFound || (vi != nil && vi.Type != searchMethod.Index.Type) {
				if vi != nil && vi.Type != searchMethod.Index.Type {
					err = collection.DeleteVectorIndex(ctx, searchMethodName)
					if err != nil {
						logger.Err(ctx, err).
							Str("index_name", searchMethodName).
							Msg("Failed to delete vector index.")
					}
				}
				vectorIndex := &interfaces.VectorIndexWrapper{}
				switch searchMethod.Index.Type {
				case interfaces.SequentialManifestType:
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
				case interfaces.HnswManifestType:
					// TODO: Implement hnsw
					vectorIndex.Type = hnsw.HnswVectorIndexType
					vectorIndex.VectorIndex = hnsw.NewHnswVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
					// vectorIndex.Type = sequential.SequentialVectorIndexType
					// vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
				case "":
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
				default:
					logger.Err(ctx, nil).
						Str("index_type", searchMethod.Index.Type).
						Msg("Unknown index type.")
					continue
				}

				err = collection.SetVectorIndex(ctx, searchMethodName, vectorIndex)
				if err != nil {
					logger.Err(ctx, err).
						Str("index_name", searchMethodName).
						Msg("Failed to create vector index.")
				}
			} else if vi.GetEmbedderName() != searchMethod.Embedder {
				//TODO: figure out what to actually do if the embedder is different, for now just updating the name
				// what if the user changes the internals? -> they want us to reindex, model might have diff dimensions
				// but what if they just chaned the name? -> they want us to just update the name. but we cant know that
				// imo we should just update the name and let the user reindex if they want to
				err := vi.SetEmbedderName(searchMethod.Embedder)
				if err != nil {
					logger.Err(ctx, err).
						Str("index_name", searchMethodName).
						Msg("Failed to update vector index.")
				}
			}

		}
	}
}

func deleteIndexesNotInManifest(ctx context.Context, Manifest manifest.HypermodeManifest) {
	for collectionName := range GlobalCollectionFactory.GetCollectionMap() {
		if _, ok := Manifest.Collections[collectionName]; !ok {
			err := GlobalCollectionFactory.Remove(ctx, collectionName)
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to remove collection.")
			}
			continue
		}
		collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
		if err != nil {
			logger.Err(ctx, err).
				Str("collection_name", collectionName).
				Msg("Failed to find collection.")
			continue
		}
		vectorIndexMap := collection.GetVectorIndexMap()
		if vectorIndexMap == nil {
			continue
		}
		for searchMethodName := range vectorIndexMap {
			_, ok := Manifest.Collections[collectionName].SearchMethods[searchMethodName]
			if !ok {
				err = collection.DeleteVectorIndex(ctx, searchMethodName)
				if err != nil {
					logger.Err(ctx, err).
						Str("collectionName", collectionName).
						Str("search_method_name", searchMethodName).
						Msg("Failed to remove vector index.")
					continue
				}
			}
		}
	}
}

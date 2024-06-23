package collections

import (
	"context"

	"hmruntime/collections/in_mem"
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

func ProcessText(ctx context.Context, collection interfaces.Collection, vectorIndex interfaces.VectorIndex, key, text string) error {
	executionInfo, err := wasmhost.CallFunction(ctx, vectorIndex.GetEmbedderName(), text)
	if err != nil {
		return err
	}

	result := executionInfo.Result

	textVec, err := utils.ConvertToFloat32Array(result)
	if err != nil {
		return err
	}

	id, err := collection.GetExternalId(ctx, key)
	if err != nil {
		return err
	}
	err = vectorIndex.InsertVector(ctx, id, textVec)
	if err != nil {
		return err
	}
	return nil
}

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	textMap, err := collection.GetTextMap(ctx)
	if err != nil {
		logger.Err(ctx, err).
			Str("colletion_name", collection.GetCollectionName()).
			Msg("Failed to get text map.")
	}
	for key, text := range textMap {
		executionInfo, err := wasmhost.CallFunction(ctx, embedder, text)
		if err != nil {
			return err
		}

		result := executionInfo.Result

		textVec, err := utils.ConvertToFloat32Array(result)
		if err != nil {
			return err
		}

		id, err := collection.GetExternalId(ctx, key)
		if err != nil {
			return err
		}
		err = vectorIndex.InsertVector(ctx, id, textVec)
		if err != nil {
			return err
		}
	}
	return nil
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
			_, err := collection.GetVectorIndex(ctx, searchMethodName)

			// if the index does not exist, create it
			if err == index.ErrVectorIndexNotFound {
				vectorIndex := &interfaces.VectorIndexWrapper{}
				switch searchMethod.Index.Type {
				case interfaces.SequentialManifestType:
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
				case interfaces.HnswManifestType:
					// TODO: Implement hnsw
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex(collectionName, searchMethodName, searchMethod.Embedder)
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
		}
		collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
		if err != nil {
			logger.Err(ctx, err).
				Str("collection_name", collectionName).
				Msg("Failed to find collection.")
		}
		vectorIndexMap := collection.GetVectorIndexMap()
		if vectorIndexMap == nil {
			continue
		}
		for searchMethodName := range vectorIndexMap {
			_, ok := Manifest.Collections[collectionName].SearchMethods[searchMethodName]
			if !ok {
				collection, err := GlobalCollectionFactory.Find(ctx, collectionName)
				if err != nil {
					logger.Err(ctx, err).
						Str("collection_name", collectionName).
						Msg("Failed to find collection.")
				}
				err = collection.DeleteVectorIndex(ctx, searchMethodName)
				if err != nil {
					logger.Err(ctx, err).
						Str("collectionName", collectionName).
						Str("search_method_name", searchMethodName).
						Msg("Failed to remove vector index.")
				}
			}
		}
	}
}

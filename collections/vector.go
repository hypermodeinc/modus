package collections

import (
	"context"
	"strings"

	"hmruntime/collections/in_mem"
	"hmruntime/collections/in_mem/sequential"
	"hmruntime/collections/index/interfaces"
	"hmruntime/collections/utils"
	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/modules"

	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

type EmbedderFnCall struct {
	EmbedderFnName   string
	CollectionName   string
	SearchMethodName string
}

var FnCallChannel = make(chan EmbedderFnCall)

func ProcessTextMap(ctx context.Context, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range collection.GetTextMap() {
		result, err := modules.CallFunctionByName(ctx, embedder, text)
		if err != nil {
			return err
		}

		textVec, err := utils.ConvertToFloat32Array(result)
		if err != nil {
			return err
		}

		_, err = vectorIndex.InsertVector(ctx, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range collection.GetTextMap() {
		result, err := modules.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			return err
		}

		textVec, err := utils.ConvertToFloat32Array(result)
		if err != nil {
			return err
		}

		_, err = vectorIndex.InsertVector(ctx, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

func CleanAndProcessManifest(ctx context.Context) error {
	deleteIndexesNotInManifest(manifestdata.Manifest)
	processManifestCollections(ctx, manifestdata.Manifest)
	return nil
}

func processManifestCollections(ctx context.Context, Manifest manifest.HypermodeManifest) {
	for collectionName, collectionInfo := range Manifest.Collections {
		collection, err := GlobalCollectionFactory.Find(collectionName)
		if err == ErrCollectionNotFound {
			// forces all users to use in-memory index for now
			// TODO implement other types of indexes based on manifest info
			collection, err = GlobalCollectionFactory.Create(collectionName, in_mem.NewCollection())
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to create vector index.")
			}
		}
		for searchMethodName, searchMethod := range collectionInfo.SearchMethods {
			_, err := collection.GetVectorIndex(searchMethodName)

			// if the index does not exist, create it
			// TODO also populate the vector index by running the embedding function to compute vectors ahead of time
			if err == in_mem.ErrVectorIndexNotFound {
				vectorIndex := &interfaces.VectorIndexWrapper{}
				switch searchMethod.Index.Type {
				case interfaces.SequentialManifestType:
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex()
				case interfaces.HnswManifestType:
					// TODO: Implement hnsw
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex()
				case "":
					vectorIndex.Type = sequential.SequentialVectorIndexType
					vectorIndex.VectorIndex = sequential.NewSequentialVectorIndex()
				default:
					logger.Err(ctx, nil).
						Str("index_type", searchMethod.Index.Type).
						Msg("Unknown index type.")
					continue
				}

				err = collection.SetVectorIndex(searchMethodName, vectorIndex)
				if err != nil {
					logger.Err(ctx, err).
						Str("index_name", searchMethodName).
						Msg("Failed to create vector index.")
				}

				// populate index in background
				go func() {
					if len(collection.GetTextMap()) != 0 {
						err = ProcessTextMap(ctx, collection, searchMethod.Embedder, collection.GetVectorIndexMap()[searchMethodName])
						if err != nil {
							if strings.Contains(err.Error(), "no function registered named ") {
								FnCallChannel <- EmbedderFnCall{
									EmbedderFnName:   searchMethod.Embedder,
									CollectionName:   collectionName,
									SearchMethodName: searchMethodName,
								}
							} else {
								logger.Err(ctx, err).
									Str("index_name", searchMethodName).
									Msg("Failed to process text map.")
							}
						}
					}
				}()
			}
		}
	}
}

func deleteIndexesNotInManifest(Manifest manifest.HypermodeManifest) {
	for indexName := range GlobalCollectionFactory.GetCollectionMap() {
		if _, ok := Manifest.Collections[indexName]; !ok {
			err := GlobalCollectionFactory.Remove(indexName)
			if err != nil {
				logger.Err(context.Background(), err).
					Str("index_name", indexName).
					Msg("Failed to remove vector index.")
			}
		}
		vectorIndexMap := GlobalCollectionFactory.GetCollectionMap()[indexName].GetVectorIndexMap()
		if vectorIndexMap == nil {
			continue
		}
		for searchMethodName := range vectorIndexMap {
			_, ok := Manifest.Collections[indexName].SearchMethods[searchMethodName]
			if !ok {
				err := GlobalCollectionFactory.GetCollectionMap()[indexName].DeleteVectorIndex(searchMethodName)
				if err != nil {
					logger.Err(context.Background(), err).
						Str("index_name", indexName).
						Str("search_method_name", searchMethodName).
						Msg("Failed to remove vector index.")
				}
			}
		}
	}
}

func CatchEmbedderReqs(ctx context.Context) {
	go func() {
		for functionCall := range FnCallChannel {
			collection, err := GlobalCollectionFactory.Find(functionCall.CollectionName)
			if err != nil {
				logger.Err(context.Background(), err).Msg("Error finding collection")
				continue
			}

			err = ProcessTextMap(ctx, collection, functionCall.EmbedderFnName,
				collection.GetVectorIndexMap()[functionCall.SearchMethodName])

			if err != nil {
				logger.Err(context.Background(), err).Msg("Error processing text map")
			}
		}
	}()
}

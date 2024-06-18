/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"

	"hmruntime/collections"
	"hmruntime/collections/in_mem"
	"hmruntime/collections/in_mem/sequential"
	"hmruntime/collections/index/interfaces"
	"hmruntime/logger"
	"hmruntime/storage"
	"hmruntime/utils"

	"github.com/hypermodeAI/manifest"
)

const manifestFileName = "hypermode.json"

var Manifest manifest.HypermodeManifest = manifest.HypermodeManifest{}

func MonitorManifestFile(ctx context.Context) {
	loadFile := func(file storage.FileInfo) error {
		if file.Name != manifestFileName {
			return nil
		}
		err := loadManifest(ctx)
		if err == nil {
			logger.Info(ctx).
				Str("filename", file.Name).
				Msg("Loaded manifest file.")
		} else {
			logger.Err(ctx, err).
				Str("filename", file.Name).
				Msg("Failed to load manifest file.")
		}

		return err
	}

	// NOTE: Removing the manifest file entirely is not currently supported.

	sm := storage.NewStorageMonitor(".json")
	sm.Added = loadFile
	sm.Modified = loadFile
	sm.Start(ctx)
}

func loadManifest(ctx context.Context) error {
	transaction, ctx := utils.NewSentryTransactionForCurrentFunc(ctx)
	defer transaction.Finish()

	bytes, err := storage.GetFileContents(ctx, manifestFileName)
	if err != nil {
		return err
	}

	man, err := manifest.ReadManifest(bytes)
	if err != nil {
		return err
	}

	if !man.IsCurrentVersion() {
		logger.Warn(ctx).
			Str("filename", manifestFileName).
			Int("manifest_version", man.Version).
			Msg("The manifest file is in a deprecated format.  Please update it to the current format.")
	}

	// remove indexes that are not in the manifest
	deleteIndexesNotInManifest(man)
	// add processing of manifest collections to create vector indexes
	processManifestCollections(ctx, man)

	// Only update the Manifest global when we have successfully read the manifest.
	Manifest = man

	return nil
}

func processManifestCollections(ctx context.Context, Manifest manifest.HypermodeManifest) {
	for collectionName, collectionInfo := range Manifest.Collections {
		collection, err := collections.GlobalCollectionFactory.Find(collectionName)
		if err == collections.ErrCollectionNotFound {
			// forces all users to use in-memory index for now
			// TODO implement other types of indexes based on manifest info
			collection, err = collections.GlobalCollectionFactory.Create(collectionName, in_mem.NewCollection())
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
						err = collections.ProcessTextMap(ctx, collection, searchMethod.Embedder, collection.GetVectorIndexMap()[searchMethodName])
						if err != nil {
							logger.Err(ctx, err).
								Str("index_name", searchMethodName).
								Msg("Failed to process text map.")
						}
					}
				}()
			}
		}
	}
}

func deleteIndexesNotInManifest(Manifest manifest.HypermodeManifest) {
	for indexName := range collections.GlobalCollectionFactory.GetCollectionMap() {
		if _, ok := Manifest.Collections[indexName]; !ok {
			err := collections.GlobalCollectionFactory.Remove(indexName)
			if err != nil {
				logger.Err(context.Background(), err).
					Str("index_name", indexName).
					Msg("Failed to remove vector index.")
			}
		}
		vectorIndexMap := collections.GlobalCollectionFactory.GetCollectionMap()[indexName].GetVectorIndexMap()
		if vectorIndexMap == nil {
			continue
		}
		for searchMethodName := range vectorIndexMap {
			_, ok := Manifest.Collections[indexName].SearchMethods[searchMethodName]
			if !ok {
				err := collections.GlobalCollectionFactory.GetCollectionMap()[indexName].DeleteVectorIndex(searchMethodName)
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

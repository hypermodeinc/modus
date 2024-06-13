/*
 * Copyright 2024 Hypermode, Inc.
 */

package manifestdata

import (
	"context"

	"hmruntime/logger"
	"hmruntime/storage"
	"hmruntime/utils"
	"hmruntime/vector"
	"hmruntime/vector/in_mem"
	"hmruntime/vector/index"

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

	// add processing of manifest collections to create vector indexes
	processManifestCollections(ctx, man)
	deleteIndexesNotInManifest(man)

	// Only update the Manifest global when we have successfully read the manifest.
	Manifest = man

	return nil
}

func processManifestCollections(ctx context.Context, Manifest manifest.HypermodeManifest) {
	for collectionName, collection := range Manifest.Collections {
		textIndex, err := vector.GlobalTextIndexFactory.Find(collectionName)
		if err == vector.ErrTextIndexNotFound {
			textIndex, err = vector.GlobalTextIndexFactory.Create(collectionName, &in_mem.InMemTextIndex[float64]{})
			if err != nil {
				logger.Err(ctx, err).
					Str("collection_name", collectionName).
					Msg("Failed to create vector index.")
			}
		}
		for searchMethodName, searchMethod := range collection.SearchMethods {
			_, err := textIndex.GetVectorIndex(searchMethodName)

			// if the index does not exist, create it
			if err == in_mem.ErrVectorIndexAlreadyExists {
				var index index.VectorIndex[float64]
				switch searchMethod.Index.Type {
				case "sequential":
					index = &in_mem.InMemBruteForceIndex{}
				case "":
					index = &in_mem.InMemBruteForceIndex{}
				default:
					logger.Err(ctx, nil).
						Str("index_type", searchMethod.Index.Type).
						Msg("Unknown index type.")
					continue
				}

				textIndex.SetVectorIndex(searchMethodName, index)

				if err != nil {
					logger.Err(ctx, err).
						Str("index_name", searchMethodName).
						Msg("Failed to create vector index.")
				}
			}
		}
	}
}

func deleteIndexesNotInManifest(Manifest manifest.HypermodeManifest) {
	for indexName := range vector.GlobalTextIndexFactory.GetTextIndexMap() {
		if _, ok := Manifest.Collections[indexName]; !ok {
			vector.GlobalTextIndexFactory.Remove(indexName)
		}
		for searchMethodName := range vector.GlobalTextIndexFactory.GetTextIndexMap()[indexName].GetVectorIndexMap() {
			if _, ok := Manifest.Collections[indexName].SearchMethods[searchMethodName]; !ok {
				vector.GlobalTextIndexFactory.GetTextIndexMap()[indexName].DeleteVectorIndex(searchMethodName)
			}
		}
	}
}

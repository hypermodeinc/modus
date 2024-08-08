/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/manifestdata"
	"hmruntime/sqlclient"
	"hmruntime/utils"
	"hmruntime/wasmhost"

	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

const hostModuleName string = "hypermode"

func RegisterHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	instantiateHostFunctions(ctx)
	instantiateWasiFunctions(ctx)

	manifestdata.RegisterManifestLoadedCallback(func(ctx context.Context) error {
		sqlclient.ShutdownPGPools()
		return nil
	})
}

func instantiateHostFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := wasmhost.RuntimeInstance.NewHostModuleBuilder(hostModuleName)

	// Misc host functions
	b.NewFunctionBuilder().WithFunc(hostLog).Export("log")
	b.NewFunctionBuilder().WithFunc(hostFetch).Export("httpFetch")
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostDatabaseQuery).Export("databaseQuery")

	// Model host functions
	b.NewFunctionBuilder().WithFunc(hostLookupModel).Export("lookupModel")
	b.NewFunctionBuilder().WithFunc(hostInvokeModel).Export("invokeModel")

	// Legacy model host functions
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")
	b.NewFunctionBuilder().WithFunc(hostInvokeTextGenerator).Export("invokeTextGenerator")

	// Collection host functions
	b.NewFunctionBuilder().WithFunc(hostUpsertToCollection).Export("upsertToCollection")
	b.NewFunctionBuilder().WithFunc(hostUpsertToCollectionV2).Export("upsertToCollection_v2")
	b.NewFunctionBuilder().WithFunc(hostDeleteFromCollection).Export("deleteFromCollection")
	b.NewFunctionBuilder().WithFunc(hostDeleteFromCollectionV2).Export("deleteFromCollection_v2")
	b.NewFunctionBuilder().WithFunc(hostSearchCollection).Export("searchCollection")
	b.NewFunctionBuilder().WithFunc(hostSearchCollectionV2).Export("searchCollection_v2")
	b.NewFunctionBuilder().WithFunc(hostNnClassifyCollection).Export("nnClassifyCollection")
	b.NewFunctionBuilder().WithFunc(hostNnClassifyCollectionV2).Export("nnClassifyCollection_v2")
	b.NewFunctionBuilder().WithFunc(hostRecomputeSearchMethod).Export("recomputeSearchMethod")
	b.NewFunctionBuilder().WithFunc(hostRecomputeSearchMethodV2).Export("recomputeSearchMethod_v2")
	b.NewFunctionBuilder().WithFunc(hostComputeDistance).Export("computeSimilarity") // Deprecated
	b.NewFunctionBuilder().WithFunc(hostComputeDistance).Export("computeDistance")
	b.NewFunctionBuilder().WithFunc(hostComputeDistanceV2).Export("computeDistance_v2")
	b.NewFunctionBuilder().WithFunc(hostGetTextFromCollection).Export("getTextFromCollection")
	b.NewFunctionBuilder().WithFunc(hostGetTextFromCollectionV2).Export("getTextFromCollection_v2")
	b.NewFunctionBuilder().WithFunc(hostGetTextsFromCollection).Export("getTextsFromCollection")
	b.NewFunctionBuilder().WithFunc(hostGetTextsFromCollectionV2).Export("getTextsFromCollection_v2")

	if _, err := b.Instantiate(ctx); err != nil {
		logger.Fatal(ctx).Err(err).
			Str("module", hostModuleName).
			Msg("Failed to instantiate the host module.  Exiting.")
	}
}

func instantiateWasiFunctions(ctx context.Context) {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := wasmhost.RuntimeInstance.NewHostModuleBuilder(wasi.ModuleName)
	wasi.NewFunctionExporter().ExportFunctions(b)

	// If we ever need to override any of the WASI functions, we can do so here.

	if _, err := b.Instantiate(ctx); err != nil {
		logger.Fatal(ctx).Err(err).
			Str("module", wasi.ModuleName).
			Msg("Failed to instantiate the host module.  Exiting.")
	}
}

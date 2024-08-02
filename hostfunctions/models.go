/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"

	"hmruntime/logger"
	"hmruntime/models"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostLookupModel(ctx context.Context, mod wasm.Module, pModelName uint32) (pModelInfo uint32) {
	var modelName string
	err := readParam(ctx, mod, pModelName, &modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	info, err := models.GetModelInfo(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model info.")
		return 0
	}

	offset, err := writeResult(ctx, mod, info)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}

func hostInvokeModel(ctx context.Context, mod wasm.Module, pModelName uint32, pInput uint32) (pOutput uint32) {
	var modelName, input string
	err := readParams2(ctx, mod, pModelName, pInput, &modelName, &input)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	output, err := models.InvokeModel(ctx, modelName, input)
	if err != nil {
		logger.Err(ctx, err).Msg("Error invoking model.")
		return 0
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

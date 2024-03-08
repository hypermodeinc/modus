/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/config"
	"hmruntime/dgraph"
	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/models/openai"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const HostModuleName string = "hypermode"

func InstantiateHostFunctions(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(HostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteDQL).Export("executeDQL")
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")
	b.NewFunctionBuilder().WithFunc(hostInvokeTextGenerator).Export("invokeTextGenerator")

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", HostModuleName, err)
	}

	return nil
}

func hostExecuteDQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32, isMutation uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading DQL statement from wasm memory.")
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading DQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteDQL[string](ctx, stmt, vars, isMutation != 0)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing DQL statement.")
		return 0
	}

	return writeString(ctx, mod, result)
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL query string from wasm memory.")
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteGQL[string](ctx, stmt, vars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing GraphQL operation.")
		return 0
	}

	return writeString(ctx, mod, result)
}

type ClassifierResult struct {
	Probabilities []ClassifierLabel `json:"probabilities"`
}

type ClassifierLabel struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	modelSpec, err := getModelSpec(mem, pModelName, config.ClassificationModelType)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model spec.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[map[string]ClassifierResult](ctx, sentenceMap, modelSpec)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from model.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling classification result.")
		return 0
	}

	return writeString(ctx, mod, string(res))
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	modelSpec, err := getModelSpec(mem, pModelName, config.EmbeddingModelType)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model spec.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[map[string]string](ctx, sentenceMap, modelSpec)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Error(ctx).Msg("Empty result returned from model.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling embedding result.")
		return 0
	}

	return writeString(ctx, mod, string(res))
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32) uint32 {
	// invoke modelspec endpoint or
	// https://api.openai.com/v1/chat/completions if endpoint is "openai"

	mem := mod.Memory()

	modelSpec, err := getModelSpec(mem, pModelName, config.GeneratorModelType)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model spec.")
		return 0
	}

	// we are assuming gpt-3.5-turbo
	// to do : define how to pass the model name in the modelSpec

	sentence, err := readString(mem, pSentence)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading sentence string from wasm memory.")
		return 0
	}

	instruction, err := readString(mem, pInstruction)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading instruction string from wasm memory.")
		return 0
	}

	// Call appropriate implementation depending on the modelSpec.Provider
	// TO DO: use Provider when it will be available in the modelSpec

	result, err := openai.GenerateText(ctx, modelSpec, instruction, sentence)

	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to OpenAI.")
		return 0
	}
	if result.Error.Message != "" {
		err := fmt.Errorf(result.Error.Message)
		logger.Err(ctx, err).Msg("Error returned from OpenAI.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling result.")
		return 0
	}
	return writeString(ctx, mod, string(res))
}

func getModelSpec(mem wasm.Memory, pModelName uint32, modelType config.ModelType) (config.ModelSpec, error) {
	modelName, err := readString(mem, pModelName)
	if err != nil {
		err = fmt.Errorf("error reading model name from wasm memory: %w", err)
		return config.ModelSpec{}, err
	}

	return models.GetModelSpec(modelName, modelType)
}

func getSentenceMap(mem wasm.Memory, pSentenceMap uint32) (map[string]string, error) {
	sentenceMapStr, err := readString(mem, pSentenceMap)
	if err != nil {
		err = fmt.Errorf("error reading sentence map string from wasm memory: %w", err)
		return nil, err
	}

	sentenceMap := make(map[string]string)
	if err := json.Unmarshal([]byte(sentenceMapStr), &sentenceMap); err != nil {
		err = fmt.Errorf("error unmarshalling sentence map: %w", err)
		return nil, err
	}

	return sentenceMap, nil
}

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
	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const HostModuleName string = "hypermode"
const HypermodeHost string = "hypermode" // separate strings in case one needs to change

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
	Label         string            `json:"label"`
	Confidence    float64           `json:"confidence"`
	Probabilities []ClassifierLabel `json:"probabilities"`
}

type ClassifierLabel struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, config.ClassificationTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := PostToClassifierModelEndpoint(ctx, sentenceMap, model)
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

type KServeClassifierResult struct {
	Predictions []ClassifierResult `json:"predictions"`
}

func PostToClassifierModelEndpoint(ctx context.Context, sentenceMap map[string]string, model config.Model) (map[string]ClassifierResult, error) {
	// self hosted models takes in array, can optimize for parallelizing later
	if model.Host == HypermodeHost {
		endpoint := fmt.Sprintf("http://%s.%s/%s:predict", model.Task, config.ModelUrl, model.Name)
		sentences := make(map[string][]string)
		keys := []string{}
		sentences["instances"] = make([]string, 0, len(sentenceMap))

		for k, v := range sentenceMap {
			// create a map of sentences to send to the model
			sentences["instances"] = append(sentences["instances"], v)
			// create a list of keys to map the results back to the original sentences
			keys = append(keys, k)
		}
		ksRes, err := utils.PostHttp[KServeClassifierResult](endpoint, sentences, nil)
		if err != nil {
			return map[string]ClassifierResult{}, err
		}
		if len(ksRes.Predictions) != len(keys) {
			return map[string]ClassifierResult{}, fmt.Errorf("number of predictions does not match number of sentences")
		}
		// map the results back to the original sentences
		result := make(map[string]ClassifierResult)
		for i, v := range ksRes.Predictions {
			result[keys[i]] = v
		}
		return result, nil
	}
	key, err := models.GetModelKey(ctx, model)
	if err != nil {
		var result map[string]ClassifierResult
		return result, fmt.Errorf("error getting model key secret: %w", err)
	}

	headers := map[string]string{
		model.AuthHeader: key,
	}

	return utils.PostHttp[map[string]ClassifierResult](model.Endpoint, sentenceMap, headers)
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, config.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := PostToEmbeddingModelEndpoint(ctx, sentenceMap, model)
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

type KServeEmbeddingResult struct {
	Predictions [][]float64 `json:"predictions"`
}

func PostToEmbeddingModelEndpoint(ctx context.Context, sentenceMap map[string]string, model config.Model) (map[string]string, error) {
	// self hosted models takes in array, can optimize for parallelizing later
	if model.Host == HypermodeHost {
		endpoint := fmt.Sprintf("http://%s.%s/%s:predict", model.Task, config.ModelUrl, model.Name)
		sentences := make(map[string][]string)
		keys := []string{}
		sentences["instances"] = make([]string, 0, len(sentenceMap))

		for k, v := range sentenceMap {
			// create a map of sentences to send to the model
			sentences["instances"] = append(sentences["instances"], v)
			// create a list of keys to map the results back to the original sentences
			keys = append(keys, k)
		}
		ksRes, err := utils.PostHttp[KServeEmbeddingResult](endpoint, sentences, nil)
		if err != nil {
			return map[string]string{}, err
		}
		if len(ksRes.Predictions) != len(keys) {
			return map[string]string{}, fmt.Errorf("number of predictions does not match number of sentences")
		}
		// map the results back to the original sentences
		result := make(map[string]string)
		for i, v := range ksRes.Predictions {
			result[keys[i]] = fmt.Sprintf("%v", v)
		}
		return result, nil
	}
	return models.PostToModelEndpoint[map[string]string](ctx, sentenceMap, model)
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, config.GeneratorTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

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

	result, err := openai.GenerateText(ctx, model, instruction, sentence)

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

func invokeModel(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32, task config.ModelTask) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, task)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}
	result, err := models.PostToModelEndpoint[string](ctx, sentenceMap, model)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	return writeString(ctx, mod, result)
}

func getModel(mem wasm.Memory, pModelName uint32, task config.ModelTask) (config.Model, error) {
	modelName, err := readString(mem, pModelName)
	if err != nil {
		err = fmt.Errorf("error reading model name from wasm memory: %w", err)
		return config.Model{}, err
	}

	return models.GetModel(modelName, task)
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

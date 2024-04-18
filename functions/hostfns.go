/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"hmruntime/appdata"
	"hmruntime/dgraph"
	"hmruntime/functions/assemblyscript"
	"hmruntime/logger"
	"hmruntime/models"

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
	b.NewFunctionBuilder().WithFunc(hostInvokeTextGeneratorV2).Export("invokeTextGenerator_v2")

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", HostModuleName, err)
	}

	return nil
}

func hostExecuteDQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32, isMutation uint32) uint32 {
	mem := mod.Memory()
	stmt, err := assemblyscript.ReadString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading DQL statement from wasm memory.")
		return 0
	}

	sVars, err := assemblyscript.ReadString(mem, pVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading DQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]any)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteDQL[string](ctx, stmt, vars, isMutation != 0)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing DQL statement.")
		return 0
	}

	return assemblyscript.WriteString(ctx, mod, result)
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()
	stmt, err := assemblyscript.ReadString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL query string from wasm memory.")
		return 0
	}

	sVars, err := assemblyscript.ReadString(mem, pVars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading GraphQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]any)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		logger.Err(ctx, err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteGQL[string](ctx, stmt, vars)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing GraphQL operation.")
		return 0
	}

	return assemblyscript.WriteString(ctx, mod, result)
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

	model, err := getModel(mem, pModelName, appdata.ClassificationTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[ClassifierResult](ctx, sentenceMap, model)
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

	return assemblyscript.WriteString(ctx, mod, string(res))
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, appdata.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}
	var result map[string][]float64
	if (models.OpenAIHost == model.Host) || (models.MistralHost == model.Host) {
		// use LLM
		llm, err := models.CreateLlmService(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error instanciating LLM")
			return 0
		}

		result, err = llm.Embedding(ctx, sentenceMap, model)
		if err != nil {
			logger.Err(ctx, err).Msg("Error embeddings with LLM.")
			return 0
		}

	} else {

		result, err = models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model)
		if err != nil {
			logger.Err(ctx, err).Msg("Error posting to model endpoint.")
			return 0
		}

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

	return assemblyscript.WriteString(ctx, mod, string(res))
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32) uint32 {
	return hostInvokeTextGeneratorV2(ctx, mod, pModelName, pInstruction, pSentence, 0)
}

func hostInvokeTextGeneratorV2(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32, pFormat uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, appdata.GeneratorTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	sentence, err := assemblyscript.ReadString(mem, pSentence)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading sentence string from wasm memory.")
		return 0
	}

	instruction, err := assemblyscript.ReadString(mem, pInstruction)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading instruction string from wasm memory.")
		return 0
	}

	// Default to text output format for backwards compatibility.
	// V2 allows for a format to be specified.
	var outputFormat models.OutputFormat
	if pFormat == 0 {
		outputFormat = models.OutputFormatText
	} else {
		format, err := assemblyscript.ReadString(mem, pFormat)
		if err != nil {
			logger.Err(ctx, err).Msg("Error reading format string from wasm memory.")
			return 0
		}
		outputFormat = models.OutputFormat(format)
	}

	if models.OutputFormatText != outputFormat && models.OutputFormatJson != outputFormat {
		logger.Err(ctx, err).Msg("Unsupported output format.")
		return 0
	}
	llm, err := models.CreateLlmService(model.Host)
	if err != nil {
		logger.Err(ctx, err).Msg("Error instanciating LLM")
		return 0
	}

	result, err := llm.ChatCompletion(ctx, model, instruction, sentence, outputFormat)
	if err != nil {
		logger.Err(ctx, err)
		return 0
	}
	if result.Error.Message != "" {
		err := fmt.Errorf(result.Error.Message)
		logger.Err(ctx, err)
		return 0
	}
	if models.OutputFormatJson == outputFormat {
		// safeguard: test is the output is a valid json
		// test every Choices.Message.Content
		for _, choice := range result.Choices {
			_, err := json.Marshal(choice.Message.Content)
			if err != nil {
				logger.Err(ctx, err).Msg("One of the generated message is not a valid JSON.")
				return 0
			}
		}
	}

	res, err := json.Marshal(result)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling result.")
		return 0
	}
	return assemblyscript.WriteString(ctx, mod, string(res))
}

func getLLM(s string) {
	panic("unimplemented")
}

func getModel(mem wasm.Memory, pModelName uint32, task appdata.ModelTask) (appdata.Model, error) {
	modelName, err := assemblyscript.ReadString(mem, pModelName)
	if err != nil {
		err = fmt.Errorf("error reading model name from wasm memory: %w", err)
		return appdata.Model{}, err
	}

	return models.GetModel(modelName, task)
}

func getSentenceMap(mem wasm.Memory, pSentenceMap uint32) (map[string]string, error) {
	sentenceMapStr, err := assemblyscript.ReadString(mem, pSentenceMap)
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

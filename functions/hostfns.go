/*
 * Copyright 2023 Hypermode, Inc.
 */

package functions

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/appdata"
	"hmruntime/connections"
	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/models/openai"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const HostModuleName string = "hypermode"
const HypermodeHostName string = "hypermode"

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

func hostExecuteDQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVars uint32, isMutation uint32) uint32 {
	mem := mod.Memory()
	stmt, err := assemblyscript.ReadString(mem, pStmt)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading DQL statement from wasm memory.")
		return 0
	}

	host, err := getHost(mem, pHostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
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

	result, err := connections.ExecuteDQL[string](ctx, host, stmt, vars, isMutation != 0)
	if err != nil {
		logger.Err(ctx, err).Msg("Error executing DQL statement.")
		return 0
	}

	return assemblyscript.WriteString(ctx, mod, result)
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pHostName uint32, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()

	host, err := getHost(mem, pHostName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting host.")
		return 0
	}

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

	result, err := connections.ExecuteGraphqlApi[string](ctx, host, stmt, vars)
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

	var host appdata.Host
	if model.Host != HypermodeHostName {
		host, err = hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[ClassifierResult](ctx, sentenceMap, model, host)
	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from model.")
		return 0
	}

	resultMap := make(map[string]map[string]float64)
	for k, v := range result {
		resultMap[k] = make(map[string]float64)
		for _, label := range v.Probabilities {
			resultMap[k][label.Label] = label.Probability
		}
	}

	resBytes, err := json.Marshal(resultMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling classification result.")
		return 0
	}

	return assemblyscript.WriteString(ctx, mod, string(resBytes))
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelName uint32, pSentenceMap uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, appdata.EmbeddingTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	var host appdata.Host
	if model.Host != HypermodeHostName {
		host, err = hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}
	}

	sentenceMap, err := getSentenceMap(mem, pSentenceMap)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting sentence map.")
		return 0
	}

	result, err := models.PostToModelEndpoint[[]float64](ctx, sentenceMap, model, host)
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

	return assemblyscript.WriteString(ctx, mod, string(res))
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32, pFormat uint32) uint32 {
	mem := mod.Memory()

	model, err := getModel(mem, pModelName, appdata.GenerationTask)
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

	var host appdata.Host
	if model.Host != HypermodeHostName {
		host, err = hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}
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

	var result models.ChatResponse
	switch model.Host {
	case hosts.OpenAIHost:
		result, err := openai.ChatCompletion(ctx, model, host, instruction, sentence, outputFormat)
		if err != nil {
			logger.Err(ctx, err).Msg("Error posting to OpenAI.")
			return 0
		}
		if result.Error.Message != "" {
			err := fmt.Errorf(result.Error.Message)
			logger.Err(ctx, err).Msg("Error returned from OpenAI.")
			return 0
		}
	default:
		err := fmt.Errorf("unsupported model host: %s", model.Host)
		logger.Err(ctx, err).Msg("Unsupported model host.")
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

	// return the first chat response
	if len(result.Choices) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from OpenAI.")
		return 0
	}
	firstMsgContent := result.Choices[0].Message.Content

	resBytes, err := json.Marshal(firstMsgContent)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling result.")
		return 0
	}
	return assemblyscript.WriteString(ctx, mod, string(resBytes))
}

func getHost(mem wasm.Memory, pHostName uint32) (appdata.Host, error) {
	hostName, err := assemblyscript.ReadString(mem, pHostName)
	if err != nil {
		err = fmt.Errorf("error reading host name from wasm memory: %w", err)
		return appdata.Host{}, err
	}

	return hosts.GetHost(hostName)
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

/*
 * Copyright 2023 Hypermode, Inc.
 */
package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"hmruntime/aws"
	"hmruntime/dgraph"
	"hmruntime/utils"

	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const (
	HostModuleName  string = "hypermode"
	classifierModel string = "classifier"
	embeddingModel  string = "embedding"
	generatorModel  string = "generator"
)

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
		log.Err(err).Msg("Error reading DQL statement from wasm memory.")
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		log.Err(err).Msg("Error reading DQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		log.Err(err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteDQL[string](ctx, stmt, vars, isMutation != 0)
	if err != nil {
		log.Err(err).Msg("Error executing DQL statement.")
		return 0
	}

	return writeString(ctx, mod, result)
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		log.Err(err).Msg("Error reading GraphQL query string from wasm memory.")
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		log.Err(err).Msg("Error reading GraphQL variables string from wasm memory.")
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		log.Err(err).Msg("Error unmarshalling GraphQL variables.")
		return 0
	}

	result, err := dgraph.ExecuteGQL[string](ctx, stmt, vars)
	if err != nil {
		log.Err(err).Msg("Error executing GraphQL operation.")
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

func textModelSetup(mod wasm.Module, pModelId uint32, pSentenceMap uint32) (dgraph.ModelSpec, map[string]string, error) {
	mem := mod.Memory()
	modelSpec, err := textModelSpecSetup(mod, pModelId)
	if err != nil {
		return dgraph.ModelSpec{}, nil, err
	}

	sentenceMapStr, err := readString(mem, pSentenceMap)
	if err != nil {
		err = fmt.Errorf("error reading sentence map string from wasm memory: %w", err)
		return dgraph.ModelSpec{}, nil, err
	}

	sentenceMap := make(map[string]string)
	if err := json.Unmarshal([]byte(sentenceMapStr), &sentenceMap); err != nil {
		err = fmt.Errorf("error unmarshalling sentence map: %w", err)
		return dgraph.ModelSpec{}, nil, err
	}

	return modelSpec, sentenceMap, nil
}
func textModelSpecSetup(mod wasm.Module, pModelId uint32) (dgraph.ModelSpec, error) {
	mem := mod.Memory()
	modelId, err := readString(mem, pModelId)
	if err != nil {
		err = fmt.Errorf("error reading model id from wasm memory: %w", err)
		return dgraph.ModelSpec{}, err
	}

	modelSpec, err := dgraph.GetModelSpec(modelId)
	if err != nil {
		err = fmt.Errorf("error getting model endpoint: %w", err)
		return dgraph.ModelSpec{}, err
	}

	return modelSpec, nil
}

func postToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, modelSpec dgraph.ModelSpec) (TResult, error) {

	key, err := aws.GetSecretString(ctx, modelSpec.ID)
	if err != nil {
		var result TResult
		return result, fmt.Errorf("error getting model key '%s': %w", modelSpec.ID, err)
	}

	headers := map[string]string{
		"x-api-key": key,
	}

	return utils.PostHttp[TResult](modelSpec.Endpoint, sentenceMap, headers)
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, pModelId uint32, pSentenceMap uint32) uint32 {
	modelSpec, sentenceMap, err := textModelSetup(mod, pModelId, pSentenceMap)
	if err != nil {
		log.Err(err).Msg("Error setting up text model.")
		return 0
	}

	if modelSpec.Type != classifierModel {
		log.Error().Msg("Model type is not 'classifier'.")
		return 0
	}

	result, err := postToModelEndpoint[map[string]ClassifierResult](ctx, sentenceMap, modelSpec)
	if err != nil {
		log.Err(err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		log.Err(err).Msg("Empty result returned from model.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Err(err).Msg("Error marshalling classifier result.")
		return 0
	}

	return writeString(ctx, mod, string(res))
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, pModelId uint32, pSentenceMap uint32) uint32 {
	modelSpec, sentenceMap, err := textModelSetup(mod, pModelId, pSentenceMap)
	if err != nil {
		log.Err(err).Msg("Error setting up text model.")
		return 0
	}

	if modelSpec.Type != embeddingModel {
		log.Error().Msg("Model type is not 'embedding'.")
		return 0
	}

	result, err := postToModelEndpoint[map[string]string](ctx, sentenceMap, modelSpec)
	if err != nil {
		log.Err(err).Msg("Error posting to model endpoint.")
		return 0
	}

	if len(result) == 0 {
		log.Error().Msg("Empty result returned from model.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Err(err).Msg("Error marshalling embedding result.")
		return 0
	}

	return writeString(ctx, mod, string(res))
}

type ChatContext struct {
	Model    string        `json:"model"`
	Messages []ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type MessageChoice struct {
	Message ChatMessage `json:"message"`
}

type InvokeError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param"`
	Code    string `json:"code"`
}
type ChatResponse struct {
	Choices []MessageChoice `json:"choices"`
	Error   InvokeError     `json:"error"`
}

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelId uint32, pInstruction uint32, pSentence uint32) uint32 {
	// invoke modelspec endpoint or
	// https://api.openai.com/v1/chat/completions if endpoint is "openai"
	mem := mod.Memory()
	modelSpec, err := textModelSpecSetup(mod, pModelId)

	if err != nil {
		log.Err(err)
		return 0
	}
	if modelSpec.Type != generatorModel {
		log.Error().Msg(fmt.Sprintf("Model type is not '%s'.", generatorModel))
		return 0
	}

	// we are assuming gpt-3.5-turbo
	// to do : define how to pass the model name in the modelSpec

	sentence, err := readString(mem, pSentence)
	if err != nil {
		log.Print("error reading sentence string from wasm memory:", err)
		return 0
	}
	instruction, err := readString(mem, pInstruction)
	if err != nil {
		log.Print("error reading instruction string from wasm memory:", err)
		return 0
	}

	jinstruction, _ := json.Marshal(instruction)
	jsentence, _ := json.Marshal(sentence)

	key, err := aws.GetSecretString(ctx, modelSpec.ID)
	if err != nil {
		log.Print("error getting model key:", err)
		return 0
	}
	// Call appropriate impelmentation depending on the modelSpec.Provider
	// TO DO: use Provider when it will be available in the modelSpec

	result, err := openaiTextGenerator(ctx, modelSpec, string(jinstruction), string(jsentence), key)

	if err != nil {
		log.Err(err).Msg("Error posting to openai.")
		return 0
	}
	if result.Error.Message != "" {
		err := fmt.Errorf(result.Error.Message)
		log.Err(err).Msg("Error returned from openai.")
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Err(err).Msg("Error marshalling result.")
		return 0
	}
	return writeString(ctx, mod, string(res))
}

func openaiTextGenerator(ctx context.Context, modelSpec dgraph.ModelSpec, instruction string, sentence string, key string) (ChatResponse, error) {
	// invoke modelspec endpoint or
	// https://api.openai.com/v1/chat/completions if endpoint is "openai"

	// we are assuming gpt-3.5-turbo
	// should be modelSpec.baseModel when it will be available in the modelSpec
	const model = "gpt-3.5-turbo"

	// build the request body following openai API

	reqBody := ChatContext{
		Model: model,
		Messages: []ChatMessage{
			{Role: "system", Content: instruction},
			{Role: "user", Content: sentence},
		},
	}

	// We ignore modelSpec.endpoint and use the openai endpoint
	const endpoint = "https://api.openai.com/v1/chat/completions"
	headers := map[string]string{
		"Authorization": "Bearer " + key,
		"Content-Type":  "application/json",
	}

	return utils.PostHttp[ChatResponse](endpoint, reqBody, headers)

}

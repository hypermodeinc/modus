/*
 * Copyright 2023 Hypermode, Inc.
 */
package functions

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hmruntime/aws"
	"hmruntime/dgraph"
	"hmruntime/utils"
	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const (
	HostModuleName  string = "hypermode"
	classifierModel string = "classifier"
	embeddingModel  string = "embedding"
)

func InstantiateHostFunctions(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(HostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteDQL).Export("executeDQL")
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")
	b.NewFunctionBuilder().WithFunc(hostInvokeOpenaiChat).Export("invokeOpenaiChat")

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
	modelId, err := readString(mem, pModelId)
	if err != nil {
		err = fmt.Errorf("error reading model id from wasm memory: %w", err)
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

	modelSpec, err := dgraph.GetModelSpec(modelId)
	if err != nil {
		err = fmt.Errorf("error getting model endpoint: %w", err)
		return dgraph.ModelSpec{}, nil, err
	}

	return modelSpec, sentenceMap, nil
}

func postToModelEndpoint[TResult any](ctx context.Context, sentenceMap map[string]string, modelSpec dgraph.ModelSpec) (TResult, error) {

	key, err := aws.GetSecretString(ctx, modelSpec.ID)
	if err != nil {
		var result TResult
		return result, fmt.Errorf("error getting model key: %w", err)
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

func hostInvokeOpenaiChat(ctx context.Context, mod wasm.Module, pmodel uint32, pinstruction uint32, psentence uint32) uint32 {
	// invoke https://api.openai.com/v1/chat/completions

	mem := mod.Memory()
	model, err := readString(mem, pmodel)

	if err != nil {
		log.Print("error reading model name from wasm memory:", err)
		return 0
	}
	// model values "gpt-3.5-turbo" etc ...
	log.Print("reading model name from wasm memory:", model)
	sentence, err := readString(mem, psentence)
	if err != nil {
		log.Print("error reading sentence string from wasm memory:", err)
		return 0
	}
	instruction, err := readString(mem, pinstruction)
	if err != nil {
		log.Print("error reading instruction string from wasm memory:", err)
		return 0
	}
	// fetch the model
	//modelQuery := ""
	//r, err := queryGQL(ctx, modelQuery)
	//if err != nil {
	//	return "", fmt.Errorf("error getting GraphQL schema from Dgraph: %v", err)
	//}

	//var sr ModelResponse
	//err = json.Unmarshal(r, &sr)

	// POST to embedding service
	jinstruction, _ := json.Marshal(instruction)
	jsentence, _ := json.Marshal(sentence)
	// modelid to get secret is hardcoded to "openai"
	key, err := aws.GetSecretString(ctx, "openai")
	if err != nil {
		log.Print("error getting model key: %w", err)
		return 0
	}

	reqBody := `{ 
		"model": "` + model + `",
		"messages": [
		   { "role": "system", "content": ` + string(jinstruction) + `},
		   { "role": "user", "content": ` + string(jsentence) + `}
		 ]
	 }`
	log.Print("body: %v", reqBody)
	var buffer bytes.Buffer
	buffer.WriteString(reqBody)
	//Leverage Go's HTTP Post function to make request

	req, err := http.NewRequest(
		http.MethodPost,
		"https://api.openai.com/v1/chat/completions",
		&buffer,
	)
	if err != nil {
		log.Print("error buidling request:", err)
		return 0
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return 0
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return 0
	}
	// log.Printf("Openai response %s", string(body))
	// return a string
	return writeString(ctx, mod, string(body))
}

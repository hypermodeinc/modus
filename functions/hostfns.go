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
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const (
	HostModuleName  string = "hypermode"
	classifierModel string = "classifier"
	embeddingModel  string = "embedding"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func InstantiateHostFunctions(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(HostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteDQL).Export("executeDQL")
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")

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
		log.Println("error reading DQL statement from wasm memory:", err)
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		log.Println("error reading DQL variables string from wasm memory:", err)
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		log.Println("error unmarshaling GraphQL variables:", err)
		return 0
	}

	r, err := dgraph.ExecuteDQL(ctx, stmt, vars, isMutation != 0)
	if err != nil {
		log.Println("error executing DQL statement:", err)
		return 0
	}

	return writeString(ctx, mod, string(r))
}

func hostExecuteGQL(ctx context.Context, mod wasm.Module, pStmt uint32, pVars uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		log.Println("error reading GraphQL query string from wasm memory:", err)
		return 0
	}

	sVars, err := readString(mem, pVars)
	if err != nil {
		log.Println("error reading GraphQL variables string from wasm memory:", err)
		return 0
	}

	vars := make(map[string]string)
	if err := json.Unmarshal([]byte(sVars), &vars); err != nil {
		log.Println("error unmarshaling GraphQL variables:", err)
		return 0
	}

	r, err := dgraph.ExecuteGQL(ctx, stmt, vars)
	if err != nil {
		log.Println("error executing GraphQL operation:", err)
		return 0
	}

	return writeString(ctx, mod, string(r))
}

type ClassifierResult struct {
	Probabilities []ClassifierLabel `json:"probabilities"`
}

type ClassifierLabel struct {
	Label       string  `json:"label"`
	Probability float64 `json:"probability"`
}

func textModelSetup(mod wasm.Module, modelId uint32, psentenceMap uint32) (string, map[string]string, dgraph.ModelSpec, error) {
	mem := mod.Memory()
	mid, err := readString(mem, modelId)
	fmt.Println("reading model id from wasm memory:", mid)
	if err != nil {
		log.Println("error reading model id from wasm memory:", err)
		return "", map[string]string{}, dgraph.ModelSpec{}, err
	}
	sentenceMapStr, err := readString(mem, psentenceMap)
	if err != nil {
		log.Println("error reading sentence map string from wasm memory:", err)
		return "", map[string]string{}, dgraph.ModelSpec{}, err
	}

	sentenceMap := make(map[string]string)
	if err := json.Unmarshal([]byte(sentenceMapStr), &sentenceMap); err != nil {
		log.Println("error unmarshaling sentence map:", err)
		return "", map[string]string{}, dgraph.ModelSpec{}, err
	}

	spec, err := dgraph.GetModelSpec(mid)
	if err != nil {
		log.Println("error getting model endpoint:", err)
		return "", map[string]string{}, dgraph.ModelSpec{}, err
	}

	return mid, sentenceMap, spec, nil
}

func postToModelEndpoint(sentenceMap map[string]string, spec dgraph.ModelSpec, mid string) ([]byte, error) {
	// POST to model endpoint
	postBody, _ := json.Marshal(sentenceMap)
	requestBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request

	req, err := http.NewRequest(
		http.MethodPost,
		spec.Endpoint,
		requestBody)
	if err != nil {
		log.Println("error buidling request:", err)
		return []byte{}, err
	}
	svc, err := aws.GetSecretManagerSession()
	if err != nil {
		log.Println("error getting secret manager session:", err)
		return []byte{}, err
	}
	secretValue, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &mid,
	})
	if err != nil {
		log.Println("error getting secret:", err)
		return []byte{}, err
	}
	if secretValue.SecretString == nil {
		log.Println("secret string was empty")
		return []byte{}, err
	}

	modelKey := *secretValue.SecretString

	req.Header.Set("x-api-key", modelKey)
	resp, err := httpClient.Do(req)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return []byte{}, err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return []byte{}, err
	}
	return body, nil
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, modelId uint32, psentenceMap uint32) uint32 {
	mid, sentenceMap, spec, err := textModelSetup(mod, modelId, psentenceMap)

	if spec.Type != classifierModel {
		log.Println("error: model type is not classifier")
		return 0
	}

	body, err := postToModelEndpoint(sentenceMap, spec, mid)
	if err != nil {
		log.Println("error posting to model endpoint:", err)
		return 0
	}

	var result map[string]ClassifierResult
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Can not unmarshal JSON with error %v", err)
		return 0
	}

	if len(result) == 0 {
		log.Printf("Unexpected body returned from classifier, body: %v", string(body))
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Println("error marshalling classifier result:", err)
		return 0
	}
	// return a string
	return writeString(ctx, mod, string(res))
}

func hostComputeEmbedding(ctx context.Context, mod wasm.Module, modelId uint32, psentenceMap uint32) uint32 {
	mid, sentenceMap, spec, err := textModelSetup(mod, modelId, psentenceMap)

	if spec.Type != embeddingModel {
		log.Println("error: model type is not embedding")
		return 0
	}

	body, err := postToModelEndpoint(sentenceMap, spec, mid)
	if err != nil {
		log.Println("error posting to model endpoint:", err)
		return 0
	}

	var result map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("Can not unmarshal JSON with error %v", err)
		return 0
	}

	if len(result) == 0 {
		log.Printf("Unexpected body returned from embedding, body: %v", string(body))
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Println("error marshalling embedding result:", err)
		return 0
	}
	// return a string
	return writeString(ctx, mod, string(res))
}

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

const HostModuleName = "hypermode"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func InstantiateHostFunctions(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(HostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteDQL).Export("executeDQL")
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")

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

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, modelId uint32, psentenceMap uint32) uint32 {
	mem := mod.Memory()
	mid, err := readString(mem, modelId)
	fmt.Println("reading model id from wasm memory:", mid)
	if err != nil {
		log.Println("error reading model id from wasm memory:", err)
		return 0
	}
	sentenceMapStr, err := readString(mem, psentenceMap)
	if err != nil {
		log.Println("error reading sentence map string from wasm memory:", err)
		return 0
	}

	sentenceMap := make(map[string]string)
	if err := json.Unmarshal([]byte(sentenceMapStr), &sentenceMap); err != nil {
		log.Println("error unmarshaling sentence map:", err)
		return 0
	}

	endpoint, err := dgraph.GetModelEndpoint(mid)
	if err != nil {
		log.Println("error getting model endpoint:", err)
		return 0
	}

	// POST to model endpoint
	postBody, _ := json.Marshal(sentenceMap)
	requestBody := bytes.NewBuffer(postBody)
	//Leverage Go's HTTP Post function to make request

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		requestBody)
	if err != nil {
		log.Println("error buidling request:", err)
		return 0
	}
	svc, err := aws.GetSecretManagerSession()
	if err != nil {
		log.Println("error getting secret manager session:", err)
		return 0
	}
	secretValue, err := svc.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: &mid,
	})
	if err != nil {
		log.Println("error getting secret:", err)
		return 0
	}
	if secretValue.SecretString == nil {
		log.Println("secret string was empty")
		return 0
	}

	modelKey := *secretValue.SecretString

	req.Header.Set("x-api-key", modelKey)
	resp, err := httpClient.Do(req)

	//Handle Error
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return 0
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("An Error Occured %v", err)
		return 0
	}

	// snippet only
	var result map[string]ClassifierResult
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		log.Printf("Can not unmarshal JSON with error %v", err)
		return 0
	}

	if len(result) == 0 {
		log.Printf("Unexpected body returned from classifier, body: %v", string(body))
		return 0
	}

	res, err := json.Marshal(result)
	if err != nil {
		log.Println("error marshaling classifier result:", err)
		return 0
	}
	// return a string
	return writeString(ctx, mod, string(res))

}

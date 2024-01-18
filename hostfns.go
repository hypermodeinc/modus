/*
 * Copyright 2023 Hypermode, Inc.
 */
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"hmruntime/aws"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const HostModuleName = "hypermode"

func instantiateHostFunctions(ctx context.Context, runtime wazero.Runtime) error {
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

func hostExecuteDQL(ctx context.Context, mod wasm.Module, pStmt uint32, isMutation uint32) uint32 {
	mem := mod.Memory()
	stmt, err := readString(mem, pStmt)
	if err != nil {
		log.Println("error reading DQL statement from wasm memory:", err)
		return 0
	}

	r, err := executeDQL(ctx, stmt, isMutation != 0)
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
		log.Println("error reading GraphQL string from wasm memory:", err)
		return 0
	}

	varsBytes, err := readBuffer(mem, pVars)

	vars := make(map[string]string)
	if err := json.Unmarshal(varsBytes, &vars); err != nil {
		log.Println("error unmarshaling GraphQL variables:", err)
		return 0
	}

	r, err := executeGQL(ctx, stmt, vars)
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

type ClassifierResponse struct {
	Uid ClassifierResult `json:"uid"`
}

func hostInvokeClassifier(ctx context.Context, mod wasm.Module, modelId uint32, psentence uint32) uint32 {
	mem := mod.Memory()
	mid, err := readString(mem, modelId)
	if err != nil {
		log.Println("error reading model id from wasm memory:", err)
		return 0
	}
	fmt.Println("reading model id from wasm memory:", mid)
	sentence, err := readString(mem, psentence)
	if err != nil {
		log.Println("error reading sentence string from wasm memory:", err)
		return 0
	}

	endpoint, err := getModelEndpoint(mid)
	if err != nil {
		log.Println("error getting model endpoint:", err)
		return 0
	}

	// POST to model endpoint
	postBody, _ := json.Marshal(map[string]string{
		"uid": sentence,
	})
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
	var result ClassifierResponse
	if err := json.Unmarshal(body, &result); err != nil { // Parse []byte to go struct pointer
		fmt.Println("Can not unmarshal JSON")
	}
	str, _ := json.Marshal(result.Uid)
	// return a string
	return writeString(ctx, mod, string(str))

}

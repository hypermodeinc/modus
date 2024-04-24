/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/manifest"
	"hmruntime/models"

	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

const hostModuleName string = "hypermode"
const hypermodeHostName string = "hypermode"

func Instantiate(ctx context.Context, runtime wazero.Runtime) error {
	b := runtime.NewHostModuleBuilder(hostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")
	b.NewFunctionBuilder().WithFunc(hostInvokeTextGenerator).Export("invokeTextGenerator")

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", hostModuleName, err)
	}

	return nil
}

func getHost(mem wasm.Memory, pHostName uint32) (manifest.Host, error) {
	hostName, err := assemblyscript.ReadString(mem, pHostName)
	if err != nil {
		err = fmt.Errorf("error reading host name from wasm memory: %w", err)
		return manifest.Host{}, err
	}

	return hosts.GetHost(hostName)
}

func getModel(mem wasm.Memory, pModelName uint32, task manifest.ModelTask) (manifest.Model, error) {
	modelName, err := assemblyscript.ReadString(mem, pModelName)
	if err != nil {
		err = fmt.Errorf("error reading model name from wasm memory: %w", err)
		return manifest.Model{}, err
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

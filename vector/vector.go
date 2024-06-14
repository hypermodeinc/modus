package vector

import (
	"context"

	"hmruntime/vector/index/interfaces"
	"hmruntime/wasmhost/module"

	wasm "github.com/tetratelabs/wazero/api"
)

func ProcessTextMap(ctx context.Context, textIndex interfaces.TextIndex, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range textIndex.GetTextMap() {
		result, err := module.CallFunctionByName(ctx, embedder, text)
		if err != nil {
			return err
		}

		resultArr := result.([]interface{})

		textVec := make([]float64, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(float64)
		}

		_, err = vectorIndex.InsertVector(ctx, nil, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, textIndex interfaces.TextIndex, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range textIndex.GetTextMap() {
		result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			return err
		}

		resultArr := result.([]interface{})

		textVec := make([]float64, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(float64)
		}

		_, err = vectorIndex.InsertVector(ctx, nil, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

package vector

import (
	"context"

	c "hmruntime/vector/constraints"
	"hmruntime/vector/index"
	"hmruntime/wasmhost/module"

	wasm "github.com/tetratelabs/wazero/api"
)

func ProcessTextMap[T c.Float](ctx context.Context, textIndex index.TextIndex[T], embedder string, vectorIndex index.VectorIndex[T]) error {

	for uuid, text := range textIndex.GetTextMap() {
		result, err := module.CallFunctionByName(ctx, embedder, text)
		if err != nil {
			return err
		}

		resultArr := result.([]interface{})

		textVec := make([]T, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(T)
		}

		_, err = vectorIndex.InsertVector(ctx, nil, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessTextMapWithModule[T c.Float](ctx context.Context, mod wasm.Module, textIndex index.TextIndex[T], embedder string, vectorIndex index.VectorIndex[T]) error {

	for uuid, text := range textIndex.GetTextMap() {
		result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			return err
		}

		resultArr := result.([]interface{})

		textVec := make([]T, len(resultArr))
		for i, val := range resultArr {
			textVec[i] = val.(T)
		}

		_, err = vectorIndex.InsertVector(ctx, nil, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

package collections

import (
	"context"

	"hmruntime/collections/index/interfaces"
	"hmruntime/collections/utils"
	"hmruntime/wasmhost/module"

	wasm "github.com/tetratelabs/wazero/api"
)

type EmbedderFnCall struct {
	EmbedderFnName   string
	CollectionName   string
	SearchMethodName string
}

var FnCallChannel = make(chan EmbedderFnCall)

func ProcessTextMap(ctx context.Context, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range collection.GetTextMap() {
		result, err := module.CallFunctionByName(ctx, embedder, text)
		if err != nil {
			return err
		}

		textVec, err := utils.ConvertToFloat32Array(result)
		if err != nil {
			return err
		}

		_, err = vectorIndex.InsertVector(ctx, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

func ProcessTextMapWithModule(ctx context.Context, mod wasm.Module, collection interfaces.Collection, embedder string, vectorIndex interfaces.VectorIndex) error {

	for uuid, text := range collection.GetTextMap() {
		result, err := module.CallFunctionByNameWithModule(ctx, mod, embedder, text)
		if err != nil {
			return err
		}

		textVec, err := utils.ConvertToFloat32Array(result)
		if err != nil {
			return err
		}

		_, err = vectorIndex.InsertVector(ctx, uuid, textVec)
		if err != nil {
			return err
		}
	}
	return nil
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	"hmruntime/plugins"

	wasm "github.com/tetratelabs/wazero/api"
)

func readObject(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (any, error) {
	switch typ.Name {
	case "string":
		return ReadString(mem, offset)

	case "Date":
		return readDate(mem, offset)
	}

	def, err := getTypeDefinition(ctx, typ.Path)
	if err != nil {
		return nil, err
	}

	id, _ := mem.ReadUint32Le(offset - 8)
	if id != def.Id {
		return nil, fmt.Errorf("pointer is not to a %s", typ.Name)
	}

	if isArrayType(typ.Path) {
		return readArray(ctx, mem, def, offset)
	} else if isMapType(typ.Path) {
		return readMap(ctx, mem, def, offset)
	}

	return readClass(ctx, mem, def, offset)
}

func writeObject(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val any) (uint32, error) {
	// TODO: handle arrays and maps

	def, err := getTypeDefinition(ctx, typ.Path)
	if err != nil {
		return 0, err
	}

	switch v := val.(type) {
	case map[string]any:
		return writeClass(ctx, mod, def, v)
	}

	return 0, fmt.Errorf("unsupported object type")
}

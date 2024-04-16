/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"strings"

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

	if strings.HasPrefix(typ.Path, "~lib/array/Array<") {
		return readArray(ctx, mem, def, offset)
	}

	return readClass(ctx, mem, def, offset)
}

func writeObject(ctx context.Context, mod wasm.Module, bytes []byte, classId uint32) uint32 {
	offset := allocateWasmMemory(ctx, mod, len(bytes), classId)
	mod.Memory().Write(offset, bytes)
	return offset
}

func allocateWasmMemory(ctx context.Context, mod wasm.Module, len int, classId uint32) uint32 {
	// Allocate memory within the AssemblyScript module.
	// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
	// See https://www.assemblyscript.org/runtime.html#interface
	newFn := mod.ExportedFunction("__new")
	res, _ := newFn.Call(ctx, uint64(len), uint64(classId))
	return uint32(res[0])
}

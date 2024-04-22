/*
 * Copyright 2023 Hypermode, Inc.
 */

package assemblyscript

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/plugins"

	wasm "github.com/tetratelabs/wazero/api"
)

func readClass(ctx context.Context, mem wasm.Memory, def plugins.TypeDefinition, offset uint32) (map[string]any, error) {
	classMap := make(map[string]any)
	for _, field := range def.Fields {
		fieldOffset := offset + field.Offset
		val, err := readField(ctx, mem, field.Type, fieldOffset)
		if err != nil {
			return nil, err
		}
		classMap[field.Name] = val
	}
	return classMap, nil
}

func readField(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (any, error) {
	var result any
	var ok bool
	switch typ.Path {
	case "bool":
		var val byte
		val, ok = mem.ReadByte(offset)
		result = val != 0

	case "u8":
		result, ok = mem.ReadByte(offset)

	case "u16":
		result, ok = mem.ReadUint16Le(offset)

	case "u32":
		result, ok = mem.ReadUint32Le(offset)

	case "u64":
		result, ok = mem.ReadUint64Le(offset)

	case "i8":
		var val byte
		val, ok = mem.ReadByte(offset)
		result = int8(val)

	case "i16":
		var val uint16
		val, ok = mem.ReadUint16Le(offset)
		result = int16(val)

	case "i32":
		var val uint32
		val, ok = mem.ReadUint32Le(offset)
		result = int32(val)

	case "i64":
		var val uint64
		val, ok = mem.ReadUint64Le(offset)
		result = int64(val)

	case "f32":
		result, ok = mem.ReadFloat32Le(offset)

	case "f64":
		result, ok = mem.ReadFloat64Le(offset)

	default:
		// Managed types have pointers to the actual data.
		p, ok := mem.ReadUint32Le(offset)
		if !ok {
			return nil, fmt.Errorf("error reading %s pointer from wasm memory", typ.Name)
		}

		// Handle null values if the type is nullable
		if isNullable(typ.Path) {
			if p == 0 {
				return nil, nil
			}
			typ = plugins.TypeInfo{
				Name: typ.Name[:len(typ.Name)-7], // remove " | null"
				Path: typ.Path[:len(typ.Path)-5], // remove "|null"
			}
		} else if p == 0 {
			return nil, fmt.Errorf("null pointer encountered for non-nullable type %s", typ.Name)
		}

		// Read the actual data.
		return readObject(ctx, mem, typ, p)
	}

	if !ok {
		return nil, fmt.Errorf("error reading %s from wasm memory", typ.Name)
	}

	return result, nil
}

func writeClass(ctx context.Context, mod wasm.Module, def plugins.TypeDefinition, data any) (uint32, error) {
	switch data := data.(type) {
	case map[string]any:
		offset, err := allocateWasmMemory(ctx, mod, int(def.Size), def.Id)
		if err != nil {
			return 0, err
		}
		for _, field := range def.Fields {
			val := data[field.Name]
			fieldOffset := offset + field.Offset
			err := writeField(ctx, mod, field.Type, fieldOffset, val)
			if err != nil {
				return 0, err
			}
		}
		return offset, nil
	default:
		// Convert the data to a map and try again.
		j, err := json.Marshal(data)
		if err != nil {
			return 0, err
		}

		dec := json.NewDecoder(bytes.NewReader(j))
		dec.UseNumber()

		var m map[string]any
		err = dec.Decode(&m)
		if err != nil {
			return 0, err
		}
		return writeClass(ctx, mod, def, m)
	}
}

func writeField(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, offset uint32, val any) error {
	enc, err := EncodeValue(ctx, mod, typ, val)
	if err != nil {
		return err
	}

	mem := mod.Memory()

	var ok bool
	switch typ.Path {
	case "bool", "i8", "u8":
		ok = mem.WriteByte(offset, byte(enc))
	case "i16", "u16":
		ok = mem.WriteUint16Le(offset, uint16(enc))
	case "i64", "u64", "f64":
		ok = mem.WriteUint64Le(offset, enc)
	default:
		ok = mem.WriteUint32Le(offset, uint32(enc))
	}

	if !ok {
		return fmt.Errorf("error writing %s to wasm memory", typ.Name)
	}

	return nil
}

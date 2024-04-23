/*
 * Copyright 2023 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func readClass(ctx context.Context, mem wasm.Memory, def plugins.TypeDefinition, offset uint32) (data map[string]any, err error) {
	data = make(map[string]any)
	for _, field := range def.Fields {
		fieldOffset := offset + field.Offset
		val, err := readField(ctx, mem, field.Type, fieldOffset)
		if err != nil {
			return nil, err
		}
		data[field.Name] = val
	}
	return data, nil
}

func readField(ctx context.Context, mem wasm.Memory, typ plugins.TypeInfo, offset uint32) (data any, err error) {
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

func writeClass(ctx context.Context, mod wasm.Module, def plugins.TypeDefinition, data any) (offset uint32, err error) {
	switch data := data.(type) {
	case map[string]any:

		// unpin everything when done
		pins := make([]uint32, 0, len(def.Fields)+1)
		defer func() {
			for _, ptr := range pins {
				err = unpinWasmMemory(ctx, mod, ptr)
				if err != nil {
					break
				}
			}
		}()

		// Allocate memory for the object
		offset, err = allocateWasmMemory(ctx, mod, def.Size, def.Id)
		if err != nil {
			return 0, err
		}

		// we need to pin the object in memory so it doesn't get garbage collected
		// if we allocate more memory when writing a field before returning the object
		err = pinWasmMemory(ctx, mod, offset)
		if err != nil {
			return 0, err
		}
		pins = append(pins, offset)

		// write fields
		for _, field := range def.Fields {
			val := data[field.Name]
			fieldOffset := offset + field.Offset
			ptr, err := writeField(ctx, mod, field.Type, fieldOffset, val)
			if err != nil {
				return 0, err
			}

			// If we allocated memory for the field, we need to pin it too.
			if ptr != 0 {
				err = pinWasmMemory(ctx, mod, ptr)
				if err != nil {
					return 0, err
				}
				pins = append(pins, ptr)
			}
		}

		return offset, nil

	default:
		m, err := utils.ConvertToMap(data)
		if err != nil {
			return 0, err
		}
		return writeClass(ctx, mod, def, m)
	}
}

func writeField(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, offset uint32, val any) (ptr uint32, err error) {
	enc, err := EncodeValue(ctx, mod, typ, val)
	if err != nil {
		return 0, err
	}

	mem := mod.Memory()

	var ok bool
	switch typ.Path {
	case "bool", "i8", "u8":
		ok = mem.WriteByte(offset, byte(enc))
	case "i16", "u16":
		ok = mem.WriteUint16Le(offset, uint16(enc))
	case "i32", "u32", "f32":
		ok = mem.WriteUint32Le(offset, uint32(enc))
	case "i64", "u64", "f64":
		ok = mem.WriteUint64Le(offset, enc)
	default: // managed types
		ptr = uint32(enc) // return pointer to the managed object
		ok = mem.WriteUint32Le(offset, uint32(enc))
	}

	if !ok {
		return ptr, fmt.Errorf("error writing %s to wasm memory", typ.Name)
	}

	return ptr, nil
}

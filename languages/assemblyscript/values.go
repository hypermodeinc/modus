/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	wasm "github.com/tetratelabs/wazero/api"
)

func (wa *wasmAdapter) EncodeValue(ctx context.Context, mod wasm.Module, data any) (uint64, error) {

	rt := reflect.TypeOf(data)
	typ, err := wa.typeInfo.getAssemblyScriptType(rt)
	if err != nil {
		return 0, err
	}

	return wa.encodeValue(ctx, mod, typ, data)
}

func (wa *wasmAdapter) encodeValue(ctx context.Context, mod wasm.Module, typ string, data any) (uint64, error) {
	// For most calls, we don't need to pin the memory.
	// If it needs to be pinned, the caller will do it.
	return wa.doEncodeValue(ctx, mod, typ, data, false)
}

func (wa *wasmAdapter) encodeValueForParameter(ctx context.Context, mod wasm.Module, typ string, data any) (uint64, error) {
	// For the inbound parameters, we need to pin the memory.
	// Otherwise, just allocating more parameters could cause the GC to run and free the memory we just allocated.
	// Note we don't bother tracking these to unpin later, because we discard the module instance and its memory after the call.
	return wa.doEncodeValue(ctx, mod, typ, data, true)
}

func (wa *wasmAdapter) doEncodeValue(ctx context.Context, mod wasm.Module, typ string, data any, pin bool) (val uint64, err error) {

	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			}
		}
	}()

	// Handle null values if the type is nullable
	if wa.typeInfo.IsNullable(typ) {
		if data == nil {
			return 0, nil
		}
		rv := reflect.ValueOf(data)
		switch rv.Kind() {
		case reflect.Pointer, reflect.Interface, reflect.Slice, reflect.Map:
			if rv.IsNil() {
				return 0, nil
			}
		}

		typ = wa.typeInfo.GetUnderlyingType(typ)
	}

	// Primitive types
	switch typ {
	case "bool":
		b, ok := data.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "u8", "u16", "u32":
		var x uint32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = uint32(n)
		case uint:
			x = uint32(v)
		case uint8:
			x = uint32(v)
		case uint16:
			x = uint32(v)
		case uint32:
			x = uint32(v)
		default:
			return 0, fmt.Errorf("input value is not an unsigned integer")
		}

		return wasm.EncodeU32(x), nil

	case "i8", "i16", "i32":
		var x int32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = int32(n)
		case int:
			x = int32(v)
		case int8:
			x = int32(v)
		case int16:
			x = int32(v)
		case int32:
			x = v
		default:
			return 0, fmt.Errorf("input value is not an signed integer")
		}

		return wasm.EncodeI32(x), nil

	case "i64", "u64":
		var x int64
		switch v := data.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, fmt.Errorf("input value is not an int")
			}
			x = int64(n)
		case int64:
			x = v
		case uint64:
			x = int64(v)
		default:
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI64(x), nil

	case "f32":
		var x float32
		switch v := data.(type) {
		case json.Number:
			n, err := v.Float64()
			if err != nil {
				return 0, fmt.Errorf("input value is not a float")
			}
			x = float32(n)
		case float32:
			x = v
		default:
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF32(x), nil

	case "f64":
		var x float64
		switch v := data.(type) {
		case json.Number:
			n, err := v.Float64()
			if err != nil {
				return 0, fmt.Errorf("input value is not a float")
			}
			x = float64(n)
		case float64:
			x = v
		default:
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF64(x), nil
	}

	// Managed types need to be written to wasm memory
	offset, err := wa.writeObject(ctx, mod, typ, data)
	if err != nil {
		return 0, err
	}

	// Pin the memory if requested
	if pin {
		err := wa.pinWasmMemory(ctx, mod, offset)
		if err != nil {
			return 0, fmt.Errorf("failed to pin wasm memory: %w", err)
		}
	}

	return uint64(offset), nil
}

func (wa *wasmAdapter) DecodeValue(ctx context.Context, mod wasm.Module, val uint64, data any) error {
	if val == 0 {
		return nil
	}

	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("expected pointer, got %T", data)
	}

	rvElem := rv.Elem()

	typ, err := wa.typeInfo.getAssemblyScriptType(rvElem.Type())
	if err != nil {
		return err
	}

	d, err := wa.decodeValue(ctx, mod, typ, val)
	if err != nil {
		return err
	}

	if rvElem.Kind() == reflect.Struct {
		if m, ok := d.(map[string]any); ok {
			if err := mapToStruct(m, rv.Interface()); err != nil {
				return err
			}
			return nil
		}
	}

	rd := reflect.ValueOf(d)
	if rd.Type() == rvElem.Type() {
		rvElem.Set(rd)
		return nil
	}

	return fmt.Errorf("expected %s, got %s", rvElem.Type(), rd.Type())
}

func (wa *wasmAdapter) decodeValueForResult(ctx context.Context, mod wasm.Module, typ string, val uint64) (data any, err error) {
	return wa.decodeValue(ctx, mod, typ, val)
}

func (wa *wasmAdapter) decodeValue(ctx context.Context, mod wasm.Module, typ string, val uint64) (data any, err error) {

	// Recover from panics and convert them to errors
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			} else if e, ok := r.(string); ok {
				err = errors.New(e)
			}
		}
	}()

	// Handle null values if the type is nullable
	if wa.typeInfo.IsNullable(typ) {
		if val == 0 {
			return nil, nil
		}
		typ = wa.typeInfo.GetUnderlyingType(typ)
	}

	// Primitive types
	switch typ {
	case "void":
		return 0, nil

	case "bool":
		return val != 0, nil

	case "u8", "u16", "u32":
		return wasm.DecodeU32(val), nil

	case "i8", "i16", "i32":
		return wasm.DecodeI32(val), nil

	case "i64", "u64":
		return int64(val), nil

	case "f32":
		return wasm.DecodeF32(val), nil

	case "f64":
		return wasm.DecodeF64(val), nil
	}

	// Managed types are read from wasm memory
	mem := mod.Memory()
	return wa.readObject(ctx, mem, typ, uint32(val))
}

func mapToStruct(m map[string]any, result any) error {

	config := &mapstructure.DecoderConfig{
		Result: result,
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(m)
}

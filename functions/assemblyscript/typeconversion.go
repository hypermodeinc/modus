/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

func EncodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val any) (uint64, error) {

	switch typ.Name {

	case "bool":
		b, ok := val.(bool)
		if !ok {
			return 0, fmt.Errorf("input value is not a bool")
		}

		if b {
			return 1, nil
		} else {
			return 0, nil
		}

	case "i8", "i16", "i32", "u8", "u16", "u32":
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI32(int32(n)), nil

	case "i64", "u64":
		n, err := val.(json.Number).Int64()
		if err != nil {
			return 0, fmt.Errorf("input value is not an int")
		}

		return wasm.EncodeI64(n), nil

	case "f32":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF32(float32(n)), nil

	case "f64":
		n, err := val.(json.Number).Float64()
		if err != nil {
			return 0, fmt.Errorf("input value is not a float")
		}

		return wasm.EncodeF64(n), nil

	case "string":
		s, ok := val.(string)
		if !ok {
			return 0, fmt.Errorf("input value is not a string")
		}

		// Note, strings are passed as a pointer to a string in wasm memory
		ptr := WriteString(ctx, mod, s)
		return uint64(ptr), nil

	case "Date":
		var t time.Time
		switch v := val.(type) {
		case json.Number:
			n, err := v.Int64()
			if err != nil {
				return 0, err
			}
			t = time.UnixMilli(n)
		case string:
			var err error
			t, err = utils.ParseTime(v)
			if err != nil {
				return 0, err
			}
		}

		ptr, err := writeDate(ctx, mod, t)
		if err != nil {
			return 0, err
		}

		return uint64(ptr), nil

	// TODO: custom types

	default:
		return 0, fmt.Errorf("unknown parameter type: %s", typ)
	}
}

func DecodeValue(ctx context.Context, mod wasm.Module, typ plugins.TypeInfo, val uint64) (any, error) {

	// Handle null values if the type is nullable
	if strings.HasSuffix(typ.Path, "|null") {
		if val == 0 {
			return nil, nil
		}
		typ = plugins.TypeInfo{
			Name: typ.Name[:len(typ.Name)-7], // remove " | null"
			Path: typ.Path[:len(typ.Path)-5], // remove "|null"
		}
	}

	// Primitive types
	switch typ.Path {
	case "void":
		return 0, nil

	case "bool":
		return val != 0, nil

	case "i8", "i16", "i32", "u8", "u16", "u32":
		return wasm.DecodeI32(val), nil

	case "i64", "u64":
		return int64(val), nil

	case "f32":
		return wasm.DecodeF32(val), nil

	case "f64":
		return wasm.DecodeF64(val), nil
	}

	// Managed types, held in wasm memory
	mem := mod.Memory()
	return readObject(ctx, mem, typ, uint32(val))
}

func isPrimitive(t string) bool {
	switch t {
	case "bool", "i8", "i16", "i32", "u8", "u16", "u32", "i64", "u64", "f32", "f64":
		return true
	default:
		return false
	}
}

var typeMap = map[string]string{
	"~lib/string/String":       "string",
	"~lib/array/Array":         "Array",
	"~lib/map/Map":             "Map",
	"~lib/date/Date":           "Date",
	"~lib/wasi_date/wasi_Date": "Date",
}

var mapRegex = regexp.MustCompile(`^~lib/map/Map<(\w+<.+>|.+),\s*(\w+<.+>|.+)>$`)

func getTypeInfo(path string) plugins.TypeInfo {

	var name string
	if isPrimitive(path) {
		name = path
	} else if t, ok := typeMap[path]; ok {
		name = t
	} else if strings.HasSuffix(path, "|null") {
		name = getTypeInfo(path[:len(path)-5]).Name + " | null"
	} else if strings.HasPrefix(path, "~lib/array/Array<") {
		t := getTypeInfo(path[17 : len(path)-1])
		if strings.HasSuffix(t.Path, "|null") {
			name = "(" + t.Name + ")[]"
		} else {
			name += t.Name + "[]"
		}
	} else if strings.HasPrefix(path, "~lib/map/Map<") {
		matches := mapRegex.FindStringSubmatch(path)
		kt := getTypeInfo(matches[1])
		vt := getTypeInfo(matches[2])
		name = "Map<" + kt.Name + ", " + vt.Name + ">"
	} else {
		name = path[strings.LastIndex(path, "/")+1:]
	}

	return plugins.TypeInfo{
		Name: name,
		Path: path,
	}
}

func getTypeDefinition(ctx context.Context, typePath string) (plugins.TypeDefinition, error) {
	if isPrimitive(typePath) {
		return plugins.TypeDefinition{
			Name: typePath,
			Path: typePath,
		}, nil
	}

	plugin := ctx.Value(utils.PluginContextKey).(*plugins.Plugin)

	types := plugin.Types
	info, ok := types[typePath]
	if !ok {
		return plugins.TypeDefinition{}, fmt.Errorf("info for type %s not found in plugin %s", typePath, plugin.Name())
	}

	return info, nil
}

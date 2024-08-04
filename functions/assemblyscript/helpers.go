/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"hmruntime/plugins"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

// These are the managed types that we handle directly.
var typeMap = map[string]string{
	"~lib/arraybuffer/ArrayBuffer": "ArrayBuffer",
	"~lib/string/String":           "string",
	"~lib/array/Array":             "Array",
	"~lib/map/Map":                 "Map",
	"~lib/date/Date":               "Date",
	"~lib/wasi_date/wasi_Date":     "Date",
}

var arrayBufferType = plugins.TypeInfo{
	Name: "ArrayBuffer",
	Path: "~lib/arraybuffer/ArrayBuffer",
}

var stringType = plugins.TypeInfo{
	Name: "string",
	Path: "~lib/string/String",
}

// Allocate memory within the AssemblyScript module.
// This uses the `__new` function exported by the AssemblyScript runtime, so it will be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func allocateWasmMemory(ctx context.Context, mod wasm.Module, len uint32, classId uint32) (uint32, error) {
	fn := mod.ExportedFunction("__new")
	res, err := fn.Call(ctx, uint64(len), uint64(classId))
	if err != nil {
		return 0, fmt.Errorf("failed to allocate WASM memory: %w", err)
	}
	return uint32(res[0]), nil
}

// Pin a managed object in memory within the AssemblyScript module.
// This prevents it from being garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func pinWasmMemory(ctx context.Context, mod wasm.Module, ptr uint32) error {
	fn := mod.ExportedFunction("__pin")
	_, err := fn.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to pin object in WASM memory: %w", err)
	}
	return nil
}

// Unpin a previously-pinned managed object in memory within the AssemblyScript module.
// This allows it to be garbage collected.
// See https://www.assemblyscript.org/runtime.html#interface
func unpinWasmMemory(ctx context.Context, mod wasm.Module, ptr uint32) error {
	fn := mod.ExportedFunction("__unpin")
	_, err := fn.Call(ctx, uint64(ptr))
	if err != nil {
		return fmt.Errorf("failed to unpin object in WASM memory: %w", err)
	}
	return nil
}

// Sets that arguments length before calling a function that takes a variable number of arguments.
// See https://www.assemblyscript.org/runtime.html#optional-arguments
func SetArgumentsLength(ctx context.Context, mod wasm.Module, length int) error {
	fn := mod.ExportedFunction("__setArgumentsLength")
	_, err := fn.Call(ctx, uint64(length))
	if err != nil {
		return fmt.Errorf("failed to set arguments length: %w", err)
	}
	return nil
}

func isPrimitive(t string) bool {
	switch t {
	case "bool", "i8", "i16", "i32", "u8", "u16", "u32", "i64", "u64", "f32", "f64":
		return true
	default:
		return false
	}
}

func getItemSize(typ plugins.TypeInfo) uint32 {
	switch typ.Path {
	case "u64", "i64", "f64":
		return 8
	case "u32", "i32", "f32":
		return 4
	case "u16", "i16":
		return 2
	case "u8", "i8", "bool":
		return 1
	default:
		return 4 // pointer
	}
}

func isArrayType(path string) bool {
	return strings.HasPrefix(path, "~lib/array/Array<")
}

func isMapType(path string) bool {
	return strings.HasPrefix(path, "~lib/map/Map<")
}

func isNullable(t string) bool {
	return strings.HasSuffix(t, "|null") || strings.HasSuffix(t, " | null")
}

func getArraySubtypeInfo(path string) plugins.TypeInfo {
	return getTypeInfo(path[17 : len(path)-1])
}

func getMapSubtypeInfo(path string) (plugins.TypeInfo, plugins.TypeInfo) {
	k, v := GetMapParts(path)
	return getTypeInfo(k), getTypeInfo(v)
}

func GetMapParts(t string) (string, string) {
	prefix := "~lib/map/Map<"
	if !strings.HasPrefix(t, prefix) {
		prefix = "Map<"
		if !strings.HasPrefix(t, prefix) {
			return "", ""
		}
	}

	n := 1
	c := 0
	for i := len(prefix); i < len(t); i++ {
		switch t[i] {
		case '<':
			n++
		case ',':
			if n == 1 {
				c = i
			}
		case '>':
			n--
			if n == 0 {
				k := strings.TrimSpace(t[len(prefix):c])
				v := strings.TrimSpace(t[c+1 : i])
				return k, v
			}
		}
	}

	return "", ""
}

func getTypeInfo(path string) plugins.TypeInfo {
	var name string
	if isPrimitive(path) {
		name = path
	} else if t, ok := typeMap[path]; ok {
		name = t
	} else if strings.HasSuffix(path, "|null") {
		name = getTypeInfo(path[:len(path)-5]).Name + " | null"
	} else if isArrayType(path) {
		t := getArraySubtypeInfo(path)
		if strings.HasSuffix(t.Path, "|null") {
			name = "(" + t.Name + ")[]"
		} else {
			name += t.Name + "[]"
		}
	} else if isMapType(path) {
		kt, vt := getMapSubtypeInfo(path)
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

func getTypeInfoForType[T any]() (plugins.TypeInfo, error) {
	var v T

	switch any(v).(type) {
	case bool:
		return getTypeInfo("bool"), nil
	case int8:
		return getTypeInfo("i8"), nil
	case int16:
		return getTypeInfo("i16"), nil
	case int32, int:
		return getTypeInfo("i32"), nil
	case int64:
		return getTypeInfo("i64"), nil
	case uint8:
		return getTypeInfo("u8"), nil
	case uint16:
		return getTypeInfo("u16"), nil
	case uint32, uint:
		return getTypeInfo("u32"), nil
	case uint64:
		return getTypeInfo("u64"), nil
	case float32:
		return getTypeInfo("f32"), nil
	case float64:
		return getTypeInfo("f64"), nil
	case []byte:
		return arrayBufferType, nil
	case string, *string:
		return stringType, nil
	}

	t := reflect.TypeFor[T]()
	return getTypeInfoForReflectedType(t)
}

func getTypeInfoForReflectedType(t reflect.Type) (plugins.TypeInfo, error) {
	switch t.Kind() {
	case reflect.Bool:
		return getTypeInfo("bool"), nil
	case reflect.Int8:
		return getTypeInfo("i8"), nil
	case reflect.Int16:
		return getTypeInfo("i16"), nil
	case reflect.Int32, reflect.Int:
		return getTypeInfo("i32"), nil
	case reflect.Int64:
		return getTypeInfo("i64"), nil
	case reflect.Uint8:
		return getTypeInfo("u8"), nil
	case reflect.Uint16:
		return getTypeInfo("u16"), nil
	case reflect.Uint32, reflect.Uint:
		return getTypeInfo("u32"), nil
	case reflect.Uint64:
		return getTypeInfo("u64"), nil
	case reflect.Float32:
		return getTypeInfo("f32"), nil
	case reflect.Float64:
		return getTypeInfo("f64"), nil
	case reflect.String:
		return stringType, nil

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return arrayBufferType, nil
		}

		elemType, err := getTypeInfoForReflectedType(t.Elem())
		if err != nil {
			return plugins.TypeInfo{}, err
		}

		return plugins.TypeInfo{
			Name: elemType.Name + "[]",
			Path: "~lib/array/Array<" + elemType.Path + ">",
		}, nil

	case reflect.Map:
		keyType, err := getTypeInfoForReflectedType(t.Key())
		if err != nil {
			return plugins.TypeInfo{}, err
		}

		valueType, err := getTypeInfoForReflectedType(t.Elem())
		if err != nil {
			return plugins.TypeInfo{}, err
		}

		return plugins.TypeInfo{
			Name: "Map<" + keyType.Name + ", " + valueType.Name + ">",
			Path: "~lib/map/Map<" + keyType.Path + "," + valueType.Path + ">",
		}, nil
	case reflect.Ptr:
		return getTypeInfoForReflectedType(t.Elem())
	case reflect.Struct:
		id := t.PkgPath() + "." + t.Name()
		info, ok := hostTypes[id]
		if ok {
			return info, nil
		}

		return plugins.TypeInfo{}, fmt.Errorf("struct missing from host types map: %s", id)
	}

	return plugins.TypeInfo{}, fmt.Errorf("unsupported type kind %s", t.Kind())
}

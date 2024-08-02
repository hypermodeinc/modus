/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"hmruntime/plugins/metadata"
	"hmruntime/utils"
	"reflect"
	"strings"
	"time"
)

func (ti *typeInfoProvider) GetArraySubtype(t string) string {
	if strings.HasPrefix(t, "~lib/array/Array<") {
		return t[17 : len(t)-1]
	}

	if strings.HasSuffix(t, "[]") {
		return t[:len(t)-2]
	}

	return ""
}

func (ti *typeInfoProvider) GetMapSubtypes(t string) (string, string) {
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
				return t[len(prefix):c], t[c+1 : i]
			}
		}
	}

	return "", ""
}

func (ti *typeInfoProvider) GetNameForType(t string) string {
	s := ti.GetUnderlyingType(t)
	return s[strings.LastIndex(s, "/")+1:]
}

func (ti *typeInfoProvider) GetUnderlyingType(t string) string {
	return strings.TrimSuffix(t, "|null")
}

func (ti *typeInfoProvider) IsArrayType(t string) bool {
	return strings.HasPrefix(t, "~lib/array/Array<") || strings.HasSuffix(t, "[]")
}

func (ti *typeInfoProvider) IsBooleanType(t string) bool {
	return t == "bool"
}

func (ti *typeInfoProvider) IsByteSequenceType(t string) bool {
	switch t {
	case
		"~lib/arraybuffer/ArrayBuffer",
		"~lib/typedarray/Uint8Array",
		"~lib/array/Array<u8>", "u8[]":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsFloatType(t string) bool {
	switch t {
	case "f32", "f64":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsIntegerType(t string) bool {
	switch t {
	case "i8", "i16", "i32", "i64",
		"u8", "u16", "u32", "u64",
		"isize", "usize":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsMapType(t string) bool {
	return strings.HasPrefix(t, "~lib/map/Map<") || strings.HasPrefix(t, "Map<")
}

func (ti *typeInfoProvider) IsNullable(t string) bool {
	return strings.HasSuffix(t, "|null")
}

func (ti *typeInfoProvider) IsSignedIntegerType(t string) bool {
	switch t {
	case "i8", "i16", "i32", "i64", "isize":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsStringType(t string) bool {
	switch t {
	case "string", "~lib/string/String":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsTimestampType(t string) bool {
	switch t {
	case "Date", "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) SizeOfType(t string) uint32 {
	switch t {
	case "u64", "i64", "f64":
		return 8
	case "u16", "i16":
		return 2
	case "u8", "i8", "bool":
		return 1
	default:
		// 32-bit types, including primitive types and pointers to managed objects
		// we only support wasm32, so we can assume 32-bit pointers
		return 4
	}
}

func (ti *typeInfoProvider) getTypeDefinition(ctx context.Context, typeName string) (*metadata.TypeDefinition, error) {

	md := ctx.Value(utils.MetadataContextKey).(*metadata.Metadata)
	typ, ok := md.Types[typeName]
	if !ok {
		return nil, fmt.Errorf("info for type %s not found in plugin %s", typeName, md.Name())
	}

	return typ, nil
}

func (ti *typeInfoProvider) getAssemblyScriptType(t reflect.Type) (string, error) {
	switch t.Kind() {
	case reflect.Bool:
		return "bool", nil
	case reflect.Int:
		return "isize", nil
	case reflect.Int8:
		return "i8", nil
	case reflect.Int16:
		return "i16", nil
	case reflect.Int32:
		return "i32", nil
	case reflect.Int64:
		return "i64", nil
	case reflect.Uint:
		return "usize", nil
	case reflect.Uint8:
		return "u8", nil
	case reflect.Uint16:
		return "u16", nil
	case reflect.Uint32:
		return "u32", nil
	case reflect.Uint64:
		return "u64", nil
	case reflect.Float32:
		return "f32", nil
	case reflect.Float64:
		return "f64", nil
	case reflect.String:
		return "string", nil

	case reflect.Slice:
		if t.Elem().Kind() == reflect.Uint8 {
			return "~lib/arraybuffer/ArrayBuffer", nil
		}

		// TODO: should we use Uint8Array and other typed array buffers?

		elemType, err := ti.getAssemblyScriptType(t.Elem())
		if err == nil {
			return "", err
		}

		return "~lib/array/Array<" + elemType + ">", nil

	case reflect.Map:
		keyType, err := ti.getAssemblyScriptType(t.Key())
		if err != nil {
			return "", err
		}

		valueType, err := ti.getAssemblyScriptType(t.Elem())
		if err != nil {
			return "", err
		}

		return "~lib/map/Map<" + keyType + "," + valueType + ">", nil

	case reflect.Ptr:
		return ti.getAssemblyScriptType(t.Elem())

	case reflect.Struct:
		id := t.PkgPath() + "." + t.Name()
		typ, found := hostTypes[id]
		if found {
			return typ, nil
		}

		return "", fmt.Errorf("struct missing from host types map: %s", id)
	}

	return "", fmt.Errorf("unsupported type kind %s", t.Kind())
}

func (ti *typeInfoProvider) getGoType(asType string) (reflect.Type, error) {
	if goType, ok := asToGoTypeMap[asType]; ok {
		return goType, nil
	}

	if ti.IsArrayType(asType) {
		et := ti.GetArraySubtype(asType)
		if et == "" {
			return nil, fmt.Errorf("invalid array type: %s", asType)
		}

		elementType, err := ti.getGoType(et)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if ti.IsMapType(asType) {
		kt, vt := ti.GetMapSubtypes(asType)
		if kt == "" || vt == "" {
			return nil, fmt.Errorf("invalid map type: %s", asType)
		}

		keyType, err := ti.getGoType(kt)
		if err != nil {
			return nil, err
		}
		valType, err := ti.getGoType(vt)
		if err != nil {
			return nil, err
		}

		return reflect.MapOf(keyType, valType), nil
	}

	// All other types are custom classes, which are represented as a map[string]any
	return goMapStringAnyType, nil
}

var goMapStringAnyType = reflect.TypeOf(map[string]any{})

// map all AssemblyScript types to Go types
var asToGoTypeMap = map[string]reflect.Type{
	"bool":                              reflect.TypeOf(bool(false)),
	"usize":                             reflect.TypeOf(uint(0)),
	"u8":                                reflect.TypeOf(uint8(0)),
	"u16":                               reflect.TypeOf(uint16(0)),
	"u32":                               reflect.TypeOf(uint32(0)),
	"u64":                               reflect.TypeOf(uint64(0)),
	"isize":                             reflect.TypeOf(int(0)),
	"i8":                                reflect.TypeOf(int8(0)),
	"i16":                               reflect.TypeOf(int16(0)),
	"i32":                               reflect.TypeOf(int32(0)),
	"i64":                               reflect.TypeOf(int64(0)),
	"f32":                               reflect.TypeOf(float32(0)),
	"f64":                               reflect.TypeOf(float64(0)),
	"string":                            reflect.TypeOf(string("")),
	"~lib/string/String":                reflect.TypeOf(string("")),
	"~lib/arraybuffer/ArrayBuffer":      reflect.TypeOf([]byte{}),
	"~lib/typedarray/Uint8ClampedArray": reflect.TypeOf([]uint8{}),
	"~lib/typedarray/Uint8Array":        reflect.TypeOf([]uint8{}),
	"~lib/typedarray/Uint16Array":       reflect.TypeOf([]uint16{}),
	"~lib/typedarray/Uint32Array":       reflect.TypeOf([]uint32{}),
	"~lib/typedarray/Uint64Array":       reflect.TypeOf([]uint64{}),
	"~lib/typedarray/Int8Array":         reflect.TypeOf([]int8{}),
	"~lib/typedarray/Int16Array":        reflect.TypeOf([]int16{}),
	"~lib/typedarray/Int32Array":        reflect.TypeOf([]int32{}),
	"~lib/typedarray/Int64Array":        reflect.TypeOf([]int64{}),
	"~lib/typedarray/Float32Array":      reflect.TypeOf([]float32{}),
	"~lib/typedarray/Float64Array":      reflect.TypeOf([]float64{}),
	"~lib/date/Date":                    reflect.TypeOf(time.Time{}),
	"~lib/wasi_date/wasi_Date":          reflect.TypeOf(time.Time{}),
}

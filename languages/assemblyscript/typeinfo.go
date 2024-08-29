/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
	"reflect"
	"strings"
	"time"
)

func (ti *typeInfoProvider) GetListSubtype(typ string) string {
	if strings.HasPrefix(typ, "~lib/array/Array<") {
		return typ[17 : len(typ)-1]
	}

	if strings.HasSuffix(typ, "[]") {
		return typ[:len(typ)-2]
	}

	return ""
}

func (ti *typeInfoProvider) GetMapSubtypes(typ string) (string, string) {
	prefix := "~lib/map/Map<"
	if !strings.HasPrefix(typ, prefix) {
		prefix = "Map<"
		if !strings.HasPrefix(typ, prefix) {
			return "", ""
		}
	}

	n := 1
	c := 0
	for i := len(prefix); i < len(typ); i++ {
		switch typ[i] {
		case '<':
			n++
		case ',':
			if n == 1 {
				c = i
			}
		case '>':
			n--
			if n == 0 {
				return typ[len(prefix):c], typ[c+1 : i]
			}
		}
	}

	return "", ""
}

func (ti *typeInfoProvider) GetNameForType(typ string) string {
	s := ti.GetUnderlyingType(typ)

	if ti.IsListType(s) {
		return ti.GetNameForType(ti.GetListSubtype(s)) + "[]"
	}

	if ti.IsMapType(s) {
		kt, vt := ti.GetMapSubtypes(s)
		return "Map<" + ti.GetNameForType(kt) + "," + ti.GetNameForType(vt) + ">"
	}

	return s[strings.LastIndex(s, "/")+1:]
}

func (ti *typeInfoProvider) GetUnderlyingType(typ string) string {
	return strings.TrimSuffix(typ, "|null")
}

func (ti *typeInfoProvider) IsListType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/array/Array<") || strings.HasSuffix(typ, "[]")
}

func (ti *typeInfoProvider) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (ti *typeInfoProvider) IsByteSequenceType(typ string) bool {
	switch typ {
	case
		"~lib/arraybuffer/ArrayBuffer",
		"~lib/typedarray/Uint8Array",
		"~lib/array/Array<u8>", "u8[]":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsFloatType(typ string) bool {
	switch typ {
	case "f32", "f64":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64",
		"u8", "u16", "u32", "u64",
		"isize", "usize":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsMapType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/map/Map<") || strings.HasPrefix(typ, "Map<")
}

func (ti *typeInfoProvider) IsNullable(typ string) bool {
	return strings.HasSuffix(typ, "|null")
}

func (ti *typeInfoProvider) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64", "isize":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsStringType(typ string) bool {
	switch typ {
	case "string", "~lib/string/String":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsTimestampType(typ string) bool {
	switch typ {
	case "Date", "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
	switch typ {
	case "u64", "i64", "f64":
		return 8, nil
	case "u16", "i16":
		return 2, nil
	case "u8", "i8", "bool":
		return 1, nil
	default:
		// 32-bit types, including primitive types and pointers to managed objects
		// we only support wasm32, so we can assume 32-bit pointers
		return 4, nil
	}
}

func (ti *typeInfoProvider) GetTypeDefinition(ctx context.Context, typ string) (*metadata.TypeDefinition, error) {

	md := ctx.Value(utils.MetadataContextKey).(*metadata.Metadata)
	def, ok := md.Types[typ]
	if !ok {
		return nil, fmt.Errorf("info for type %s not found in plugin %s", typ, md.Name())
	}

	return def, nil
}

func (ti *typeInfoProvider) getAssemblyScriptType(rt reflect.Type, customTypes map[reflect.Type]string) (string, error) {
	if customTypes != nil {
		if typ, ok := customTypes[rt]; ok {
			return typ, nil
		}
	}

	switch rt.Kind() {
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
		return "~lib/string/String", nil

	case reflect.Slice:
		if rt.Elem().Kind() == reflect.Uint8 {
			return "~lib/arraybuffer/ArrayBuffer", nil
		}

		// TODO: should we use Uint8Array and other typed array buffers?

		elemType, err := ti.getAssemblyScriptType(rt.Elem(), customTypes)
		if err != nil {
			return "", err
		}

		return "~lib/array/Array<" + elemType + ">", nil

	case reflect.Map:
		keyType, err := ti.getAssemblyScriptType(rt.Key(), customTypes)
		if err != nil {
			return "", err
		}

		valueType, err := ti.getAssemblyScriptType(rt.Elem(), customTypes)
		if err != nil {
			return "", err
		}

		return "~lib/map/Map<" + keyType + "," + valueType + ">", nil

	case reflect.Ptr:
		return ti.getAssemblyScriptType(rt.Elem(), customTypes)

	case reflect.Struct:
		id := rt.PkgPath() + "." + rt.Name()
		typ, found := hostTypes[id]
		if found {
			return typ, nil
		}

		return "", fmt.Errorf("struct missing from host types map: %s", id)
	}

	return "", fmt.Errorf("unsupported type kind %s", rt.Kind())
}

func (ti *typeInfoProvider) getReflectedType(typ string, customTypes map[string]reflect.Type) (reflect.Type, error) {
	if customTypes != nil {
		if rt, ok := customTypes[typ]; ok {
			return rt, nil
		}
	}

	if rt, ok := asToGoTypeMap[typ]; ok {
		return rt, nil
	}

	if ti.IsListType(typ) {
		et := ti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid array type: %s", typ)
		}

		elementType, err := ti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if ti.IsMapType(typ) {
		kt, vt := ti.GetMapSubtypes(typ)
		if kt == "" || vt == "" {
			return nil, fmt.Errorf("invalid map type: %s", typ)
		}

		keyType, err := ti.getReflectedType(kt, customTypes)
		if err != nil {
			return nil, err
		}
		valType, err := ti.getReflectedType(vt, customTypes)
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

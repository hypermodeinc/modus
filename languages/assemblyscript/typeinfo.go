/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"hypruntime/langsupport"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

var _typeInfo = &typeInfo{}

func TypeInfo() langsupport.TypeInfo {
	return _typeInfo
}

type typeInfo struct{}

func (ti *typeInfo) GetListSubtype(typ string) string {
	typ = ti.GetUnderlyingType(typ)
	if strings.HasPrefix(typ, "~lib/array/Array<") {
		return typ[17 : len(typ)-1]
	}

	return ""
}

func (ti *typeInfo) GetMapSubtypes(typ string) (string, string) {
	typ = ti.GetUnderlyingType(typ)
	prefix := "~lib/map/Map<"
	if !strings.HasPrefix(typ, prefix) {
		return "", ""
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

func (ti *typeInfo) GetNameForType(typ string) string {
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

func (ti *typeInfo) GetUnderlyingType(typ string) string {
	if s, ok := strings.CutSuffix(typ, " | null"); ok {
		return s
	} else {
		return strings.TrimSuffix(typ, "|null")
	}
}

func (ti *typeInfo) IsNullable(typ string) bool {
	return strings.HasSuffix(typ, " | null") || strings.HasSuffix(typ, "|null")
}

func (ti *typeInfo) IsListType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/array/Array<")
}

func (ti *typeInfo) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (ti *typeInfo) IsByteSequenceType(typ string) bool {
	typ = ti.GetUnderlyingType(typ)
	switch typ {
	case
		"~lib/arraybuffer/ArrayBuffer",
		"~lib/typedarray/Uint8Array",
		"~lib/typedarray/Uint8ClampedArray",
		"~lib/array/Array<u8>":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsFloatType(typ string) bool {
	switch typ {
	case "f32", "f64":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64",
		"u8", "u16", "u32", "u64",
		"isize", "usize":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsMapType(typ string) bool {
	typ = ti.GetUnderlyingType(typ)
	return strings.HasPrefix(typ, "~lib/map/Map<")
}

func (ti *typeInfo) IsPrimitiveType(typ string) bool {
	return ti.IsBooleanType(typ) || ti.IsIntegerType(typ) || ti.IsFloatType(typ)
}

func (ti *typeInfo) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64", "isize":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsStringType(typ string) bool {
	return ti.GetUnderlyingType(typ) == "~lib/string/String"
}

func (ti *typeInfo) IsArrayBufferType(typ string) bool {
	return ti.GetUnderlyingType(typ) == "~lib/arraybuffer/ArrayBuffer"
}

func (ti *typeInfo) IsTypedArrayType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/typedarray/")
}

func (ti *typeInfo) IsTimestampType(typ string) bool {
	typ = ti.GetUnderlyingType(typ)
	switch typ {
	case "Date", "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
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

// Deprecated: use metadata.Metadata.GetTypeDefinition instead
func (ti *typeInfo) GetTypeDefinition(ctx context.Context, typ string) (*metadata.TypeDefinition, error) {
	md := ctx.Value(utils.MetadataContextKey).(*metadata.Metadata)
	return md.GetTypeDefinition(typ)
}

func (ti *typeInfo) GetReflectedType(ctx context.Context, typ string) (reflect.Type, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesContextKey).(map[string]reflect.Type); ok {
		return ti.getReflectedType(typ, customTypes)
	} else {
		return ti.getReflectedType(typ, nil)
	}
}

func (ti *typeInfo) getReflectedType(typ string, customTypes map[string]reflect.Type) (reflect.Type, error) {

	if ti.IsNullable(typ) {
		rt, err := ti.getReflectedType(ti.GetUnderlyingType(typ), customTypes)
		if err != nil {
			return nil, err
		}
		if utils.CanBeNil(rt) {
			return rt, nil
		} else {
			return reflect.PointerTo(rt), nil
		}
	}

	if customTypes != nil {
		if rt, ok := customTypes[typ]; ok {
			return rt, nil
		}
	}

	if rt, ok := reflectedTypeMap[typ]; ok {
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
	return rtMapStringAny, nil
}

var rtMapStringAny = reflect.TypeFor[map[string]any]()
var reflectedTypeMap = map[string]reflect.Type{
	"bool":                              reflect.TypeFor[bool](),
	"usize":                             reflect.TypeFor[uint](),
	"u8":                                reflect.TypeFor[uint8](),
	"u16":                               reflect.TypeFor[uint16](),
	"u32":                               reflect.TypeFor[uint32](),
	"u64":                               reflect.TypeFor[uint64](),
	"isize":                             reflect.TypeFor[int](),
	"i8":                                reflect.TypeFor[int8](),
	"i16":                               reflect.TypeFor[int16](),
	"i32":                               reflect.TypeFor[int32](),
	"i64":                               reflect.TypeFor[int64](),
	"f32":                               reflect.TypeFor[float32](),
	"f64":                               reflect.TypeFor[float64](),
	"~lib/string/String":                reflect.TypeFor[string](),
	"~lib/arraybuffer/ArrayBuffer":      reflect.TypeFor[[]byte](),
	"~lib/typedarray/Uint8ClampedArray": reflect.TypeFor[[]uint8](),
	"~lib/typedarray/Uint8Array":        reflect.TypeFor[[]uint8](),
	"~lib/typedarray/Uint16Array":       reflect.TypeFor[[]uint16](),
	"~lib/typedarray/Uint32Array":       reflect.TypeFor[[]uint32](),
	"~lib/typedarray/Uint64Array":       reflect.TypeFor[[]uint64](),
	"~lib/typedarray/Int8Array":         reflect.TypeFor[[]int8](),
	"~lib/typedarray/Int16Array":        reflect.TypeFor[[]int16](),
	"~lib/typedarray/Int32Array":        reflect.TypeFor[[]int32](),
	"~lib/typedarray/Int64Array":        reflect.TypeFor[[]int64](),
	"~lib/typedarray/Float32Array":      reflect.TypeFor[[]float32](),
	"~lib/typedarray/Float64Array":      reflect.TypeFor[[]float64](),
	"~lib/date/Date":                    reflect.TypeFor[time.Time](),
	"~lib/wasi_date/wasi_Date":          reflect.TypeFor[time.Time](),
}

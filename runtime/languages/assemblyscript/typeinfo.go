/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"
)

var _langTypeInfo = &langTypeInfo{}

func LanguageTypeInfo() langsupport.LanguageTypeInfo {
	return _langTypeInfo
}

func GetTypeInfo(ctx context.Context, typeName string, typeCache map[string]langsupport.TypeInfo) (langsupport.TypeInfo, error) {
	return langsupport.GetTypeInfo(ctx, _langTypeInfo, typeName, typeCache)
}

type langTypeInfo struct{}

func (lti *langTypeInfo) GetListSubtype(typ string) string {
	typ = lti.GetUnderlyingType(typ)
	if strings.HasPrefix(typ, "~lib/array/Array<") {
		return typ[17 : len(typ)-1]
	}

	return ""
}

func (lti *langTypeInfo) GetMapSubtypes(typ string) (string, string) {
	typ = lti.GetUnderlyingType(typ)
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

func (lti *langTypeInfo) GetNameForType(typ string) string {
	s := lti.GetUnderlyingType(typ)

	if lti.IsListType(s) {
		return lti.GetNameForType(lti.GetListSubtype(s)) + "[]"
	}

	if lti.IsMapType(s) {
		kt, vt := lti.GetMapSubtypes(s)
		return "Map<" + lti.GetNameForType(kt) + "," + lti.GetNameForType(vt) + ">"
	}

	return s[strings.LastIndex(s, "/")+1:]
}

func (lti *langTypeInfo) GetUnderlyingType(typ string) string {
	if s, ok := strings.CutSuffix(typ, " | null"); ok {
		return s
	} else {
		return strings.TrimSuffix(typ, "|null")
	}
}

func (lti *langTypeInfo) IsNullableType(typ string) bool {
	return strings.HasSuffix(typ, " | null") || strings.HasSuffix(typ, "|null")
}

func (lti *langTypeInfo) IsListType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/array/Array<")
}

func (lti *langTypeInfo) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (lti *langTypeInfo) IsByteSequenceType(typ string) bool {
	typ = lti.GetUnderlyingType(typ)
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

func (lti *langTypeInfo) IsFloatType(typ string) bool {
	switch typ {
	case "f32", "f64":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64",
		"u8", "u16", "u32", "u64",
		"isize", "usize":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsMapType(typ string) bool {
	typ = lti.GetUnderlyingType(typ)
	return strings.HasPrefix(typ, "~lib/map/Map<")
}

func (lti *langTypeInfo) IsObjectType(typ string) bool {
	return !lti.IsPrimitiveType(typ) &&
		!lti.IsListType(typ) &&
		!lti.IsMapType(typ) &&
		!lti.IsStringType(typ) &&
		!lti.IsTimestampType(typ) &&
		!lti.IsArrayBufferType(typ) &&
		!lti.IsTypedArrayType(typ)
}

func (lti *langTypeInfo) IsPointerType(typ string) bool {
	return false // assemblyscript does not have pointer types
}

func (lti *langTypeInfo) IsPrimitiveType(typ string) bool {
	return lti.IsBooleanType(typ) || lti.IsIntegerType(typ) || lti.IsFloatType(typ)
}

func (lti *langTypeInfo) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "i8", "i16", "i32", "i64", "isize":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsStringType(typ string) bool {
	return lti.GetUnderlyingType(typ) == "~lib/string/String"
}

func (lti *langTypeInfo) IsArrayBufferType(typ string) bool {
	return lti.GetUnderlyingType(typ) == "~lib/arraybuffer/ArrayBuffer"
}

func (lti *langTypeInfo) IsTypedArrayType(typ string) bool {
	return strings.HasPrefix(typ, "~lib/typedarray/")
}

func (lti *langTypeInfo) IsTimestampType(typ string) bool {
	typ = lti.GetUnderlyingType(typ)
	switch typ {
	case "Date", "~lib/date/Date", "~lib/wasi_date/wasi_Date":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
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

func (lti *langTypeInfo) GetAlignmentOfType(ctx context.Context, typ string) (uint32, error) {
	return lti.GetSizeOfType(ctx, typ)
}

func (lti *langTypeInfo) ObjectsUseMaxFieldAlignment() bool {
	// AssemblyScript classes are aligned to the pointer size (4 bytes), not the max field alignment.
	return false
}

func (lti *langTypeInfo) GetDataSizeOfType(ctx context.Context, typ string) (uint32, error) {
	switch typ {
	case "u64", "i64", "f64":
		return 8, nil
	case "u32", "i32", "f32", "isize", "usize":
		return 4, nil
	case "u16", "i16":
		return 2, nil
	case "u8", "i8", "bool":
		return 1, nil
	}

	if lti.IsListType(typ) {
		return 16, nil
	} else if lti.IsMapType(typ) {
		return 24, nil
	} else if lti.IsTypedArrayType(typ) {
		return 12, nil
	} else if lti.IsTimestampType(typ) {
		return 20, nil
	}

	// calculate dynamic size later
	return 0, nil
}

func (lti *langTypeInfo) GetEncodingLengthOfType(ctx context.Context, typ string) (uint32, error) {
	return 1, nil
}

func (lti *langTypeInfo) GetReflectedType(ctx context.Context, typ string) (reflect.Type, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesContextKey).(map[string]reflect.Type); ok {
		return lti.getReflectedType(typ, customTypes)
	} else {
		return lti.getReflectedType(typ, nil)
	}
}

func (lti *langTypeInfo) getReflectedType(typ string, customTypes map[string]reflect.Type) (reflect.Type, error) {

	if lti.IsNullableType(typ) {
		rt, err := lti.getReflectedType(lti.GetUnderlyingType(typ), customTypes)
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

	if lti.IsListType(typ) {
		et := lti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid array type: %s", typ)
		}

		elementType, err := lti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if lti.IsMapType(typ) {
		kt, vt := lti.GetMapSubtypes(typ)
		if kt == "" || vt == "" {
			return nil, fmt.Errorf("invalid map type: %s", typ)
		}

		keyType, err := lti.getReflectedType(kt, customTypes)
		if err != nil {
			return nil, err
		}
		valType, err := lti.getReflectedType(vt, customTypes)
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

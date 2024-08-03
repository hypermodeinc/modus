/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"fmt"
	"reflect"
	"time"
)

func getGoType(asType string) (reflect.Type, error) {
	if goType, ok := asToGoTypeMap[asType]; ok {
		return goType, nil
	}

	if isArrayType(asType) {
		et := getArraySubtypeInfo(asType)
		if et.Path == "" {
			return nil, fmt.Errorf("invalid array type: %s", asType)
		}

		elementType, err := getGoType(et.Path)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if isMapType(asType) {
		kt, vt := getMapSubtypeInfo(asType)
		if kt.Path == "" || vt.Path == "" {
			return nil, fmt.Errorf("invalid map type: %s", asType)
		}

		keyType, err := getGoType(kt.Path)
		if err != nil {
			return nil, err
		}
		valType, err := getGoType(vt.Path)
		if err != nil {
			return nil, err
		}

		return reflect.MapOf(keyType, valType), nil
	}

	return nil, fmt.Errorf("unsupported AssemblyScript type: %s", asType)
}

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

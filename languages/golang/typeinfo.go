/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

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
	return typ[strings.Index(typ, "]")+1:]
}

func (ti *typeInfo) GetMapSubtypes(typ string) (string, string) {
	const prefix = "map["
	if !strings.HasPrefix(typ, prefix) {
		return "", ""
	}

	n := 1
	for i := len(prefix); i < len(typ); i++ {
		switch typ[i] {
		case '[':
			n++
		case ']':
			n--
			if n == 0 {
				return typ[len(prefix):i], typ[i+1:]
			}
		}
	}

	return "", ""
}

func (ti *typeInfo) GetNameForType(typ string) string {
	// "github.com/hypermodeAI/functions-go/examples/simple.Person" -> "Person"

	if ti.IsPointerType(typ) {
		return "*" + ti.GetNameForType(ti.GetUnderlyingType(typ))
	}

	if ti.IsListType(typ) {
		return "[]" + ti.GetNameForType(ti.GetListSubtype(typ))
	}

	if ti.IsMapType(typ) {
		kt, vt := ti.GetMapSubtypes(typ)
		return "map[" + ti.GetNameForType(kt) + "]" + ti.GetNameForType(vt)
	}

	return typ[strings.LastIndex(typ, ".")+1:]
}

func (ti *typeInfo) GetUnderlyingType(typ string) string {
	return strings.TrimPrefix(typ, "*")
}

func (ti *typeInfo) IsListType(typ string) bool {
	// covers both slices and arrays
	return len(typ) > 2 && typ[0] == '['
}

func (ti *typeInfo) IsSliceType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] == ']'
}

func (ti *typeInfo) IsArrayType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] != ']'
}

func (ti *typeInfo) ArrayLength(typ string) (int, error) {
	i := strings.Index(typ, "]")
	if i == -1 {
		return -1, fmt.Errorf("invalid array type: %s", typ)
	}

	size := typ[1:i]
	if size == "" {
		return -1, fmt.Errorf("invalid array type: %s", typ)
	}

	return strconv.Atoi(size)
}

func (ti *typeInfo) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (ti *typeInfo) IsByteSequenceType(typ string) bool {
	switch typ {
	case "[]byte", "[]uint8":
		return true
	}

	if ti.IsArrayType(typ) {
		switch ti.GetListSubtype(typ) {
		case "byte", "uint8":
			return true
		}
	}

	return false
}

func (ti *typeInfo) IsFloatType(typ string) bool {
	switch typ {
	case "float32", "float64":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "byte", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsMapType(typ string) bool {
	return strings.HasPrefix(typ, "map[")
}

func (ti *typeInfo) IsNullable(typ string) bool {
	return ti.IsPointerType(typ) || ti.IsSliceType(typ) || ti.IsMapType(typ)
}

func (ti *typeInfo) IsPointerType(typ string) bool {
	return strings.HasPrefix(typ, "*")
}

func (ti *typeInfo) IsPrimitiveType(typ string) bool {
	return ti.IsBooleanType(typ) || ti.IsIntegerType(typ) || ti.IsFloatType(typ)
}

func (ti *typeInfo) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (ti *typeInfo) IsStringType(typ string) bool {
	return typ == "string"
}

func (ti *typeInfo) IsTimestampType(typ string) bool {
	return typ == "time.Time"
}

func (ti *typeInfo) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
	switch typ {
	case "int8", "uint8", "bool", "byte":
		return 1, nil
	case "int16", "uint16":
		return 2, nil
	case "int32", "uint32", "float32", "rune",
		"int", "uint", "uintptr", "unsafe.Pointer": // we only support 32-bit wasm
		return 4, nil
	case "int64", "uint64", "float64", "time.Duration":
		return 8, nil
	case "string":
		// string header is a 4 byte pointer and 4 byte length
		return 8, nil
	}

	if ti.IsPointerType(typ) {
		return 4, nil
	}

	if ti.IsSliceType(typ) {
		// slice header is a 4 byte pointer, 4 byte length, 4 byte capacity
		return 12, nil
	}

	if ti.IsArrayType(typ) {
		arrSize, err := ti.ArrayLength(typ)
		if err != nil {
			return 0, err
		}

		t := ti.GetListSubtype(typ)
		elementSize, err := ti.GetSizeOfType(ctx, t)
		if err != nil {
			return 0, err
		}

		return uint32(arrSize) * elementSize, nil
	}

	if ti.IsMapType(typ) {
		// maps are passed by reference using a 4 byte pointer
		return 4, nil
	}

	if ti.IsTimestampType(typ) {
		// time.Time has 3 fields: 8 byte uint64, 8 byte int64, 4 byte pointer
		return 20, nil
	}

	return ti.getSizeOfStruct(ctx, typ)
}

func (ti *typeInfo) getSizeOfStruct(ctx context.Context, typ string) (uint32, error) {
	def, err := ti.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, err
	}

	offset := uint32(0)
	for _, field := range def.Fields {
		size, err := ti.GetSizeOfType(ctx, field.Type)
		if err != nil {
			return 0, err
		}

		pad := langsupport.GetAlignmentPadding(offset, size)
		offset += size + pad
	}

	return offset, nil
}

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
	if customTypes != nil {
		if rt, ok := customTypes[typ]; ok {
			return rt, nil
		}
	}

	if rt, ok := reflectedTypeMap[typ]; ok {
		return rt, nil
	}

	if ti.IsPointerType(typ) {
		tt := ti.GetUnderlyingType(typ)
		targetType, err := ti.getReflectedType(tt, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.PointerTo(targetType), nil
	}

	if ti.IsSliceType(typ) {
		et := ti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid slice type: %s", typ)
		}

		elementType, err := ti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if ti.IsArrayType(typ) {
		et := ti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid array type: %s", typ)
		}

		size, err := ti.ArrayLength(typ)
		if err != nil {
			return nil, err
		}

		elementType, err := ti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.ArrayOf(size, elementType), nil
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
	"bool":           reflect.TypeFor[bool](),
	"byte":           reflect.TypeFor[byte](),
	"uint":           reflect.TypeFor[uint](),
	"uint8":          reflect.TypeFor[uint8](),
	"uint16":         reflect.TypeFor[uint16](),
	"uint32":         reflect.TypeFor[uint32](),
	"uint64":         reflect.TypeFor[uint64](),
	"uintptr":        reflect.TypeFor[uintptr](),
	"int":            reflect.TypeFor[int](),
	"int8":           reflect.TypeFor[int8](),
	"int16":          reflect.TypeFor[int16](),
	"int32":          reflect.TypeFor[int32](),
	"int64":          reflect.TypeFor[int64](),
	"rune":           reflect.TypeFor[rune](),
	"float32":        reflect.TypeFor[float32](),
	"float64":        reflect.TypeFor[float64](),
	"string":         reflect.TypeFor[string](),
	"time.Time":      reflect.TypeFor[time.Time](),
	"time.Duration":  reflect.TypeFor[time.Duration](),
	"unsafe.Pointer": reflect.TypeFor[unsafe.Pointer](),
}

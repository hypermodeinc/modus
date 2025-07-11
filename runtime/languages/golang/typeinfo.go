/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package golang

import (
	"context"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/hypermodeinc/modus/lib/metadata"
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
	return typ[strings.Index(typ, "]")+1:]
}

func (lti *langTypeInfo) GetMapSubtypes(typ string) (string, string) {
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

func (lti *langTypeInfo) GetNameForType(typ string) string {
	// "github.com/hypermodeinc/modus/sdk/go/examples/simple.Person" -> "Person"

	if lti.IsPointerType(typ) {
		return "*" + lti.GetNameForType(lti.GetUnderlyingType(typ))
	}

	if lti.IsListType(typ) {
		return "[]" + lti.GetNameForType(lti.GetListSubtype(typ))
	}

	if lti.IsMapType(typ) {
		kt, vt := lti.GetMapSubtypes(typ)
		return "map[" + lti.GetNameForType(kt) + "]" + lti.GetNameForType(vt)
	}

	return typ[strings.LastIndex(typ, ".")+1:]
}

func (lti *langTypeInfo) IsObjectType(typ string) bool {
	return !lti.IsPrimitiveType(typ) &&
		!lti.IsListType(typ) &&
		!lti.IsMapType(typ) &&
		!lti.IsStringType(typ) &&
		!lti.IsTimestampType(typ) &&
		!lti.IsPointerType(typ)
}

func (lti *langTypeInfo) GetUnderlyingType(typ string) string {
	return strings.TrimPrefix(typ, "*")
}

func (lti *langTypeInfo) IsListType(typ string) bool {
	// covers both slices and arrays
	return len(typ) > 2 && typ[0] == '['
}

func (lti *langTypeInfo) IsSliceType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] == ']'
}

func (lti *langTypeInfo) IsArrayType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] != ']'
}

func (lti *langTypeInfo) ArrayLength(typ string) (int, error) {
	i := strings.Index(typ, "]")
	if i == -1 {
		return -1, fmt.Errorf("invalid array type: %s", typ)
	}

	size := typ[1:i]
	if size == "" {
		return -1, fmt.Errorf("invalid array type: %s", typ)
	}

	parsedSize, err := strconv.Atoi(size)
	if err != nil {
		return -1, err
	}
	if parsedSize < 0 || parsedSize > math.MaxUint32 {
		return -1, fmt.Errorf("array size out of bounds: %s", size)
	}

	return parsedSize, nil
}

func (lti *langTypeInfo) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (lti *langTypeInfo) IsByteSequenceType(typ string) bool {
	switch typ {
	case "[]byte", "[]uint8":
		return true
	}

	if lti.IsArrayType(typ) {
		switch lti.GetListSubtype(typ) {
		case "byte", "uint8":
			return true
		}
	}

	return false
}

func (lti *langTypeInfo) IsFloatType(typ string) bool {
	switch typ {
	case "float32", "float64":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "byte", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsMapType(typ string) bool {
	return strings.HasPrefix(typ, "map[")
}

func (lti *langTypeInfo) IsNullableType(typ string) bool {
	return lti.IsPointerType(typ) || lti.IsSliceType(typ) || lti.IsMapType(typ)
}

func (lti *langTypeInfo) IsPointerType(typ string) bool {
	return strings.HasPrefix(typ, "*")
}

func (lti *langTypeInfo) IsPrimitiveType(typ string) bool {
	return lti.IsBooleanType(typ) || lti.IsIntegerType(typ) || lti.IsFloatType(typ)
}

func (lti *langTypeInfo) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (lti *langTypeInfo) IsStringType(typ string) bool {
	return typ == "string"
}

func (lti *langTypeInfo) IsTimestampType(typ string) bool {
	return typ == "time.Time"
}

func (lti *langTypeInfo) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
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
	}

	if lti.IsStringType(typ) {
		// string header is a 4 byte pointer and 4 byte length
		return 8, nil
	}

	if lti.IsPointerType(typ) {
		return 4, nil
	}

	if lti.IsMapType(typ) {
		// maps are passed by reference using a 4 byte pointer
		return 4, nil
	}

	if lti.IsSliceType(typ) {
		// slice header is a 4 byte pointer, 4 byte length, 4 byte capacity
		return 12, nil
	}

	if lti.IsTimestampType(typ) {
		// time.Time has 3 fields: 8 byte uint64, 8 byte int64, 4 byte pointer
		return 20, nil
	}

	if lti.IsArrayType(typ) {
		return lti.getSizeOfArray(ctx, typ)
	}

	return lti.getSizeOfStruct(ctx, typ)
}

func (lti *langTypeInfo) getSizeOfArray(ctx context.Context, typ string) (uint32, error) {
	// array size is the element size times the number of elements, aligned to the element size
	arrSize, err := lti.ArrayLength(typ)
	if err != nil {
		return 0, err
	}
	if arrSize == 0 {
		return 0, nil
	}

	t := lti.GetListSubtype(typ)
	elementAlignment, err := lti.GetAlignmentOfType(ctx, t)
	if err != nil {
		return 0, err
	}
	elementSize, err := lti.GetSizeOfType(ctx, t)
	if err != nil {
		return 0, err
	}

	size := langsupport.AlignOffset(elementSize, elementAlignment)*uint32(arrSize-1) + elementSize
	return size, nil
}

func (lti *langTypeInfo) getSizeOfStruct(ctx context.Context, typ string) (uint32, error) {
	def, err := lti.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, err
	}
	if len(def.Fields) == 0 {
		return 0, nil
	}

	offset := uint32(0)
	maxAlign := uint32(1)
	for _, field := range def.Fields {
		size, err := lti.GetSizeOfType(ctx, field.Type)
		if err != nil {
			return 0, err
		}
		alignment, err := lti.GetAlignmentOfType(ctx, field.Type)
		if err != nil {
			return 0, err
		}
		if alignment > maxAlign {
			maxAlign = alignment
		}
		offset = langsupport.AlignOffset(offset, alignment)
		offset += size
	}

	size := langsupport.AlignOffset(offset, maxAlign)
	return size, nil
}

func (lti *langTypeInfo) GetAlignmentOfType(ctx context.Context, typ string) (uint32, error) {

	// reference: https://github.com/tinygo-org/tinygo/blob/release/compiler/sizes.go

	// primitives align to their natural size
	if lti.IsPrimitiveType(typ) {
		return lti.GetSizeOfType(ctx, typ)
	}

	// arrays align to the alignment of their element type
	if lti.IsArrayType(typ) {
		t := lti.GetListSubtype(typ)
		return lti.GetAlignmentOfType(ctx, t)
	}

	// reference types align to the pointer size (4 bytes on 32-bit wasm)
	if lti.IsPointerType(typ) || lti.IsSliceType(typ) || lti.IsStringType(typ) || lti.IsMapType(typ) {
		return 4, nil
	}

	// time.Time has 3 fields, the maximum alignment is 8 bytes
	if lti.IsTimestampType(typ) {
		return 8, nil
	}

	// structs align to the maximum alignment of their fields
	return lti.getAlignmentOfStruct(ctx, typ)
}

func (lti *langTypeInfo) ObjectsUseMaxFieldAlignment() bool {
	// Go structs are aligned to the maximum alignment of their fields
	return true
}

func (lti *langTypeInfo) getAlignmentOfStruct(ctx context.Context, typ string) (uint32, error) {
	def, err := lti.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, err
	}

	max := uint32(1)
	for _, field := range def.Fields {
		align, err := lti.GetAlignmentOfType(ctx, field.Type)
		if err != nil {
			return 0, err
		}
		if align > max {
			max = align
		}
	}

	return max, nil
}

func (lti *langTypeInfo) GetDataSizeOfType(ctx context.Context, typ string) (uint32, error) {
	return lti.GetSizeOfType(ctx, typ)
}

func (lti *langTypeInfo) GetEncodingLengthOfType(ctx context.Context, typ string) (uint32, error) {
	if lti.IsPrimitiveType(typ) || lti.IsPointerType(typ) || lti.IsMapType(typ) {
		return 1, nil
	} else if lti.IsStringType(typ) {
		return 2, nil
	} else if lti.IsSliceType(typ) || lti.IsTimestampType(typ) {
		return 3, nil
	} else if lti.IsArrayType(typ) {
		return lti.getEncodingLengthOfArray(ctx, typ)
	} else if lti.IsObjectType(typ) {
		return lti.getEncodingLengthOfStruct(ctx, typ)
	}

	return 0, fmt.Errorf("unable to determine encoding length for type: %s", typ)
}

func (lti *langTypeInfo) getEncodingLengthOfArray(ctx context.Context, typ string) (uint32, error) {
	arrSize, err := lti.ArrayLength(typ)
	if err != nil {
		return 0, err
	}
	if arrSize == 0 {
		return 0, nil
	}

	t := lti.GetListSubtype(typ)
	elementLen, err := lti.GetEncodingLengthOfType(ctx, t)
	if err != nil {
		return 0, err
	}

	return uint32(arrSize) * elementLen, nil
}

func (lti *langTypeInfo) getEncodingLengthOfStruct(ctx context.Context, typ string) (uint32, error) {
	def, err := lti.GetTypeDefinition(ctx, typ)
	if err != nil {
		return 0, err
	}

	total := uint32(0)
	for _, field := range def.Fields {
		len, err := lti.GetEncodingLengthOfType(ctx, field.Type)
		if err != nil {
			return 0, err
		}
		total += len
	}

	return total, nil
}

func (lti *langTypeInfo) GetTypeDefinition(ctx context.Context, typ string) (*metadata.TypeDefinition, error) {
	md := ctx.Value(utils.MetadataContextKey).(*metadata.Metadata)
	return md.GetTypeDefinition(typ)
}

func (lti *langTypeInfo) GetReflectedType(ctx context.Context, typ string) (reflect.Type, error) {
	if customTypes, ok := ctx.Value(utils.CustomTypesContextKey).(map[string]reflect.Type); ok {
		return lti.getReflectedType(typ, customTypes)
	} else {
		return lti.getReflectedType(typ, nil)
	}
}

func (lti *langTypeInfo) getReflectedType(typ string, customTypes map[string]reflect.Type) (reflect.Type, error) {
	if customTypes != nil {
		if rt, ok := customTypes[typ]; ok {
			return rt, nil
		}
	}

	if rt, ok := reflectedTypeMap[typ]; ok {
		return rt, nil
	}

	if lti.IsPointerType(typ) {
		tt := lti.GetUnderlyingType(typ)
		targetType, err := lti.getReflectedType(tt, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.PointerTo(targetType), nil
	}

	if lti.IsSliceType(typ) {
		et := lti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid slice type: %s", typ)
		}

		elementType, err := lti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.SliceOf(elementType), nil
	}

	if lti.IsArrayType(typ) {
		et := lti.GetListSubtype(typ)
		if et == "" {
			return nil, fmt.Errorf("invalid array type: %s", typ)
		}

		size, err := lti.ArrayLength(typ)
		if err != nil {
			return nil, err
		}

		elementType, err := lti.getReflectedType(et, customTypes)
		if err != nil {
			return nil, err
		}
		return reflect.ArrayOf(size, elementType), nil
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

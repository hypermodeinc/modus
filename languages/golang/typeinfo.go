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

	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

func (ti *typeInfoProvider) GetListSubtype(typ string) string {
	return typ[strings.Index(typ, "]")+1:]
}

func (ti *typeInfoProvider) GetMapSubtypes(typ string) (string, string) {
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

func (ti *typeInfoProvider) GetNameForType(typ string) string {
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

func (ti *typeInfoProvider) GetUnderlyingType(typ string) string {
	return strings.TrimPrefix(typ, "*")
}

func (ti *typeInfoProvider) IsListType(typ string) bool {
	// covers both slices and arrays
	return len(typ) > 2 && typ[0] == '['
}

func (ti *typeInfoProvider) IsSliceType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] == ']'
}

func (ti *typeInfoProvider) IsArrayType(typ string) bool {
	return len(typ) > 2 && typ[0] == '[' && typ[1] != ']'
}

func (ti *typeInfoProvider) ArraySize(typ string) (int, error) {
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

func (ti *typeInfoProvider) IsBooleanType(typ string) bool {
	return typ == "bool"
}

func (ti *typeInfoProvider) IsByteSequenceType(typ string) bool {
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

func (ti *typeInfoProvider) IsFloatType(typ string) bool {
	switch typ {
	case "float32", "float64":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64",
		"uintptr", "byte", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsMapType(typ string) bool {
	return strings.HasPrefix(typ, "map[")
}

func (ti *typeInfoProvider) IsNullable(typ string) bool {
	return ti.IsPointerType(typ) || ti.IsSliceType(typ) || ti.IsMapType(typ)
}

func (ti *typeInfoProvider) IsPointerType(typ string) bool {
	return strings.HasPrefix(typ, "*")
}

func (ti *typeInfoProvider) IsPrimitiveType(typ string) bool {
	return ti.IsBooleanType(typ) || ti.IsIntegerType(typ) || ti.IsFloatType(typ)
}

func (ti *typeInfoProvider) IsSignedIntegerType(typ string) bool {
	switch typ {
	case "int", "int8", "int16", "int32", "int64", "rune", "time.Duration":
		return true
	default:
		return false
	}
}

func (ti *typeInfoProvider) IsStringType(typ string) bool {
	return typ == "string"
}

func (ti *typeInfoProvider) IsTimestampType(typ string) bool {
	return typ == "time.Time"
}

func (ti *typeInfoProvider) GetSizeOfType(ctx context.Context, typ string) (uint32, error) {
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
		arrSize, err := ti.ArraySize(typ)
		if err != nil {
			return 0, err
		}

		t := ti.GetListSubtype(typ)
		itemSize, err := ti.GetSizeOfType(ctx, t)
		if err != nil {
			return 0, err
		}

		return uint32(arrSize) * itemSize, nil
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

func (ti *typeInfoProvider) getSizeOfStruct(ctx context.Context, typ string) (uint32, error) {
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

		pad := ti.getAlignmentPadding(offset, size)
		offset += size + pad
	}

	return offset, nil
}

func (ti *typeInfoProvider) getAlignmentPadding(offset, size uint32) uint32 {

	// maximum alignment is 4 bytes on 32-bit wasm
	if size > 4 {
		size = 4
	}

	mask := size - 1
	if offset&mask != 0 {
		return (offset | mask) + 1 - offset
	}

	return 0
}

func (ti *typeInfoProvider) GetTypeDefinition(ctx context.Context, typ string) (*metadata.TypeDefinition, error) {

	md := ctx.Value(utils.MetadataContextKey).(*metadata.Metadata)
	def, ok := md.Types[typ]
	if !ok {
		return nil, fmt.Errorf("info for type %s not found in plugin %s", typ, md.Name())
	}

	return def, nil
}

func (ti *typeInfoProvider) getReflectedType(typ string, customTypes map[string]reflect.Type) (reflect.Type, error) {
	if customTypes != nil {
		if rt, ok := customTypes[typ]; ok {
			return rt, nil
		}
	}

	if rt, ok := goTypeMap[typ]; ok {
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

		size, err := ti.ArraySize(typ)
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
	return goMapStringAnyType, nil
}

var goMapStringAnyType = reflect.TypeOf(map[string]any{})

// map of all types that can be converted to Go types
var goTypeMap = map[string]reflect.Type{
	"bool":          reflect.TypeOf(bool(false)),
	"byte":          reflect.TypeOf(byte(0)),
	"uint":          reflect.TypeOf(uint(0)),
	"uint8":         reflect.TypeOf(uint8(0)),
	"uint16":        reflect.TypeOf(uint16(0)),
	"uint32":        reflect.TypeOf(uint32(0)),
	"uint64":        reflect.TypeOf(uint64(0)),
	"uintptr":       reflect.TypeOf(uintptr(0)),
	"int":           reflect.TypeOf(int(0)),
	"int8":          reflect.TypeOf(int8(0)),
	"int16":         reflect.TypeOf(int16(0)),
	"int32":         reflect.TypeOf(int32(0)),
	"int64":         reflect.TypeOf(int64(0)),
	"rune":          reflect.TypeOf(rune(0)),
	"float32":       reflect.TypeOf(float32(0)),
	"float64":       reflect.TypeOf(float64(0)),
	"string":        reflect.TypeOf(string("")),
	"time.Time":     reflect.TypeOf(time.Time{}),
	"time.Duration": reflect.TypeOf(time.Duration(0)),
}

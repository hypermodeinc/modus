/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"context"
	"errors"
	"reflect"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/utils"
)

type TypeInfo interface {
	Name() string
	ReflectedType() reflect.Type
	ZeroValue() any
	Size() uint32
	Alignment() uint32
	DataSize() uint32
	EncodingLength() uint32

	UnderlyingType() TypeInfo
	ListElementType() TypeInfo
	MapKeyType() TypeInfo
	MapValueType() TypeInfo
	ObjectFieldTypes() []TypeInfo
	ObjectFieldOffsets() []uint32

	IsBoolean() bool
	IsByteSequence() bool
	IsFloat() bool
	IsInteger() bool
	IsList() bool
	IsMap() bool
	IsNullable() bool
	IsObject() bool
	IsPointer() bool
	IsPrimitive() bool
	IsSignedInteger() bool
	IsString() bool
	IsTimestamp() bool
}

func GetTypeInfo(ctx context.Context, lti LanguageTypeInfo, typeName string, typeCache map[string]TypeInfo) (TypeInfo, error) {
	if t, ok := typeCache[typeName]; ok {
		return t, nil
	}

	info := &typeInfo{name: typeName}
	typeCache[typeName] = info

	flags := typeFlags(0)

	if lti.IsPrimitiveType(typeName) {
		flags |= tfPrimitive
		if lti.IsBooleanType(typeName) {
			flags |= tfBoolean
		} else if lti.IsFloatType(typeName) {
			flags |= tfFloat
		} else if lti.IsIntegerType(typeName) {
			flags |= tfInteger
			if lti.IsSignedIntegerType(typeName) {
				flags |= tfSignedInteger
			}
		}
	}

	if lti.IsPointerType(typeName) {
		flags |= tfPointer
		flags |= tfNullable

		underlyingTypeName := lti.GetUnderlyingType(typeName)
		if underlyingTypeName != typeName {
			underlyingTypeInfo, err := GetTypeInfo(ctx, lti, underlyingTypeName, typeCache)
			if err != nil {
				return nil, err
			}
			info.underlyingType = underlyingTypeInfo
		}
	} else if lti.IsNullableType(typeName) {
		flags |= tfNullable

		underlyingTypeName := lti.GetUnderlyingType(typeName)
		if underlyingTypeName != typeName {
			underlyingTypeInfo, err := GetTypeInfo(ctx, lti, underlyingTypeName, typeCache)
			if err != nil {
				return nil, err
			}
			info.underlyingType = underlyingTypeInfo

			// remaining tests are for the underlying (non-null) type
			typeName = underlyingTypeName
		}
	}

	if lti.IsByteSequenceType(typeName) {
		flags |= tfByteSequence
	} else if lti.IsStringType(typeName) {
		flags |= tfString
	} else if lti.IsTimestampType(typeName) {
		flags |= tfTimestamp
	}

	if lti.IsListType(typeName) {
		flags |= tfList
		elementTypeName := lti.GetListSubtype(typeName)
		elementTypeInfo, err := GetTypeInfo(ctx, lti, elementTypeName, typeCache)
		if err != nil {
			return nil, err
		}
		info.fieldTypes = []TypeInfo{elementTypeInfo}
	} else if lti.IsMapType(typeName) {
		flags |= tfMap
		keyTypeName, keyValueName := lti.GetMapSubtypes(typeName)
		keyTypeInfo, err := GetTypeInfo(ctx, lti, keyTypeName, typeCache)
		if err != nil {
			return nil, err
		}
		valueTypeInfo, err := GetTypeInfo(ctx, lti, keyValueName, typeCache)
		if err != nil {
			return nil, err
		}
		info.fieldTypes = []TypeInfo{keyTypeInfo, valueTypeInfo}
	} else if lti.IsObjectType(typeName) {
		flags |= tfObject

		md, err := getMetadataFromContext(ctx)
		if err != nil {
			return nil, err
		}

		def, err := md.GetTypeDefinition(typeName)
		if err != nil {
			return nil, err
		}

		offset := uint32(0)
		maxAlignment := uint32(0)
		info.fieldTypes = make([]TypeInfo, len(def.Fields))
		info.fieldOffsets = make([]uint32, len(def.Fields))
		for i, field := range def.Fields {
			fti, err := GetTypeInfo(ctx, lti, field.Type, typeCache)
			if err != nil {
				return nil, err
			}
			info.fieldTypes[i] = fti

			size := fti.Size()
			alignment := fti.Alignment()
			offset = AlignOffset(offset, alignment)
			info.fieldOffsets[i] = offset
			offset += size

			if alignment > maxAlignment {
				maxAlignment = alignment
			}
		}
		info.dataSize = AlignOffset(offset, maxAlignment)

		if lti.ObjectsUseMaxFieldAlignment() {
			info.alignment = maxAlignment
		}
	}

	info.flags = flags

	reflectedType, err := lti.GetReflectedType(ctx, typeName)
	if err != nil {
		return nil, err
	}
	info.reflectedType = reflectedType
	info.zeroValue = reflect.Zero(reflectedType).Interface()

	if info.size == 0 {
		size, err := lti.GetSizeOfType(ctx, typeName)
		if err != nil {
			return nil, err
		}
		info.size = size
	}

	if info.alignment == 0 {
		alignment, err := lti.GetAlignmentOfType(ctx, typeName)
		if err != nil {
			return nil, err
		}
		info.alignment = alignment
	}

	if info.dataSize == 0 {
		dataSize, err := lti.GetDataSizeOfType(ctx, typeName)
		if err != nil {
			return nil, err
		}
		info.dataSize = dataSize
	}

	if info.encodingLength == 0 {
		encodingLength, err := lti.GetEncodingLengthOfType(ctx, typeName)
		if err != nil {
			return nil, err
		}
		info.encodingLength = encodingLength
	}

	return info, nil
}

type typeFlags uint32

const (
	_         typeFlags = 0
	tfBoolean typeFlags = 1 << (iota - 1)
	tfByteSequence
	tfFloat
	tfInteger
	tfList
	tfMap
	tfNullable
	tfObject
	tfPointer
	tfPrimitive
	tfSignedInteger
	tfString
	tfTimestamp
)

type typeInfo struct {
	name           string
	flags          typeFlags
	reflectedType  reflect.Type
	zeroValue      any
	size           uint32
	alignment      uint32
	dataSize       uint32
	encodingLength uint32
	underlyingType TypeInfo
	fieldTypes     []TypeInfo
	fieldOffsets   []uint32
}

func (h *typeInfo) Name() string                { return h.name }
func (h *typeInfo) ReflectedType() reflect.Type { return h.reflectedType }
func (h *typeInfo) ZeroValue() any              { return h.zeroValue }
func (h *typeInfo) Size() uint32                { return h.size }
func (h *typeInfo) Alignment() uint32           { return h.alignment }
func (h *typeInfo) DataSize() uint32            { return h.dataSize }
func (h *typeInfo) EncodingLength() uint32      { return h.encodingLength }

func (h *typeInfo) IsBoolean() bool       { return h.flags&tfBoolean != 0 }
func (h *typeInfo) IsByteSequence() bool  { return h.flags&tfByteSequence != 0 }
func (h *typeInfo) IsFloat() bool         { return h.flags&tfFloat != 0 }
func (h *typeInfo) IsInteger() bool       { return h.flags&tfInteger != 0 }
func (h *typeInfo) IsList() bool          { return h.flags&tfList != 0 }
func (h *typeInfo) IsMap() bool           { return h.flags&tfMap != 0 }
func (h *typeInfo) IsNullable() bool      { return h.flags&tfNullable != 0 }
func (h *typeInfo) IsObject() bool        { return h.flags&tfObject != 0 }
func (h *typeInfo) IsPointer() bool       { return h.flags&tfPointer != 0 }
func (h *typeInfo) IsPrimitive() bool     { return h.flags&tfPrimitive != 0 }
func (h *typeInfo) IsSignedInteger() bool { return h.flags&tfSignedInteger != 0 }
func (h *typeInfo) IsString() bool        { return h.flags&tfString != 0 }
func (h *typeInfo) IsTimestamp() bool     { return h.flags&tfTimestamp != 0 }

func (h *typeInfo) UnderlyingType() TypeInfo {
	return h.underlyingType
}

func (h *typeInfo) ListElementType() TypeInfo {
	if h.IsList() {
		return h.fieldTypes[0]
	}
	return nil
}

func (h *typeInfo) MapKeyType() TypeInfo {
	if h.IsMap() {
		return h.fieldTypes[0]
	}
	return nil
}

func (h *typeInfo) MapValueType() TypeInfo {
	if h.IsMap() {
		return h.fieldTypes[1]
	}
	return nil
}

func (h *typeInfo) ObjectFieldTypes() []TypeInfo {
	if h.IsObject() {
		return h.fieldTypes
	}
	return nil
}

func (h *typeInfo) ObjectFieldOffsets() []uint32 {
	return h.fieldOffsets
}

func getMetadataFromContext(ctx context.Context) (*metadata.Metadata, error) {
	v := ctx.Value(utils.MetadataContextKey)
	if v == nil {
		return nil, errors.New("metadata not found in context")
	}
	return v.(*metadata.Metadata), nil
}

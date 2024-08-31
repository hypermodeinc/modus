/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"hypruntime/langsupport"
	"hypruntime/langsupport/primitives"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"

	"golang.org/x/exp/constraints"
)

func (p *planner) NewTypedArrayHandler(typ string, rt reflect.Type) (managedTypeHandler, error) {
	info := langsupport.NewTypeHandlerInfo(typ, rt, 12, 1)
	typ = _typeInfo.GetUnderlyingType(typ)
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}

	switch typ {
	case "~lib/typedarray/Uint8Array", "~lib/typedarray/Uint8ClampedArray":
		converter := primitives.NewPrimitiveTypeConverter[uint8]()
		return &typedArrayHandler[uint8]{info, typeDef, converter}, nil
	case "~lib/typedarray/Uint16Array":
		converter := primitives.NewPrimitiveTypeConverter[uint16]()
		return &typedArrayHandler[uint16]{info, typeDef, converter}, nil
	case "~lib/typedarray/Uint32Array":
		converter := primitives.NewPrimitiveTypeConverter[uint32]()
		return &typedArrayHandler[uint32]{info, typeDef, converter}, nil
	case "~lib/typedarray/Uint64Array":
		converter := primitives.NewPrimitiveTypeConverter[uint64]()
		return &typedArrayHandler[uint64]{info, typeDef, converter}, nil
	case "~lib/typedarray/Int8Array":
		converter := primitives.NewPrimitiveTypeConverter[int8]()
		return &typedArrayHandler[int8]{info, typeDef, converter}, nil
	case "~lib/typedarray/Int16Array":
		converter := primitives.NewPrimitiveTypeConverter[int16]()
		return &typedArrayHandler[int16]{info, typeDef, converter}, nil
	case "~lib/typedarray/Int32Array":
		converter := primitives.NewPrimitiveTypeConverter[int32]()
		return &typedArrayHandler[int32]{info, typeDef, converter}, nil
	case "~lib/typedarray/Int64Array":
		converter := primitives.NewPrimitiveTypeConverter[int64]()
		return &typedArrayHandler[int64]{info, typeDef, converter}, nil
	case "~lib/typedarray/Float32Array":
		converter := primitives.NewPrimitiveTypeConverter[float32]()
		return &typedArrayHandler[float32]{info, typeDef, converter}, nil
	case "~lib/typedarray/Float64Array":
		converter := primitives.NewPrimitiveTypeConverter[float64]()
		return &typedArrayHandler[float64]{info, typeDef, converter}, nil
	default:
		return nil, fmt.Errorf("unsupported typed array type: %s", typ)
	}
}

type typedArrayHandler[T constraints.Integer | constraints.Float] struct {
	info      langsupport.TypeHandlerInfo
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
}

func (h *typedArrayHandler[T]) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *typedArrayHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	// The base offset is supposed to be the address to the ArrayBuffer that backs the typed array.
	// However, it appears to be identical to the data start address.
	// We don't need it anyway, because we only read the part that's in view.

	dataStart, ok := wa.Memory().ReadUint32Le(offset + 4)
	if !ok {
		return nil, fmt.Errorf("failed to read data start address for %s", h.info.TypeName())
	}

	byteLen, ok := wa.Memory().ReadUint32Le(offset + 8)
	if !ok {
		return nil, fmt.Errorf("failed to read byte length for %s", h.info.TypeName())
	} else if byteLen == 0 {
		// empty typed array
		return []T{}, nil
	}

	buf, ok := wa.Memory().Read(dataStart, byteLen)
	if !ok {
		return nil, errors.New("failed to read array data")
	}

	items := h.converter.BytesToSlice(buf)
	return items, nil
}

func (h *typedArrayHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	items, ok := obj.([]T)
	if !ok {
		return nil, fmt.Errorf("expected a %T, got %T", []T{}, obj)
	} else if len(items) == 0 {
		// empty typed array
		return nil, nil
	}

	bytes := h.converter.SliceToBytes(items)

	// allocate memory for the buffer
	bufferSize := uint32(len(bytes))
	ptr, cln, err := wa.AllocateMemory(ctx, bufferSize)
	if err != nil {
		return cln, err
	}

	// write the buffer
	if ok := wa.Memory().Write(ptr, bytes); !ok {
		return cln, fmt.Errorf("failed to write typed array data for %s", h.info.TypeName())
	}

	// write the typed array object
	if ok := wa.Memory().WriteUint32Le(offset, ptr); !ok {
		return cln, errors.New("failed to write typed array buffer pointer")
	}
	if ok := wa.Memory().WriteUint32Le(offset+4, ptr); !ok {
		return cln, errors.New("failed to write typed array data start pointer")
	}
	if ok := wa.Memory().WriteUint32Le(offset+8, bufferSize); !ok {
		return cln, errors.New("failed to write typed array byte length")
	}

	return nil, nil
}

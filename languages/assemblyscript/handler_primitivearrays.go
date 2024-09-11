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
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/array.ts

func (p *planner) NewPrimitiveArrayHandler(ctx context.Context, typ string, rt reflect.Type) (managedTypeHandler, error) {

	typ = _typeInfo.GetUnderlyingType(typ)
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}

	elementType := _typeInfo.GetListSubtype(typ)
	if elementType == "" {
		return nil, errors.New("array type must have a subtype")
	}

	info := langsupport.NewTypeHandlerInfo(typ, rt, 16, 0)

	switch elementType {
	case "bool":
		converter := primitives.NewPrimitiveTypeConverter[bool]()
		return &primitiveArrayHandler[bool]{info, typeDef, converter}, nil
	case "u8":
		converter := primitives.NewPrimitiveTypeConverter[uint8]()
		return &primitiveArrayHandler[uint8]{info, typeDef, converter}, nil
	case "u16":
		converter := primitives.NewPrimitiveTypeConverter[uint16]()
		return &primitiveArrayHandler[uint16]{info, typeDef, converter}, nil
	case "u32":
		converter := primitives.NewPrimitiveTypeConverter[uint32]()
		return &primitiveArrayHandler[uint32]{info, typeDef, converter}, nil
	case "u64":
		converter := primitives.NewPrimitiveTypeConverter[uint64]()
		return &primitiveArrayHandler[uint64]{info, typeDef, converter}, nil
	case "i8":
		converter := primitives.NewPrimitiveTypeConverter[int8]()
		return &primitiveArrayHandler[int8]{info, typeDef, converter}, nil
	case "i16":
		converter := primitives.NewPrimitiveTypeConverter[int16]()
		return &primitiveArrayHandler[int16]{info, typeDef, converter}, nil
	case "i32":
		converter := primitives.NewPrimitiveTypeConverter[int32]()
		return &primitiveArrayHandler[int32]{info, typeDef, converter}, nil
	case "i64":
		converter := primitives.NewPrimitiveTypeConverter[int64]()
		return &primitiveArrayHandler[int64]{info, typeDef, converter}, nil
	case "f32":
		converter := primitives.NewPrimitiveTypeConverter[float32]()
		return &primitiveArrayHandler[float32]{info, typeDef, converter}, nil
	case "f64":
		converter := primitives.NewPrimitiveTypeConverter[float64]()
		return &primitiveArrayHandler[float64]{info, typeDef, converter}, nil
	case "isize":
		converter := primitives.NewPrimitiveTypeConverter[int]()
		return &primitiveArrayHandler[int]{info, typeDef, converter}, nil
	case "usize":
		converter := primitives.NewPrimitiveTypeConverter[uint]()
		return &primitiveArrayHandler[uint]{info, typeDef, converter}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive array type: %s", typ)
	}
}

type primitiveArrayHandler[T primitive] struct {
	info      langsupport.TypeHandlerInfo
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
}

func (h *primitiveArrayHandler[T]) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *primitiveArrayHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	data, ok := wa.Memory().ReadUint32Le(offset + 4)
	if !ok {
		return nil, errors.New("failed to read array data pointer")
	}

	arrLen, ok := wa.Memory().ReadUint32Le(offset + 12)
	if !ok {
		return nil, errors.New("failed to read array length")
	} else if arrLen == 0 {
		// empty array
		return []T{}, nil
	}

	bufferSize := arrLen * uint32(h.converter.TypeSize())
	buf, ok := wa.Memory().Read(data, bufferSize)
	if !ok {
		return nil, errors.New("failed to read array data")
	}

	items := h.converter.BytesToSlice(buf)
	return items, nil
}

func (h *primitiveArrayHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	items, ok := obj.([]T)
	if !ok {
		return nil, fmt.Errorf("expected a %T, got %T", []T{}, obj)
	}

	arrayLen := uint32(len(items))
	if arrayLen == 0 {
		// empty array
		return nil, nil
	}

	bytes := h.converter.SliceToBytes(items)

	// allocate memory for the buffer
	bufferSize := uint32(len(bytes))
	bufferOffset, cln, err := wa.AllocateMemory(ctx, bufferSize)
	if err != nil {
		return cln, err
	}

	// write the buffer
	if ok := wa.Memory().Write(bufferOffset, bytes); !ok {
		return cln, fmt.Errorf("failed to write array data for %s", h.info.TypeName())
	}

	// write array object
	if ok := wa.Memory().WriteUint32Le(offset, bufferOffset); !ok {
		return cln, errors.New("failed to write array buffer pointer")
	}
	if ok := wa.Memory().WriteUint32Le(offset+4, bufferOffset); !ok {
		return cln, errors.New("failed to write array data start pointer")
	}
	if ok := wa.Memory().WriteUint32Le(offset+8, bufferSize); !ok {
		return cln, errors.New("failed to write array bytes length")
	}
	if ok := wa.Memory().WriteUint32Le(offset+12, arrayLen); !ok {
		return cln, errors.New("failed to write array length")
	}

	return cln, nil
}

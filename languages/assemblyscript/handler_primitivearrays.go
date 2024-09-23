/*
 * Copyright 2024 Hypermode, Inc.
 */

package assemblyscript

import (
	"context"
	"errors"
	"fmt"

	"hypruntime/langsupport"
	"hypruntime/langsupport/primitives"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/array.ts

func (p *planner) NewPrimitiveArrayHandler(ti langsupport.TypeInfo) (managedTypeHandler, error) {
	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}

	switch ti.ListElementType().Name() {
	case "bool":
		return newPrimitiveArrayHandler[bool](ti, typeDef), nil
	case "u8":
		return newPrimitiveArrayHandler[uint8](ti, typeDef), nil
	case "u16":
		return newPrimitiveArrayHandler[uint16](ti, typeDef), nil
	case "u32":
		return newPrimitiveArrayHandler[uint32](ti, typeDef), nil
	case "u64":
		return newPrimitiveArrayHandler[uint64](ti, typeDef), nil
	case "i8":
		return newPrimitiveArrayHandler[int8](ti, typeDef), nil
	case "i16":
		return newPrimitiveArrayHandler[int16](ti, typeDef), nil
	case "i32":
		return newPrimitiveArrayHandler[int32](ti, typeDef), nil
	case "i64":
		return newPrimitiveArrayHandler[int64](ti, typeDef), nil
	case "f32":
		return newPrimitiveArrayHandler[float32](ti, typeDef), nil
	case "f64":
		return newPrimitiveArrayHandler[float64](ti, typeDef), nil
	case "isize":
		return newPrimitiveArrayHandler[int](ti, typeDef), nil
	case "usize":
		return newPrimitiveArrayHandler[uint](ti, typeDef), nil
	default:
		return nil, fmt.Errorf("unsupported primitive array type: %s", ti.Name())
	}
}

func newPrimitiveArrayHandler[T primitive](ti langsupport.TypeInfo, typeDef *metadata.TypeDefinition) *primitiveArrayHandler[T] {
	return &primitiveArrayHandler[T]{
		*NewTypeHandler(ti),
		typeDef,
		primitives.NewPrimitiveTypeConverter[T](),
	}
}

type primitiveArrayHandler[T primitive] struct {
	typeHandler
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
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
		return cln, fmt.Errorf("failed to write array data for %s", h.typeInfo.Name())
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

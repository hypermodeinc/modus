/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"hypruntime/langsupport"
	"hypruntime/langsupport/primitives"
	"hypruntime/utils"

	"golang.org/x/exp/constraints"
)

type primitive interface {
	constraints.Integer | constraints.Float | ~bool
}

func (p *planner) NewPrimitiveHandler(typ string, rt reflect.Type) (h langsupport.TypeHandler, err error) {
	defer func() {
		if err == nil {
			p.typeHandlers[typ] = h
		}
	}()

	switch typ {
	case "bool":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 1, 1)
		converter := primitives.NewPrimitiveTypeConverter[bool]()
		return &primitiveHandler[bool]{info, converter}, nil
	case "uint8", "byte":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 1, 1)
		converter := primitives.NewPrimitiveTypeConverter[uint8]()
		return &primitiveHandler[uint8]{info, converter}, nil
	case "uint16":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 2, 1)
		converter := primitives.NewPrimitiveTypeConverter[uint16]()
		return &primitiveHandler[uint16]{info, converter}, nil
	case "uint32":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[uint32]()
		return &primitiveHandler[uint32]{info, converter}, nil
	case "uint64":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 8, 1)
		converter := primitives.NewPrimitiveTypeConverter[uint64]()
		return &primitiveHandler[uint64]{info, converter}, nil
	case "int8":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 1, 1)
		converter := primitives.NewPrimitiveTypeConverter[int8]()
		return &primitiveHandler[int8]{info, converter}, nil
	case "int16":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 2, 1)
		converter := primitives.NewPrimitiveTypeConverter[int16]()
		return &primitiveHandler[int16]{info, converter}, nil
	case "int32", "rune":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[int32]()
		return &primitiveHandler[int32]{info, converter}, nil
	case "int64":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 8, 1)
		converter := primitives.NewPrimitiveTypeConverter[int64]()
		return &primitiveHandler[int64]{info, converter}, nil
	case "float32":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[float32]()
		return &primitiveHandler[float32]{info, converter}, nil
	case "float64":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 8, 1)
		converter := primitives.NewPrimitiveTypeConverter[float64]()
		return &primitiveHandler[float64]{info, converter}, nil
	case "int":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[int]()
		return &primitiveHandler[int]{info, converter}, nil
	case "uint":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[uint]()
		return &primitiveHandler[uint]{info, converter}, nil
	case "uintptr":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
		converter := primitives.NewPrimitiveTypeConverter[uintptr]()
		return &primitiveHandler[uintptr]{info, converter}, nil
	case "time.Duration":
		info := langsupport.NewTypeHandlerInfo(typ, rt, 8, 1)
		converter := primitives.NewPrimitiveTypeConverter[time.Duration]()
		return &primitiveHandler[time.Duration]{info, converter}, nil
	default:
		return nil, fmt.Errorf("unsupported primitive type: %s", typ)
	}
}

type primitiveHandler[T primitive] struct {
	info      langsupport.TypeHandlerInfo
	converter primitives.TypeConverter[T]
}

func (h *primitiveHandler[T]) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *primitiveHandler[T]) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	val, ok := h.converter.Read(wa.Memory(), offset)
	if !ok {
		return 0, fmt.Errorf("failed to read %s from memory", h.info.TypeName())
	}

	return val, nil
}

func (h *primitiveHandler[T]) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	val, err := utils.Cast[T](obj)
	if err != nil {
		return nil, err
	}

	if ok := h.converter.Write(wa.Memory(), offset, val); !ok {
		return nil, fmt.Errorf("failed to write %s to memory", h.info.TypeName())
	}

	return nil, nil
}

func (h *primitiveHandler[T]) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("expected 1 value, got %d", len(vals))
	}

	return h.converter.Decode(vals[0]), nil
}

func (h *primitiveHandler[T]) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	val, err := utils.Cast[T](obj)
	if err != nil {
		return nil, nil, err
	}

	return []uint64{h.converter.Encode(val)}, nil, nil
}

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
	"hypruntime/logger"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

const maxDepth = 5 // TODO: make this based on the depth requested in the query

func (p *planner) NewManagedObjectHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[managedObjectHandler](p, typ)
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 4, 1)
	handler.nullable = _typeInfo.IsNullable(typ)

	typ = _typeInfo.GetUnderlyingType(typ)
	if typeDef, err := p.metadata.GetTypeDefinition(typ); err != nil {
		return nil, err
	} else {
		handler.typeDef = typeDef
	}

	newInnerHandler := func() (managedTypeHandler, error) {

		// nullable types use the same handler as the underlying type,
		// but are passed as a pointer on the runtime side
		kind := rt.Kind()
		if handler.nullable && rt.Kind() == reflect.Ptr {
			kind = rt.Elem().Kind()
		}

		switch kind {
		case reflect.Slice, reflect.Array:
			if _typeInfo.IsTypedArrayType(typ) {
				return p.NewTypedArrayHandler(typ, rt)
			} else {
				elementType := _typeInfo.GetListSubtype(typ)
				if _typeInfo.IsPrimitiveType(elementType) {
					return p.NewPrimitiveArrayHandler(ctx, typ, rt)
				} else {
					return p.NewArrayHandler(ctx, typ, rt)
				}
			}
		case reflect.Map:
			if _typeInfo.IsMapType(typ) {
				return p.NewMapHandler(ctx, typ, rt)
			} else {
				// This is a class that is being passed as a map.
				return p.NewClassHandler(ctx, typ, rt)
			}
		case reflect.Struct:
			if _typeInfo.IsTimestampType(typ) {
				return p.NewDateHandler(typ, rt)
			} else {
				return p.NewClassHandler(ctx, typ, rt)
			}
		}

		return nil, fmt.Errorf("unsupported type for managed object: %s", rt)
	}
	if h, err := newInnerHandler(); err != nil {
		return nil, err
	} else {
		handler.innerHandler = h
	}

	return handler, nil
}

type managedTypeHandler interface {
	Info() langsupport.TypeHandlerInfo
	Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error)
	Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error)
}

type managedObjectHandler struct {
	info         langsupport.TypeHandlerInfo
	typeDef      *metadata.TypeDefinition
	innerHandler managedTypeHandler
	nullable     bool
}

func (h *managedObjectHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *managedObjectHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, fmt.Errorf("unexpected address 0 reading managed object of type %s", h.info.TypeName())
	}

	// Check for recursion
	visitedPtrs := wa.(*wasmAdapter).visitedPtrs
	if visitedPtrs[offset] >= maxDepth {
		logger.Warn(ctx).Bool("user_visible", true).Msgf("Excessive recursion detected in %s. Stopping at depth %d.", h.info.TypeName(), maxDepth)
		return nil, nil
	}
	visitedPtrs[offset]++
	defer func() {
		n := visitedPtrs[offset]
		if n == 1 {
			delete(visitedPtrs, offset)
		} else {
			visitedPtrs[offset] = n - 1
		}
	}()

	ptr, ok := wa.Memory().ReadUint32Le(offset)
	if !ok {
		return nil, errors.New("failed to read managed object pointer")
	}

	return h.doRead(ctx, wa, ptr)
}

func (h *managedObjectHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	ptr, cln, err := h.doWrite(ctx, wa, obj)
	if err != nil {
		return cln, err
	}

	if ok := wa.Memory().WriteUint32Le(offset, ptr); !ok {
		return cln, fmt.Errorf("failed to write managed object pointer to WASM memory")
	}

	return cln, nil
}

func (h *managedObjectHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	if len(vals) != 1 {
		return nil, fmt.Errorf("expected 1 value when decoding a managed object, got %d", len(vals))
	}

	return h.doRead(ctx, wa, uint32(vals[0]))
}

func (h *managedObjectHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	ptr, cln, err := h.doWrite(ctx, wa, obj)
	if err != nil {
		return nil, cln, err
	}

	return []uint64{uint64(ptr)}, cln, nil
}

func (h *managedObjectHandler) doRead(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		if h.nullable {
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected null pointer for non-nullable type %s", h.info.TypeName())
		}
	}

	return h.innerHandler.Read(ctx, wa, offset)
}

func (h *managedObjectHandler) doWrite(ctx context.Context, wa langsupport.WasmAdapter, obj any) (uint32, utils.Cleaner, error) {
	if utils.HasNil(obj) {
		if h.nullable {
			return 0, nil, nil
		} else {
			return 0, nil, fmt.Errorf("unexpected nil value for non-nullable type %s", h.info.TypeName())
		}
	}

	if h.nullable {
		obj = utils.DereferencePointer(obj)
	}

	id := h.typeDef.Id
	size := h.innerHandler.Info().TypeSize()

	ptr, cln, err := wa.(*wasmAdapter).allocateAndPinMemory(ctx, size, id)
	if err != nil {
		return 0, cln, err
	}

	c, err := h.innerHandler.Write(ctx, wa, ptr, obj)
	if err != nil {
		cln.AddCleaner(c)
		return 0, cln, err
	}

	// we can unpin the inner objects early, since they are part of the managed object
	// if err := c.Clean(); err != nil {
	// 	return 0, cln, err
	// }

	return ptr, cln, nil
}

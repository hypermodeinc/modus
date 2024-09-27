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

func (p *planner) NewManagedObjectHandler(ctx context.Context, ti langsupport.TypeInfo) (langsupport.TypeHandler, error) {

	handler := &managedObjectHandler{
		typeHandler: *NewTypeHandler(ti),
	}
	p.AddHandler(handler)

	ut := ti.UnderlyingType()
	if ut != nil {
		ti = ut
	}

	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	newInnerHandler := func() (managedTypeHandler, error) {

		// nullable types use the same handler as the underlying type,
		// but are passed as a pointer on the runtime side
		kind := ti.ReflectedType().Kind()
		if ti.IsNullable() && kind == reflect.Ptr {
			kind = ti.ReflectedType().Elem().Kind()
		}

		switch kind {
		case reflect.Slice, reflect.Array:
			if _langTypeInfo.IsTypedArrayType(ti.Name()) {
				return p.NewTypedArrayHandler(ti)
			} else if ti.ListElementType().IsPrimitive() {
				return p.NewPrimitiveArrayHandler(ti)
			} else {
				return p.NewArrayHandler(ctx, ti)
			}
		case reflect.Map:
			if ti.IsMap() {
				return p.NewMapHandler(ctx, ti)
			} else {
				// This is a class that is being passed as a map.
				return p.NewClassHandler(ctx, ti)
			}
		case reflect.Struct:
			if ti.IsTimestamp() {
				return p.NewDateHandler(ti)
			} else {
				return p.NewClassHandler(ctx, ti)
			}
		}

		return nil, fmt.Errorf("unsupported type for managed object: %s", ti.Name())
	}
	if h, err := newInnerHandler(); err != nil {
		return nil, err
	} else {
		handler.innerHandler = h
	}

	return handler, nil
}

type managedTypeHandler interface {
	Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error)
	Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error)
}

type managedObjectHandler struct {
	typeHandler
	typeDef      *metadata.TypeDefinition
	innerHandler managedTypeHandler
}

func (h *managedObjectHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, fmt.Errorf("unexpected address 0 reading managed object of type %s", h.typeInfo.Name())
	}

	// Check for recursion
	visitedPtrs := wa.(*wasmAdapter).visitedPtrs
	if visitedPtrs[offset] >= maxDepth {
		logger.Warn(ctx).Bool("user_visible", true).Msgf("Excessive recursion detected in %s. Stopping at depth %d.", h.typeInfo.Name(), maxDepth)
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
		if h.typeInfo.IsNullable() {
			return nil, nil
		} else {
			return nil, fmt.Errorf("unexpected null pointer for non-nullable type %s", h.typeInfo.Name())
		}
	}

	return h.innerHandler.Read(ctx, wa, offset)
}

func (h *managedObjectHandler) doWrite(ctx context.Context, wa langsupport.WasmAdapter, obj any) (uint32, utils.Cleaner, error) {
	if utils.HasNil(obj) {
		if h.typeInfo.IsNullable() {
			return 0, nil, nil
		} else {
			return 0, nil, fmt.Errorf("unexpected nil value for non-nullable type %s", h.typeInfo.Name())
		}
	}

	if h.typeInfo.IsNullable() {
		obj = utils.DereferencePointer(obj)
	}

	ptr, cln, err := wa.(*wasmAdapter).allocateAndPinMemory(ctx, h.typeInfo.DataSize(), h.typeDef.Id)
	if err != nil {
		return 0, cln, err
	}

	c, err := h.innerHandler.Write(ctx, wa, ptr, obj)
	if err != nil {
		cln.AddCleaner(c)
		return 0, cln, err
	}

	// TODO: the following should work, but it in some cases we are finding that the inner objects fail to unpin even though they were successfully pinned

	// we can unpin the inner objects early, since they are now referenced by the managed object
	// if c != nil {
	// 	if err := c.Clean(); err != nil {
	// 		return 0, cln, err
	// 	}
	// }

	return ptr, cln, nil
}

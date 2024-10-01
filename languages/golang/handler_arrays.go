/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"reflect"

	"hypruntime/langsupport"
	"hypruntime/utils"
)

func (p *planner) NewArrayHandler(ctx context.Context, ti langsupport.TypeInfo) (langsupport.TypeHandler, error) {
	handler := &arrayHandler{
		typeHandler: *NewTypeHandler(ti),
	}
	p.AddHandler(handler)

	elementHandler, err := p.GetHandler(ctx, ti.ListElementType().Name())
	if err != nil {
		return nil, err
	}
	handler.elementHandler = elementHandler

	arrayLen, err := _langTypeInfo.ArrayLength(ti.Name())
	if err != nil {
		return nil, err
	}
	handler.arrayLen = arrayLen

	return handler, nil
}

type arrayHandler struct {
	typeHandler
	elementHandler langsupport.TypeHandler
	arrayLen       int
}

func (h *arrayHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	elementSize := h.elementHandler.TypeInfo().Size()
	items := reflect.New(h.typeInfo.ReflectedType()).Elem()
	for i := 0; i < h.arrayLen; i++ {
		itemOffset := offset + uint32(i)*elementSize
		item, err := h.elementHandler.Read(ctx, wa, itemOffset)
		if err != nil {
			return nil, err
		}
		items.Index(i).Set(reflect.ValueOf(item))
	}

	return items.Interface(), nil
}

func (h *arrayHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	items, err := utils.ConvertToSlice(obj)
	if err != nil {
		return nil, err
	}

	cln := utils.NewCleaner()
	elementSize := h.elementHandler.TypeInfo().Size()

	// write exactly the number of items that will fit in the array
	for i := 0; i < h.arrayLen; i++ {
		if i >= len(items) {
			break
		}
		item := items[i]
		c, err := h.elementHandler.Write(ctx, wa, offset, item)
		cln.AddCleaner(c)
		if err != nil {
			return cln, err
		}
		offset += elementSize
	}

	// zero out any remaining space in the array
	remainingItems := h.arrayLen - len(items)
	if remainingItems > 0 {
		zeros := make([]byte, remainingItems*int(elementSize))
		if ok := wa.Memory().Write(offset, zeros); !ok {
			return nil, errors.New("failed to zero out remaining array space")
		}
	}

	return cln, nil
}

func (h *arrayHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	array := reflect.New(h.typeInfo.ReflectedType()).Elem()
	itemLen := int(h.elementHandler.TypeInfo().EncodingLength())
	for i := 0; i < h.arrayLen; i++ {
		data, err := h.elementHandler.Decode(ctx, wa, vals[i*itemLen:(i+1)*itemLen])
		if err != nil {
			return nil, err
		}
		array.Index(i).Set(reflect.ValueOf(data))
	}
	return array.Interface(), nil
}

func (h *arrayHandler) Encode(ctx context.Context, wa langsupport.WasmAdapter, obj any) ([]uint64, utils.Cleaner, error) {
	items, err := utils.ConvertToSlice(obj)
	if err != nil {
		return nil, nil, err
	}

	itemLen := int(h.elementHandler.TypeInfo().EncodingLength())
	res := make([]uint64, h.typeInfo.EncodingLength())
	cln := utils.NewCleaner()

	// fill the array with the encoded values, too few will be zeroed, too many will be truncated
	for i, item := range items {
		if i >= h.arrayLen {
			break
		}

		vals, c, err := h.elementHandler.Encode(ctx, wa, item)
		cln.AddCleaner(c)
		if err != nil {
			return nil, cln, err
		}

		copy(res[i*itemLen:(i+1)*itemLen], vals)
	}

	return res, cln, nil
}

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

func (p *planner) NewArrayHandler(ctx context.Context, typ string, rt reflect.Type) (langsupport.TypeHandler, error) {
	handler := NewTypeHandler[arrayHandler](p, typ)

	elementType := _typeInfo.GetListSubtype(typ)
	elementHandler, err := p.GetHandler(ctx, elementType)
	if err != nil {
		return nil, err
	}
	handler.elementHandler = elementHandler

	arrayLen, err := _typeInfo.ArrayLength(typ)
	if err != nil {
		return nil, err
	}

	elementSize := elementHandler.Info().TypeSize()
	typeSize := uint32(arrayLen) * elementSize
	encodingLen := arrayLen * elementHandler.Info().EncodingLength()

	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, typeSize, encodingLen)
	handler.arrayLen = arrayLen
	handler.elementSize = elementSize

	return handler, nil
}

type arrayHandler struct {
	info           langsupport.TypeHandlerInfo
	elementHandler langsupport.TypeHandler
	arrayLen       int
	elementSize    uint32
}

func (h *arrayHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *arrayHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	items := reflect.New(h.info.RuntimeType()).Elem()
	for i := 0; i < h.arrayLen; i++ {
		itemOffset := offset + uint32(i)*h.elementSize
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
		offset += h.elementSize
	}

	// zero out any remaining space in the array
	remainingItems := h.arrayLen - len(items)
	if remainingItems > 0 {
		zeros := make([]byte, remainingItems*int(h.elementSize))
		if ok := wa.Memory().Write(offset, zeros); !ok {
			return nil, errors.New("failed to zero out remaining array space")
		}
	}

	return cln, nil
}

func (h *arrayHandler) Decode(ctx context.Context, wa langsupport.WasmAdapter, vals []uint64) (any, error) {
	array := reflect.New(h.info.RuntimeType()).Elem()
	itemEncLen := h.elementHandler.Info().EncodingLength()
	for i := 0; i < h.arrayLen; i++ {
		data, err := h.elementHandler.Decode(ctx, wa, vals[i*itemEncLen:(i+1)*itemEncLen])
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

	itemLen := h.elementHandler.Info().EncodingLength()
	res := make([]uint64, h.info.EncodingLength())
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

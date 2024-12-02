/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package assemblyscript

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/utils"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/array.ts

func (p *planner) NewArrayHandler(ctx context.Context, ti langsupport.TypeInfo) (managedTypeHandler, error) {
	handler := &arrayHandler{
		typeHandler: *NewTypeHandler(ti),
	}

	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	elementHandler, err := p.GetHandler(ctx, ti.ListElementType().Name())
	if err != nil {
		return nil, err
	}
	handler.elementHandler = elementHandler

	handler.emptyValue = reflect.MakeSlice(ti.ReflectedType(), 0, 0).Interface()

	return handler, nil
}

type arrayHandler struct {
	typeHandler
	typeDef        *metadata.TypeDefinition
	elementHandler langsupport.TypeHandler
	emptyValue     any
}

func (h *arrayHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
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
		return h.emptyValue, nil
	}

	elementSize := h.elementHandler.TypeInfo().Size()
	items := reflect.MakeSlice(h.typeInfo.ReflectedType(), int(arrLen), int(arrLen))
	for i := uint32(0); i < arrLen; i++ {
		itemOffset := data + i*elementSize
		item, err := h.elementHandler.Read(ctx, wa, itemOffset)
		if err != nil {
			return nil, err
		}
		items.Index(int(i)).Set(reflect.ValueOf(item))
	}

	return items.Interface(), nil
}

func (h *arrayHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	items, err := utils.ConvertToSlice(obj)
	if err != nil {
		return nil, err
	}

	arrLen := uint32(len(items))
	if arrLen == 0 {
		// empty array
		return nil, nil
	}

	// allocate memory for the buffer
	elementSize := h.elementHandler.TypeInfo().Size()
	bufferSize := arrLen * elementSize
	bufferOffset, cln, err := wa.AllocateMemory(ctx, bufferSize)
	if err != nil {
		return cln, err
	}

	// write the elements to the buffer
	for i := uint32(0); i < arrLen; i++ {
		itemOffset := bufferOffset + (elementSize * i)
		c, err := h.elementHandler.Write(ctx, wa, itemOffset, items[i])
		cln.AddCleaner(c)
		if err != nil {
			return cln, fmt.Errorf("failed to write array item: %w", err)
		}
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

	if ok := wa.Memory().WriteUint32Le(offset+12, arrLen); !ok {
		return cln, errors.New("failed to write array length")
	}

	return cln, nil
}

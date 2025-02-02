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
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/hypermodeinc/modus/lib/metadata"
	"github.com/hypermodeinc/modus/runtime/langsupport"
	"github.com/hypermodeinc/modus/runtime/langsupport/primitives"
	"github.com/hypermodeinc/modus/runtime/utils"

	"golang.org/x/exp/constraints"
)

func (p *planner) NewTypedArrayHandler(ti langsupport.TypeInfo) (managedTypeHandler, error) {
	typeDef, err := p.metadata.GetTypeDefinition(ti.Name())
	if err != nil {
		return nil, err
	}

	switch ti.Name() {
	case "~lib/typedarray/Uint8Array", "~lib/typedarray/Uint8ClampedArray":
		return newTypedArrayHandler[uint8](ti, typeDef), nil
	case "~lib/typedarray/Uint16Array":
		return newTypedArrayHandler[uint16](ti, typeDef), nil
	case "~lib/typedarray/Uint32Array":
		return newTypedArrayHandler[uint32](ti, typeDef), nil
	case "~lib/typedarray/Uint64Array":
		return newTypedArrayHandler[uint64](ti, typeDef), nil
	case "~lib/typedarray/Int8Array":
		return newTypedArrayHandler[int8](ti, typeDef), nil
	case "~lib/typedarray/Int16Array":
		return newTypedArrayHandler[int16](ti, typeDef), nil
	case "~lib/typedarray/Int32Array":
		return newTypedArrayHandler[int32](ti, typeDef), nil
	case "~lib/typedarray/Int64Array":
		return newTypedArrayHandler[int64](ti, typeDef), nil
	case "~lib/typedarray/Float32Array":
		return newTypedArrayHandler[float32](ti, typeDef), nil
	case "~lib/typedarray/Float64Array":
		return newTypedArrayHandler[float64](ti, typeDef), nil
	default:
		return nil, fmt.Errorf("unsupported typed array type: %s", ti.Name())
	}
}

func newTypedArrayHandler[T constraints.Integer | constraints.Float](ti langsupport.TypeInfo, typeDef *metadata.TypeDefinition) *typedArrayHandler[T] {
	return &typedArrayHandler[T]{
		*NewTypeHandler(ti),
		typeDef,
		primitives.NewPrimitiveTypeConverter[T](),
	}
}

type typedArrayHandler[T constraints.Integer | constraints.Float] struct {
	typeHandler
	typeDef   *metadata.TypeDefinition
	converter primitives.TypeConverter[T]
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
		return nil, fmt.Errorf("failed to read data start address for %s", h.typeInfo.Name())
	}

	byteLen, ok := wa.Memory().ReadUint32Le(offset + 8)
	if !ok {
		return nil, fmt.Errorf("failed to read byte length for %s", h.typeInfo.Name())
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
	var bytes []byte
	if s, ok := obj.(string); ok {
		if b, err := base64.StdEncoding.DecodeString(s); err != nil {
			return nil, utils.NewUserError(fmt.Errorf("failed to decode base64 string: %w", err))
		} else {
			bytes = b
		}
	} else {
		items, ok := utils.ConvertToSliceOf[T](obj)
		if !ok {
			return nil, fmt.Errorf("input is invalid for type %s", h.typeInfo.Name())
		} else if len(items) == 0 {
			// empty typed array
			return nil, nil
		}
		bytes = h.converter.SliceToBytes(items)
	}

	// allocate memory for the buffer
	bufferSize := uint32(len(bytes))
	ptr, cln, err := wa.AllocateMemory(ctx, bufferSize)
	if err != nil {
		return cln, err
	}

	// write the buffer
	if ok := wa.Memory().Write(ptr, bytes); !ok {
		return cln, fmt.Errorf("failed to write typed array data for %s", h.typeInfo.Name())
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

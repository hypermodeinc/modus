/*
 * Copyright 2024 Hypermode, Inc.
 */

package golang

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"hypruntime/utils"
)

func (wa *wasmAdapter) encodeArray(ctx context.Context, typ string, obj any) ([]uint64, utils.Cleaner, error) {
	arrLen, err := wa.typeInfo.ArraySize(typ)
	if err != nil {
		return nil, nil, err
	}

	items, err := utils.ConvertToSlice(obj)
	if err != nil {
		return nil, nil, err
	}

	t := wa.typeInfo.GetListSubtype(typ)
	itemLen, err := wa.GetEncodingLength(ctx, t)
	if err != nil {
		return nil, nil, err
	}

	res := make([]uint64, arrLen*itemLen)
	cln := utils.NewCleaner()

	// fill the array with the encoded values, too few will be zeroed, too many will be truncated
	for i, item := range items {
		if i >= arrLen {
			break
		}

		vals, c, err := wa.encodeObject(ctx, t, item)
		cln.AddCleaner(c)
		if err != nil {
			return nil, cln, err
		}

		copy(res[i*itemLen:(i+1)*itemLen], vals)
	}

	return res, cln, nil
}

func (wa *wasmAdapter) decodeArray(ctx context.Context, typ string, vals []uint64) (any, error) {
	arrLen, err := wa.typeInfo.ArraySize(typ)
	if err != nil {
		return nil, err
	}

	itemType := wa.typeInfo.GetListSubtype(typ)
	goItemType, err := wa.getReflectedType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	items := reflect.New(reflect.ArrayOf(arrLen, goItemType)).Elem()

	if arrLen == 1 {
		data, err := wa.decodeObject(ctx, itemType, vals)
		if err != nil {
			return nil, err
		}
		items.Index(0).Set(reflect.ValueOf(data))
	} else if arrLen > 1 {
		// This shouldn't be reachable, because multiple result values are not supported by TinyGo,
		// and an array with length > 1 would require multiple values.
		return nil, errors.New("decoding arrays with more than one element is not yet supported")
	}

	return items.Interface(), nil
}

func (wa *wasmAdapter) getArrayEncodedLength(ctx context.Context, typ string) (int, error) {
	arrSize, err := wa.typeInfo.ArraySize(typ)
	if err != nil {
		return 0, err
	}

	t := wa.typeInfo.GetListSubtype(typ)
	itemSize, err := wa.GetEncodingLength(ctx, t)
	if err != nil {
		return 0, err
	}

	return arrSize * itemSize, nil
}

func (wa *wasmAdapter) writeArray(ctx context.Context, typ string, offset uint32, obj any) (utils.Cleaner, error) {
	itemType := wa.typeInfo.GetListSubtype(typ)
	arrLen, err := wa.typeInfo.ArraySize(typ)
	if err != nil {
		return nil, err
	}

	mem := wa.mod.Memory()

	// special case for byte arrays, because they can be written more efficiently
	if itemType == "byte" || itemType == "uint8" {
		bytes, err := convertToByteSlice(obj)
		if err != nil {
			return nil, err
		}

		// the array must be exactly the right length
		if len(bytes) < arrLen {
			bytes = append(bytes, make([]byte, arrLen-len(bytes))...)
		} else if len(bytes) > arrLen {
			bytes = bytes[:arrLen]
		}

		if ok := mem.Write(offset, bytes); !ok {
			return nil, fmt.Errorf("failed to write byte array to WASM memory")
		}
		return nil, nil
	}

	items, err := utils.ConvertToSlice(obj)
	if err != nil {
		return nil, err
	}

	itemSize, err := wa.typeInfo.GetSizeOfType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	cln := utils.NewCleaner()

	// write exactly the number of items that will fit in the array
	for i := 0; i < arrLen; i++ {
		if i >= len(items) {
			break
		}
		item := items[i]
		c, err := wa.writeObject(ctx, itemType, offset, item)
		cln.AddCleaner(c)
		if err != nil {
			return cln, err
		}
		offset += itemSize
	}

	// zero out any remaining space in the array
	remainingItems := arrLen - len(items)
	if remainingItems > 0 {
		zeros := make([]byte, remainingItems*int(itemSize))
		if ok := mem.Write(offset, zeros); !ok {
			return nil, fmt.Errorf("failed to zero out remaining array space")
		}
	}

	return cln, nil
}

func (wa *wasmAdapter) readArray(ctx context.Context, typ string, offset uint32) (any, error) {
	arrLen, err := wa.typeInfo.ArraySize(typ)
	if err != nil {
		return nil, err
	}

	itemType := wa.typeInfo.GetListSubtype(typ)

	// special case for byte arrays, because they can be read more efficiently
	if itemType == "byte" || itemType == "uint8" {
		bytes, ok := wa.mod.Memory().Read(offset, uint32(arrLen))
		if !ok {
			return nil, errors.New("failed to read data from WASM memory")
		}
		return bytes, nil
	}

	itemSize, err := wa.typeInfo.GetSizeOfType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	goItemType, err := wa.getReflectedType(ctx, itemType)
	if err != nil {
		return nil, err
	}

	items := reflect.New(reflect.ArrayOf(arrLen, goItemType)).Elem()
	for i := 0; i < arrLen; i++ {
		itemOffset := offset + uint32(i)*itemSize
		item, err := wa.readObject(ctx, itemType, itemOffset)
		if err != nil {
			return nil, err
		}
		items.Index(i).Set(reflect.ValueOf(item))

	}

	return items.Interface(), nil
}

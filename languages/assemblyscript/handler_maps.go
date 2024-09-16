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
	"hypruntime/languages/assemblyscript/hash"
	"hypruntime/plugins/metadata"
	"hypruntime/utils"
)

// Reference: https://github.com/AssemblyScript/assemblyscript/blob/main/std/assembly/map.ts

func (p *planner) NewMapHandler(ctx context.Context, typ string, rt reflect.Type) (managedTypeHandler, error) {
	handler := new(mapHandler)
	handler.info = langsupport.NewTypeHandlerInfo(typ, rt, 24, 0)

	typ = _typeInfo.GetUnderlyingType(typ)
	typeDef, err := p.metadata.GetTypeDefinition(typ)
	if err != nil {
		return nil, err
	}
	handler.typeDef = typeDef

	keyType, valueType := _typeInfo.GetMapSubtypes(typ)
	if keyType == "" || valueType == "" {
		return nil, errors.New("map type must have key and value subtypes")
	}

	keyHandler, err := p.GetHandler(ctx, keyType)
	if err != nil {
		return nil, err
	}
	handler.keyHandler = keyHandler

	valueHandler, err := p.GetHandler(ctx, valueType)
	if err != nil {
		return nil, err
	}
	handler.valueHandler = valueHandler

	rtKey := keyHandler.Info().RuntimeType()
	rtValue := valueHandler.Info().RuntimeType()
	if !rtKey.Comparable() {
		handler.usePseudoMap = true
		handler.rtPseudoMapSlice = reflect.SliceOf(reflect.StructOf([]reflect.StructField{
			{
				Name: "Key",
				Type: rtKey,
				Tag:  `json:"key"`,
			},
			{
				Name: "Value",
				Type: rtValue,
				Tag:  `json:"value"`,
			},
		}))

		handler.rtPseudoMap = reflect.StructOf([]reflect.StructField{
			{
				Name: "Data",
				Type: handler.rtPseudoMapSlice,
				Tag:  `json:"$mapdata"`,
			},
		})
	}

	return handler, nil
}

type mapHandler struct {
	info             langsupport.TypeHandlerInfo
	typeDef          *metadata.TypeDefinition
	keyHandler       langsupport.TypeHandler
	valueHandler     langsupport.TypeHandler
	usePseudoMap     bool
	rtPseudoMap      reflect.Type
	rtPseudoMapSlice reflect.Type
}

func (h *mapHandler) Info() langsupport.TypeHandlerInfo {
	return h.info
}

func (h *mapHandler) Read(ctx context.Context, wa langsupport.WasmAdapter, offset uint32) (any, error) {
	if offset == 0 {
		return nil, nil
	}

	// buckets, ok := wa.Memory().ReadUint32Le(offset)
	// if !ok {
	// 	return nil, errors.New("failed to read map buckets pointer")
	// }

	// bucketsMask, ok := wa.Memory().ReadUint32Le(offset + 4)
	// if !ok {
	// 	return nil, errors.New("failed to read map buckets mask")
	// }

	entries, ok := wa.Memory().ReadUint32Le(offset + 8)
	if !ok {
		return nil, errors.New("failed to read map entries pointer")
	}

	entriesCapacity, ok := wa.Memory().ReadUint32Le(offset + 12)
	if !ok {
		return nil, errors.New("failed to read map entries capacity")
	}

	// entriesOffset, ok := wa.Memory().ReadUint32Le(offset + 16)
	// if !ok {
	// 	return nil, errors.New("failed to read map entries offset")
	// }

	entriesCount, ok := wa.Memory().ReadUint32Le(offset + 20)
	if !ok {
		return nil, errors.New("failed to read map entries count")
	}

	// the length of array buffer is stored 4 bytes before the offset
	byteLength, ok := wa.Memory().ReadUint32Le(entries - 4)
	if !ok {
		return nil, errors.New("failed to read map entries buffer length")
	}

	mapSize := int(entriesCount)
	entrySize := byteLength / entriesCapacity
	keySize := h.keyHandler.Info().TypeSize()
	valueOffset := max(keySize, 4)

	if !h.usePseudoMap {
		// return a map
		m := reflect.MakeMapWithSize(h.info.RuntimeType(), mapSize)
		for i := uint32(0); i < entriesCount; i++ {
			p := entries + (i * entrySize)

			k, err := h.keyHandler.Read(ctx, wa, p)
			if err != nil {
				return nil, err
			}

			v, err := h.valueHandler.Read(ctx, wa, p+valueOffset)
			if err != nil {
				return nil, err
			}

			m.SetMapIndex(reflect.ValueOf(k), reflect.ValueOf(v))
		}
		return m.Interface(), nil

	} else {
		// return a pseudo-map
		s := reflect.MakeSlice(h.rtPseudoMapSlice, mapSize, mapSize)
		for i := 0; i < mapSize; i++ {
			p := entries + uint32(i)*entrySize

			k, err := h.keyHandler.Read(ctx, wa, p)
			if err != nil {
				return nil, err
			}

			v, err := h.valueHandler.Read(ctx, wa, p+keySize)
			if err != nil {
				return nil, err
			}

			s.Index(i).Field(0).Set(reflect.ValueOf(k))
			s.Index(i).Field(1).Set(reflect.ValueOf(v))
		}

		m := reflect.New(h.rtPseudoMap).Elem()
		m.Field(0).Set(s)
		return m.Interface(), nil
	}
}

func (h *mapHandler) Write(ctx context.Context, wa langsupport.WasmAdapter, offset uint32, obj any) (utils.Cleaner, error) {
	rv := reflect.ValueOf(obj)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Map {
		return nil, fmt.Errorf("unsupported map type %T", obj)
	}
	mapLen := uint32(rv.Len())

	// determine capacities and mask
	bucketsCapacity := uint32(4)
	entriesCapacity := uint32(4)
	bucketsMask := bucketsCapacity - 1
	for bucketsCapacity < mapLen {
		bucketsCapacity <<= 1
		entriesCapacity = bucketsCapacity * 8 / 3
		bucketsMask = bucketsCapacity - 1
	}

	cln := utils.NewCleanerN(int(mapLen*2) + 1)

	// create buckets array buffer
	const bucketSize = 4
	bucketsBufferSize := bucketSize * bucketsCapacity
	bucketsBufferOffset, c, err := wa.AllocateMemory(ctx, bucketsBufferSize)
	cln.AddCleaner(c)
	if err != nil {
		return cln, fmt.Errorf("failed to allocate memory for array buffer: %w", err)
	}

	// write entries array buffer
	// note: unlike arrays, an empty map DOES have array buffers
	keySize := h.keyHandler.Info().TypeSize()
	valueSize := h.valueHandler.Info().TypeSize()
	valueOffset := max(keySize, 4)

	const taggedNextSize = 4
	taggedNextOffset := valueSize + valueOffset

	entryAlign := max(keySize, valueSize, 4) - 1
	entrySize := (keySize + valueSize + taggedNextSize + entryAlign) & ^entryAlign
	entriesBufferSize := entrySize * entriesCapacity
	entriesBufferOffset, c, err := wa.AllocateMemory(ctx, entriesBufferSize)
	cln.AddCleaner(c)
	if err != nil {
		return cln, fmt.Errorf("failed to allocate memory for array buffer: %w", err)
	}

	mapKeys := rv.MapKeys()
	for i, mapKey := range mapKeys {

		entryOffset := entriesBufferOffset + (entrySize * uint32(i))

		// write entry key and calculate hash code
		var hashCode uint32
		key := mapKey.Interface()

		switch t := key.(type) {
		case string:
			// Special case for string keys, to avoid encoding to UTF16 twice.
			bytes := utils.EncodeUTF16(t)
			hashCode = hash.GetHashCode(bytes)

			ptr, c, err := h.keyHandler.(*stringHandler).doWriteBytes(ctx, wa, bytes)
			cln.AddCleaner(c)
			if err != nil {
				return cln, errors.New("failed to write map entry key")
			}
			if ok := wa.Memory().WriteUint32Le(entryOffset, ptr); !ok {
				return cln, errors.New("failed to write map entry key pointer")
			}

		default:
			hashCode = hash.GetHashCode(key)
			c, err := h.keyHandler.Write(ctx, wa, entryOffset, key)
			cln.AddCleaner(c)
			if err != nil {
				return cln, fmt.Errorf("failed to write map entry key: %w", err)
			}
		}

		// write entry value
		mapValue := rv.MapIndex(mapKey)
		value := mapValue.Interface()
		entryValueOffset := entryOffset + valueOffset
		c, err := h.valueHandler.Write(ctx, wa, entryValueOffset, value)
		cln.AddCleaner(c)
		if err != nil {
			return cln, fmt.Errorf("failed to write map entry value: %w", err)
		}

		// write to bucket and "tagged next" field
		bucketPtrBase := bucketsBufferOffset + ((hashCode & bucketsMask) * bucketSize)

		if prev, ok := wa.Memory().ReadUint32Le(bucketPtrBase); !ok {
			return cln, errors.New("failed to read previous map entry bucket pointer")
		} else if ok := wa.Memory().WriteUint32Le(entryOffset+taggedNextOffset, prev); !ok {
			return cln, errors.New("failed to write map entry tagged next field")
		}

		if ok := wa.Memory().WriteUint32Le(bucketPtrBase, entryOffset); !ok {
			return cln, errors.New("failed to write map entry bucket pointer")
		}
	}

	if ok := wa.Memory().WriteUint32Le(offset, bucketsBufferOffset); !ok {
		return cln, errors.New("failed to write map buckets pointer")
	}

	if ok := wa.Memory().WriteUint32Le(offset+4, bucketsMask); !ok {
		return cln, errors.New("failed to write map buckets mask")
	}

	if ok := wa.Memory().WriteUint32Le(offset+8, entriesBufferOffset); !ok {
		return cln, errors.New("failed to write map entries pointer")
	}

	if ok := wa.Memory().WriteUint32Le(offset+12, entriesCapacity); !ok {
		return cln, errors.New("failed to write map entries capacity")
	}

	if ok := wa.Memory().WriteUint32Le(offset+16, mapLen); !ok {
		return cln, errors.New("failed to write map entries offset")
	}

	if ok := wa.Memory().WriteUint32Le(offset+20, mapLen); !ok {
		return cln, errors.New("failed to write map entries count")
	}

	return cln, nil
}
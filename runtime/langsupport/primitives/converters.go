/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package primitives

import (
	"fmt"
	"time"
	"unsafe"

	wasm "github.com/tetratelabs/wazero/api"
	"golang.org/x/exp/constraints"
)

type primitive interface {
	constraints.Integer | constraints.Float | ~bool
}

type TypeConverter[T primitive] interface {
	TypeSize() int
	Read(wasm.Memory, uint32) (T, bool)
	Write(wasm.Memory, uint32, T) bool
	Decode(val uint64) T
	Encode(val T) uint64
	BytesToSlice([]byte) []T
	SliceToBytes([]T) []byte
}

type primitiveConverter[T primitive] struct {
	typeSize     int
	read         func(wasm.Memory, uint32) (T, bool)
	write        func(wasm.Memory, uint32, T) bool
	decode       func(uint64) T
	encode       func(T) uint64
	bytesToSlice func([]byte, int) []T
	sliceToBytes func([]T, int) []byte
}

func (c primitiveConverter[T]) TypeSize() int {
	return c.typeSize
}

func (c primitiveConverter[T]) Read(mem wasm.Memory, offset uint32) (T, bool) {
	return c.read(mem, offset)
}

func (c primitiveConverter[T]) Write(mem wasm.Memory, offset uint32, val T) bool {
	return c.write(mem, offset, val)
}

func (c primitiveConverter[T]) Decode(val uint64) T {
	return c.decode(val)
}

func (c primitiveConverter[T]) Encode(val T) uint64 {
	return c.encode(val)
}

func (c primitiveConverter[T]) BytesToSlice(b []byte) []T {
	return c.bytesToSlice(b, c.typeSize)
}

func (c primitiveConverter[T]) SliceToBytes(s []T) []byte {
	return c.sliceToBytes(s, c.typeSize)
}

func NewPrimitiveTypeConverter[T primitive]() TypeConverter[T] {
	var r any

	var c T
	switch any(c).(type) {
	case bool:
		r = primitiveConverter[bool]{1, ReadBool, WriteBool, DecodeBool, EncodeBool, bytesToSlice[bool], sliceToBytes[bool]}
	case int8:
		r = primitiveConverter[int8]{1, ReadInt8, WriteInt8, DecodeInt8, EncodeInt8, bytesToSlice[int8], sliceToBytes[int8]}
	case int16:
		r = primitiveConverter[int16]{2, ReadInt16, WriteInt16, DecodeInt16, EncodeInt16, bytesToSlice[int16], sliceToBytes[int16]}
	case int32:
		r = primitiveConverter[int32]{4, ReadInt32, WriteInt32, DecodeInt32, EncodeInt32, bytesToSlice[int32], sliceToBytes[int32]}
	case int64:
		r = primitiveConverter[int64]{8, ReadInt64, WriteInt64, DecodeInt64, EncodeInt64, bytesToSlice[int64], sliceToBytes[int64]}
	case uint8:
		r = primitiveConverter[uint8]{1, ReadUint8, WriteUint8, DecodeUint8, EncodeUint8, bytesToSlice[uint8], sliceToBytes[uint8]}
	case uint16:
		r = primitiveConverter[uint16]{2, ReadUint16, WriteUint16, DecodeUint16, EncodeUint16, bytesToSlice[uint16], sliceToBytes[uint16]}
	case uint32:
		r = primitiveConverter[uint32]{4, ReadUint32, WriteUint32, DecodeUint32, EncodeUint32, bytesToSlice[uint32], sliceToBytes[uint32]}
	case uint64:
		r = primitiveConverter[uint64]{8, ReadUint64, WriteUint64, DecodeUint64, EncodeUint64, bytesToSlice[uint64], sliceToBytes[uint64]}
	case float32:
		r = primitiveConverter[float32]{4, ReadFloat32, WriteFloat32, DecodeFloat32, EncodeFloat32, bytesToSlice[float32], sliceToBytes[float32]}
	case float64:
		r = primitiveConverter[float64]{8, ReadFloat64, WriteFloat64, DecodeFloat64, EncodeFloat64, bytesToSlice[float64], sliceToBytes[float64]}
	case int:
		r = primitiveConverter[int]{4, ReadInt, WriteInt, DecodeInt, EncodeInt, bytesToSliceFixed32[int], sliceToBytesFixed32[int]}
	case uint:
		r = primitiveConverter[uint]{4, ReadUint, WriteUint, DecodeUint, EncodeUint, bytesToSliceFixed32[uint], sliceToBytesFixed32[uint]}
	case uintptr:
		r = primitiveConverter[uintptr]{4, ReadUintptr, WriteUintptr, DecodeUintptr, EncodeUintptr, bytesToSliceFixed32[uintptr], sliceToBytesFixed32[uintptr]}
	case time.Duration:
		r = primitiveConverter[time.Duration]{8, ReadDuration, WriteDuration, DecodeDuration, EncodeDuration, bytesToSlice[time.Duration], sliceToBytes[time.Duration]}
	default:
		panic(fmt.Sprintf("unsupported primitive type %T", c))
	}

	return r.(TypeConverter[T])
}

func bytesToSlice[T primitive](b []byte, size int) []T {
	if b == nil {
		return nil
	} else if len(b) == 0 {
		return []T{}
	} else {
		return unsafe.Slice((*T)(unsafe.Pointer(&b[0])), len(b)/size)
	}
}

func sliceToBytes[T primitive](s []T, size int) []byte {
	if s == nil {
		return nil
	} else if len(s) == 0 {
		return []byte{}
	} else {
		return unsafe.Slice((*byte)(unsafe.Pointer(&s[0])), len(s)*size)
	}
}

func bytesToSliceFixed32[T ~int | ~uint | ~uintptr](b []byte, size int) []T {
	x := make([]T, len(b)/size)
	for i := 0; i < len(x); i++ {
		x[i] = T(*(*uint32)(unsafe.Pointer(&b[i*size])))
	}
	return x
}

func sliceToBytesFixed32[T ~int | ~uint | ~uintptr](s []T, size int) []byte {
	b := make([]byte, len(s)*size)
	for i := 0; i < len(s); i++ {
		u := uint32(s[i])
		bytes := *(*[4]byte)(unsafe.Pointer(&u))
		copy(b[i*size:], bytes[:size])
	}
	return b
}

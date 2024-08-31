/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import (
	"context"
	"reflect"

	"hypruntime/utils"
)

type TypeHandler interface {
	Info() TypeHandlerInfo
	Read(ctx context.Context, wa WasmAdapter, offset uint32) (any, error)
	Write(ctx context.Context, wa WasmAdapter, offset uint32, obj any) (utils.Cleaner, error)
	Decode(ctx context.Context, wa WasmAdapter, vals []uint64) (any, error)
	Encode(ctx context.Context, wa WasmAdapter, obj any) ([]uint64, utils.Cleaner, error)
}

type TypeHandlerInfo interface {
	TypeName() string
	TypeSize() uint32
	EncodingLength() int
	RuntimeType() reflect.Type
	ZeroValue() any
}

func NewTypeHandlerInfo(typeName string, runtimeType reflect.Type, typeSize uint32, encodingLen int) TypeHandlerInfo {
	zeroValue := reflect.Zero(runtimeType).Interface()
	return &handlerInfo{typeName, runtimeType, zeroValue, typeSize, encodingLen, nil}
}

type handlerInfo struct {
	typeName      string
	runtimeType   reflect.Type
	zeroValue     any
	typeSize      uint32
	encodingLen   int
	innerHandlers []TypeHandler
}

func (h *handlerInfo) TypeName() string {
	return h.typeName
}

func (h *handlerInfo) RuntimeType() reflect.Type {
	return h.runtimeType
}

func (h *handlerInfo) ZeroValue() any {
	return h.zeroValue
}

func (h *handlerInfo) TypeSize() uint32 {
	return h.typeSize
}

func (h *handlerInfo) EncodingLength() int {
	return h.encodingLen
}

func (h *handlerInfo) InnerHandlers() []TypeHandler {
	return h.innerHandlers
}

func (h *handlerInfo) AddHandler(handler TypeHandler) {
	h.innerHandlers = append(h.innerHandlers, handler)
}

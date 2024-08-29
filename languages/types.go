/*
 * Copyright 2024 Hypermode, Inc.
 */

package languages

import (
	"context"
	"errors"

	"hmruntime/plugins/metadata"
	"hmruntime/utils"

	wasm "github.com/tetratelabs/wazero/api"
)

type Language interface {
	Name() string
	TypeInfo() TypeInfo
	NewWasmAdapter(mod wasm.Module) WasmAdapter
}

type language struct {
	name               string
	typeInfo           TypeInfo
	wasmAdapterFactory func(wasm.Module) WasmAdapter
}

func (l *language) Name() string {
	return l.name
}

func (l *language) TypeInfo() TypeInfo {
	return l.typeInfo
}

func (l *language) NewWasmAdapter(mod wasm.Module) WasmAdapter {
	return l.wasmAdapterFactory(mod)
}

type TypeInfo interface {
	GetListSubtype(typ string) string
	GetMapSubtypes(typ string) (string, string)
	GetNameForType(typ string) string
	GetUnderlyingType(typ string) string
	IsListType(typ string) bool
	IsBooleanType(typ string) bool
	IsByteSequenceType(v string) bool
	IsFloatType(typ string) bool
	IsIntegerType(typ string) bool
	IsMapType(typ string) bool
	IsNullable(typ string) bool
	IsSignedIntegerType(typ string) bool
	IsStringType(typ string) bool
	IsTimestampType(typ string) bool
	GetSizeOfType(ctx context.Context, typ string) (uint32, error)
}

type WasmAdapter interface {
	InvokeFunction(ctx context.Context, function *metadata.Function, parameters map[string]any) (any, error)
	EncodeData(ctx context.Context, typ string, data any) ([]uint64, utils.Cleaner, error)
	DecodeData(ctx context.Context, typ string, vals []uint64, pData *any) error
	GetEncodingLength(ctx context.Context, typ string) (int, error)
}

type WasmAdapterWithIndirection interface {
	WasmAdapter
	ReadIndirectResults(ctx context.Context, fn *metadata.Function, offset uint32) (map[string]any, error)
	WriteIndirectResults(ctx context.Context, fn *metadata.Function, offset uint32, results []any) (err error)
}

func GetWasmAdapter(ctx context.Context) (WasmAdapter, error) {
	if wa, ok := ctx.Value(utils.WasmAdapterContextKey).(WasmAdapter); ok {
		return wa, nil
	}
	return nil, errors.New("no WasmAdapter in context")
}

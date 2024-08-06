/*
 * Copyright 2024 Hypermode, Inc.
 */

package languages

import (
	"context"

	"hmruntime/plugins/metadata"

	wasm "github.com/tetratelabs/wazero/api"
)

type Language interface {
	Name() string
	TypeInfo() TypeInfo
	WasmAdapter() WasmAdapter
}

type language struct {
	name        string
	typeInfo    TypeInfo
	wasmAdapter WasmAdapter
}

func (l *language) Name() string {
	return l.name
}

func (l *language) TypeInfo() TypeInfo {
	return l.typeInfo
}

func (l *language) WasmAdapter() WasmAdapter {
	return l.wasmAdapter
}

type TypeInfo interface {
	GetArraySubtype(t string) string
	GetMapSubtypes(t string) (string, string)
	GetNameForType(t string) string
	GetUnderlyingType(t string) string
	IsArrayType(t string) bool
	IsBooleanType(t string) bool
	IsByteSequenceType(t string) bool
	IsFloatType(t string) bool
	IsIntegerType(t string) bool
	IsMapType(t string) bool
	IsNullable(t string) bool
	IsSignedIntegerType(t string) bool
	IsStringType(t string) bool
	IsTimestampType(t string) bool
	SizeOfType(t string) uint32
}

type WasmAdapter interface {
	InvokeFunction(ctx context.Context, mod wasm.Module, function *metadata.Function, parameters map[string]any) (any, error)
	EncodeValue(ctx context.Context, mod wasm.Module, data any) (uint64, error)
	DecodeValue(ctx context.Context, mod wasm.Module, val uint64, data any) error
}

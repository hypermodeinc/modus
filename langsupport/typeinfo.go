/*
 * Copyright 2024 Hypermode, Inc.
 */

package langsupport

import "context"

type TypeInfo interface {
	GetListSubtype(typ string) string
	GetMapSubtypes(typ string) (string, string)
	GetNameForType(typ string) string
	GetSizeOfType(ctx context.Context, typ string) (uint32, error)
	GetAlignOfType(ctx context.Context, typ string) (uint32, error)
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
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"context"
	"reflect"
)

type LanguageTypeInfo interface {
	GetReflectedType(ctx context.Context, typ string) (reflect.Type, error)
	GetSizeOfType(ctx context.Context, typ string) (uint32, error)
	GetAlignmentOfType(ctx context.Context, typ string) (uint32, error)
	GetDataSizeOfType(ctx context.Context, typ string) (uint32, error)
	GetEncodingLengthOfType(ctx context.Context, typ string) (uint32, error)
	ObjectsUseMaxFieldAlignment() bool

	GetListSubtype(typ string) string
	GetMapSubtypes(typ string) (string, string)
	GetNameForType(typ string) string
	GetUnderlyingType(typ string) string

	IsBooleanType(typ string) bool
	IsByteSequenceType(typ string) bool
	IsFloatType(typ string) bool
	IsIntegerType(typ string) bool
	IsListType(typ string) bool
	IsMapType(typ string) bool
	IsObjectType(typ string) bool
	IsNullableType(typ string) bool
	IsPointerType(typ string) bool
	IsPrimitiveType(typ string) bool
	IsSignedIntegerType(typ string) bool
	IsStringType(typ string) bool
	IsTimestampType(typ string) bool
}

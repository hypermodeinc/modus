/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	"github.com/google/uuid"
)

func NilIf[T any](condition bool, val T) *T {
	if condition {
		return nil
	}
	return &val
}

func NilIfEmpty(val string) *string {
	return NilIf(val == "", val)
}

func EnvVarFlagEnabled(envVarName string) bool {
	v := os.Getenv(envVarName)
	b, err := strconv.ParseBool(v)
	return err == nil && b
}

func DebugModeEnabled() bool {
	return EnvVarFlagEnabled("MODUS_DEBUG")
}

func TraceModeEnabled() bool {
	return EnvVarFlagEnabled("MODUS_TRACE")
}

func TrimStringBefore(s string, sep string) string {
	parts := strings.SplitN(s, sep, 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return s
}

func GenerateUUIDv7() string {
	return uuid.Must(uuid.NewV7()).String()
}

func ConvertToError(e any) error {
	if e == nil {
		return nil
	}

	switch e := e.(type) {
	case error:
		return e
	case string:
		return errors.New(e)
	default:
		return fmt.Errorf("%v", e)
	}
}

func getUnsafeDataPtr(x any) unsafe.Pointer {
	type iface struct {
		typ  unsafe.Pointer
		data unsafe.Pointer
	}

	internal := *(*iface)(unsafe.Pointer(&x))
	return internal.data
}

// HasNil returns true if the given interface value is nil, or contains an object that is nil.
func HasNil(x any) bool {
	if x == nil {
		return true
	}

	// this is essentially an optimized version of reflect.ValueOf(x).IsNil()
	// that will never panic and avoids most of the reflection overhead
	p := getUnsafeDataPtr(x)
	switch reflect.TypeOf(x).Kind() {
	case reflect.Interface, reflect.Slice:
		return *(*unsafe.Pointer)(p) == nil
	default:
		return p == nil
	}
}

func CanBeNil(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

func GetStructFieldValue(rs reflect.Value, fieldName string, caseInsensitive bool) (any, error) {
	if rs.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, got %s", rs.Kind())
	}

	rsType := rs.Type()

	var field reflect.StructField
	var found bool
	if caseInsensitive {
		field, found = rsType.FieldByNameFunc(func(s string) bool { return strings.EqualFold(s, fieldName) })
	} else {
		field, found = rsType.FieldByName(fieldName)
	}

	if !found {
		return nil, fmt.Errorf("field %s not found in struct %s", fieldName, rs.Type().Name())
	}

	rf := rs.FieldByIndex(field.Index)

	// allow retrieving values of unexported fields
	// see https://stackoverflow.com/a/43918797
	if !rf.CanInterface() {
		if !rf.CanAddr() {
			rs2 := reflect.New(rsType.Elem())
			rs2.Set(rs)
			rf = rs2.FieldByIndex(field.Index)
		}
		rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	}

	return rf.Interface(), nil
}

/*
 * Copyright 2024 Hypermode, Inc.
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

func HypermodeDebugEnabled() bool {
	return EnvVarFlagEnabled("HYPERMODE_DEBUG")
}

func HypermodeTraceEnabled() bool {
	return EnvVarFlagEnabled("HYPERMODE_TRACE")
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

func GetUnsafeDataPtr(x any) unsafe.Pointer {
	type iface struct {
		typ  unsafe.Pointer
		data unsafe.Pointer
	}

	internal := *(*iface)(unsafe.Pointer(&x))
	return internal.data
}

// HasNil returns true if the given interface value is nil, or contains a nil pointer.
func HasNil(x any) bool {
	return x == nil || uintptr(GetUnsafeDataPtr(x)) == 0
}

func CanBeNil(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

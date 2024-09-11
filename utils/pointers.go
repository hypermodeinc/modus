/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"reflect"
	"time"
)

func DereferencePointer(obj any) any {
	// optimization for common types
	switch t := obj.(type) {
	case *string:
		return *t
	case *bool:
		return *t
	case *int:
		return *t
	case *int8:
		return *t
	case *int16:
		return *t
	case *int32:
		return *t
	case *int64:
		return *t
	case *uint:
		return *t
	case *uint8:
		return *t
	case *uint16:
		return *t
	case *uint32:
		return *t
	case *uint64:
		return *t
	case *uintptr:
		return *t
	case *float32:
		return *t
	case *float64:
		return *t
	case *time.Time:
		return *t
	case *time.Duration:
		return *t
	case string, bool,
		int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64, uintptr,
		float32, float64,
		time.Time, time.Duration:
		return t // already dereferenced
	}

	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr {
		return obj // already dereferenced
	}
	if rv.IsNil() {
		return nil
	}
	return rv.Elem().Interface()
}

func MakePointer(obj any) any {
	// optimization for common types
	switch t := obj.(type) {
	case string:
		return &t
	case bool:
		return &t
	case int:
		return &t
	case int8:
		return &t
	case int16:
		return &t
	case int32:
		return &t
	case int64:
		return &t
	case uint:
		return &t
	case uint8:
		return &t
	case uint16:
		return &t
	case uint32:
		return &t
	case uint64:
		return &t
	case uintptr:
		return &t
	case float32:
		return &t
	case float64:
		return &t
	case time.Time:
		return &t
	case time.Duration:
		return &t
	}

	p := reflect.New(reflect.TypeOf(obj))
	p.Elem().Set(reflect.ValueOf(obj))
	return p.Interface()
}

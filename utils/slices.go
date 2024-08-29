/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"fmt"
	"reflect"
	"time"
	"unsafe"
)

func ConvertToSlice(input any) ([]any, error) {
	// Handle slices of common types directly
	switch input := input.(type) {
	case []any:
		return input, nil
	case []string:
		return convertSlice(input)
	case []uint:
		return convertSlice(input)
	case []uint8:
		return convertSlice(input)
	case []uint16:
		return convertSlice(input)
	case []uint32:
		return convertSlice(input)
	case []uint64:
		return convertSlice(input)
	case []int:
		return convertSlice(input)
	case []int8:
		return convertSlice(input)
	case []int16:
		return convertSlice(input)
	case []int32:
		return convertSlice(input)
	case []int64:
		return convertSlice(input)
	case []float32:
		return convertSlice(input)
	case []float64:
		return convertSlice(input)
	case []bool:
		return convertSlice(input)
	case []uintptr:
		return convertSlice(input)
	case []unsafe.Pointer:
		return convertSlice(input)
	case []time.Time:
		return convertSlice(input)
	case []JSONTime:
		return convertSlice(input)
	}

	// We need to use reflection for the general case.

	rv := reflect.ValueOf(input)
	kind := rv.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil, fmt.Errorf("input is not a slice")
	}

	out := make([]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out = append(out, rv.Index(i).Interface())
	}

	return out, nil
}

func convertSlice[T any](input []T) ([]any, error) {
	out := make([]any, len(input))
	for i, v := range input {
		out[i] = v
	}
	return out, nil
}

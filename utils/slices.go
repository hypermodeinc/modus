/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"fmt"
	"reflect"
	"time"
)

func ConvertToArray(input any) ([]any, error) {
	// Handle slices of common types directly
	switch arr := input.(type) {
	case []any:
		return arr, nil
	case []string:
		return convertSlice(arr)
	case []uint:
		return convertSlice(arr)
	case []uint8:
		return convertSlice(arr)
	case []uint16:
		return convertSlice(arr)
	case []uint32:
		return convertSlice(arr)
	case []uint64:
		return convertSlice(arr)
	case []int:
		return convertSlice(arr)
	case []int8:
		return convertSlice(arr)
	case []int16:
		return convertSlice(arr)
	case []int32:
		return convertSlice(arr)
	case []int64:
		return convertSlice(arr)
	case []float32:
		return convertSlice(arr)
	case []float64:
		return convertSlice(arr)
	case []bool:
		return convertSlice(arr)
	case []time.Time:
		return convertSlice(arr)
	case []JSONTime:
		return convertSlice(arr)
	}

	// Unfortunately, we need to use reflection for the general case

	rv := reflect.ValueOf(input)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("input is not an array")
	}

	out := make([]any, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		out = append(out, rv.Index(i).Interface())
	}

	return out, nil
}

func convertSlice[T any](arr []T) ([]any, error) {
	out := make([]any, len(arr))
	for i, v := range arr {
		out[i] = v
	}
	return out, nil
}

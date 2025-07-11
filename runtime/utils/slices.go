/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
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
	for i := range rv.Len() {
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

func ConvertToSliceOf[T any](obj any) ([]T, bool) {
	switch obj := obj.(type) {
	case nil:
		return nil, true
	case []T:
		return obj, true
	case []any:
		out := make([]T, len(obj))
		for i, v := range obj {
			if t, err := Cast[T](v); err != nil {
				return nil, false
			} else {
				out[i] = t
			}
		}
		return out, true
	case [0]T:
		return obj[:], true
	case [1]T:
		return obj[:], true
	case [2]T:
		return obj[:], true
	case [3]T:
		return obj[:], true
	case [4]T:
		return obj[:], true
	case [5]T:
		return obj[:], true
	case [6]T:
		return obj[:], true
	case [7]T:
		return obj[:], true
	case [8]T:
		return obj[:], true
	case [9]T:
		return obj[:], true
	case [10]T:
		return obj[:], true
	case [11]T:
		return obj[:], true
	case [12]T:
		return obj[:], true
	case [13]T:
		return obj[:], true
	case [14]T:
		return obj[:], true
	case [15]T:
		return obj[:], true
	case [16]T:
		return obj[:], true
	case [17]T:
		return obj[:], true
	case [18]T:
		return obj[:], true
	case [19]T:
		return obj[:], true
	case [20]T:
		return obj[:], true
	case [21]T:
		return obj[:], true
	case [22]T:
		return obj[:], true
	case [23]T:
		return obj[:], true
	case [24]T:
		return obj[:], true
	case [25]T:
		return obj[:], true
	case [26]T:
		return obj[:], true
	case [27]T:
		return obj[:], true
	case [28]T:
		return obj[:], true
	case [29]T:
		return obj[:], true
	case [30]T:
		return obj[:], true
	case [31]T:
		return obj[:], true
	case [32]T:
		return obj[:], true
	}

	rv := reflect.ValueOf(obj)
	switch rv.Kind() {
	case reflect.Array:
		out := make([]T, rv.Len())
		for i := range rv.Len() {
			out[i] = rv.Index(i).Interface().(T)
		}
		return out, true
	case reflect.Slice:
		if rv.IsNil() {
			var out []T = nil
			return out, true
		}
		out := make([]T, rv.Len())
		for i := range rv.Len() {
			out[i] = rv.Index(i).Interface().(T)
		}
		return out, true
	}

	return nil, false
}

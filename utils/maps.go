/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"fmt"
	"reflect"
)

func MapKeys[M ~map[K]V, K comparable, V any](m M) []K {
	keys := make([]K, len(m))

	i := 0
	for k := range m {
		keys[i] = k
		i++
	}

	return keys
}

func MapValues[M ~map[K]V, K comparable, V any](m M) []V {
	vals := make([]V, len(m))

	i := 0
	for _, v := range m {
		vals[i] = v
		i++
	}

	return vals
}

func MapKeysAndValues[M ~map[K]V, K comparable, V any](m M) ([]K, []V) {
	size := len(m)
	keys := make([]K, size)
	vals := make([]V, size)

	i := 0
	for k, v := range m {
		vals[i] = v
		keys[i] = k
		i++
	}

	return keys, vals
}

func ConvertToMap(input any) (map[any]any, error) {
	// Handle maps of common types directly
	switch input := input.(type) {
	case map[any]any:
		return input, nil
	case map[string]any:
		return convertMap(input)
	case map[string]string:
		return convertMap(input)
	}

	// We need to use reflection for the general case.

	rv := reflect.ValueOf(input)
	kind := rv.Kind()
	if kind != reflect.Map {
		return nil, fmt.Errorf("input is not a map")
	}

	out := make(map[any]any, rv.Len())

	iter := rv.MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()
		out[k.Interface()] = v.Interface()
	}

	return out, nil
}

func convertMap[K comparable, V any](input map[K]V) (map[any]any, error) {
	out := make(map[any]any, len(input))
	for k, v := range input {
		out[k] = v
	}
	return out, nil
}

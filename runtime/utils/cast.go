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
	"fmt"

	"github.com/spf13/cast"
)

func Cast[T any](obj any) (T, error) {
	if t, ok := obj.(T); ok {
		return t, nil
	}

	var result T
	switch any(result).(type) {
	case int:
		v, e := cast.ToIntE(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case int8:
		v, e := cast.ToInt8E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case int16:
		v, e := cast.ToInt16E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case int32:
		v, e := cast.ToInt32E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case int64:
		v, e := cast.ToInt64E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case uint:
		v, e := cast.ToUintE(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)

	case uint8:
		v, e := cast.ToUint8E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case uint16:
		v, e := cast.ToUint16E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case uint32:
		v, e := cast.ToUint32E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case uint64:
		v, e := cast.ToUint64E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case float32:
		v, e := cast.ToFloat32E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case float64:
		v, e := cast.ToFloat64E(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case bool:
		v, e := cast.ToBoolE(obj)
		if e != nil {
			return result, e
		}
		result = any(v).(T)
	case uintptr:
		v, e := cast.ToUint32E(obj)
		if e != nil {
			return result, e
		}
		result = any(uintptr(v)).(T)
	default:
		return result, fmt.Errorf("unsupported type: %T", obj)
	}

	return result, nil
}

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package langsupport

import (
	"fmt"
	"reflect"

	"github.com/hypermodeinc/modus/runtime/utils"
)

// AlignOffset returns the smallest y >= x such that y % a == 0.
func AlignOffset(x, a uint32) uint32 {
	return (x + a - 1) &^ (a - 1)
}

// GetFieldObject retrieves the field object with the given name from either a map or a reflect.Value of a struct.
func GetFieldObject(fieldName string, mapObj map[string]any, rvObj reflect.Value) (any, error) {
	if mapObj != nil {
		// case sensitive when reading from map
		if obj, ok := mapObj[fieldName]; !ok {
			return nil, fmt.Errorf("field %s not found in map", fieldName)
		} else {
			return obj, nil
		}
	} else {
		// case insensitive when reading from struct
		if obj, err := utils.GetStructFieldValue(rvObj, fieldName, true); err != nil {
			return nil, err
		} else {
			return obj, nil
		}
	}
}

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package functions

import (
	"fmt"

	"github.com/hypermodeinc/modus/lib/metadata"
)

func CreateParametersMap(fn *metadata.Function, paramValues ...any) (map[string]any, error) {
	if len(paramValues) != len(fn.Parameters) {
		return nil, fmt.Errorf("function %s expects %d parameters, got %d",
			fn.Name,
			len(fn.Parameters),
			len(paramValues))
	}

	parameters := make(map[string]any, len(paramValues))
	for i, value := range paramValues {
		name := fn.Parameters[i].Name
		parameters[name] = value
	}

	return parameters, nil
}

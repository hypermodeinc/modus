/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"fmt"

	"hypruntime/plugins/metadata"
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

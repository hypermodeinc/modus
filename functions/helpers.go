/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"errors"
	"fmt"
	"strings"
)

func CreateParametersMap(info *FunctionInfo, paramValues ...any) (map[string]any, error) {
	if len(paramValues) != len(info.Function.Parameters) {
		return nil, fmt.Errorf("function %s expects %d parameters, got %d",
			info.Function.Name,
			len(info.Function.Parameters),
			len(paramValues))
	}

	parameters := make(map[string]any, len(paramValues))
	for i, value := range paramValues {
		name := info.Function.Parameters[i].Name
		parameters[name] = value
	}

	return parameters, nil
}

func VerifyFunctionSignature(info *FunctionInfo, expectedTypes ...string) error {
	if len(expectedTypes) == 0 {
		return errors.New("expectedTypes must not be empty")
	}
	l := len(expectedTypes)
	expectedSig := fmt.Sprintf("(%s):%s", strings.Join(expectedTypes[:l-1], ","), expectedTypes[l-1])

	sig := info.Function.Signature()
	if sig != expectedSig {
		return fmt.Errorf("function %s has signature %s, expected %s", info.Function.Name, sig, expectedSig)
	}

	return nil
}

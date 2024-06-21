/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"hmruntime/plugins"
	"testing"
)

func Test_VerifyFunctionSignature(t *testing.T) {

	info := FunctionInfo{
		Function: plugins.FunctionSignature{
			Name: "myFunction",
			Parameters: []plugins.Parameter{
				{Name: "param1", Type: plugins.TypeInfo{Name: "int"}},
				{Name: "param2", Type: plugins.TypeInfo{Name: "string"}},
			},
			ReturnType: plugins.TypeInfo{Name: "bool"},
		}}

	err := VerifyFunctionSignature(info, "int", "string", "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}

	info = FunctionInfo{
		Function: plugins.FunctionSignature{
			Name:       "anotherFunction",
			ReturnType: plugins.TypeInfo{Name: "bool"},
		},
	}

	err = VerifyFunctionSignature(info, "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}
}

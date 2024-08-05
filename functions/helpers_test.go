/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"hmruntime/plugins/metadata"
	"testing"
)

func Test_VerifyFunctionSignature(t *testing.T) {

	info := FunctionInfo{
		Function: metadata.Function{
			Name: "myFunction",
			Parameters: []metadata.Parameter{
				{Name: "param1", Type: metadata.TypeInfo{Name: "int"}},
				{Name: "param2", Type: metadata.TypeInfo{Name: "string"}},
			},
			ReturnType: metadata.TypeInfo{Name: "bool"},
		}}

	err := VerifyFunctionSignature(info, "int", "string", "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}

	info = FunctionInfo{
		Function: metadata.Function{
			Name:       "anotherFunction",
			ReturnType: metadata.TypeInfo{Name: "bool"},
		},
	}

	err = VerifyFunctionSignature(info, "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}
}

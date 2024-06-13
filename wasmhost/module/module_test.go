/*
 * Copyright 2024 Hypermode, Inc.
 */

package module

import (
	"hmruntime/functions"
	"hmruntime/plugins"
	"testing"
)

func Test_VerifyFunctionSignature(t *testing.T) {

	functions.Functions["myFunction"] = functions.FunctionInfo{
		Function: plugins.FunctionSignature{
			Name: "myFunction",
			Parameters: []plugins.Parameter{
				{Name: "param1", Type: plugins.TypeInfo{Name: "int"}},
				{Name: "param2", Type: plugins.TypeInfo{Name: "string"}},
			},
			ReturnType: plugins.TypeInfo{Name: "bool"},
		}}

	err := VerifyFunctionSignature("myFunction", "int", "string", "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}

	functions.Functions["anotherFunction"] = functions.FunctionInfo{
		Function: plugins.FunctionSignature{
			Name:       "anotherFunction",
			ReturnType: plugins.TypeInfo{Name: "bool"},
		},
	}

	err = VerifyFunctionSignature("anotherFunction", "bool")
	if err != nil {
		t.Errorf("verifyFunctionSignature failed: %v", err)
	}

	err = VerifyFunctionSignature("nonExistentFunction", "float64")
	if err == nil {
		t.Error("verifyFunctionSignature should have returned an error for non-existent function")
	}
}

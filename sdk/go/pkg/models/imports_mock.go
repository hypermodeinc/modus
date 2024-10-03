//go:build !wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import "github.com/hypermodeAI/functions-go/pkg/testutils"

var LookupModelCallStack = testutils.NewCallStack()
var InvokeModelCallStack = testutils.NewCallStack()

const MockResponseText = "Hello, World!"

func lookupModel(modelName *string) *ModelInfo {
	LookupModelCallStack.Push(modelName)

	return &ModelInfo{
		Name:     *modelName,
		FullName: *modelName + "-mock-model",
	}
}

func invokeModel(modelName *string, input *string) *string {
	InvokeModelCallStack.Push(modelName, input)

	output := `{"response":"` + MockResponseText + `"}`
	return &output
}

//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

var LookupModelCallStack = testutils.NewCallStack()
var InvokeModelCallStack = testutils.NewCallStack()

const MockResponseText = "Hello, World!"

func hostGetModelInfo(modelName *string) *ModelInfo {
	LookupModelCallStack.Push(modelName)

	return &ModelInfo{
		Name:     *modelName,
		FullName: *modelName + "-mock-model",
	}
}

func hostInvokeModel(modelName *string, input *string) *string {
	InvokeModelCallStack.Push(modelName, input)

	output := `{"response":"` + MockResponseText + `"}`
	return &output
}

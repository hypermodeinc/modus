//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import "github.com/hypermodeinc/modus/sdk/go/pkg/testutils"

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

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

// The openai package provides objects that conform to the OpenAI API specification.
// It can be used with OpenAI, as well as other services that conform to the OpenAI API.
package openai

// The usage statistics returned by the OpenAI API.
type Usage struct {

	// The number of completion tokens used.
	CompletionTokens int `json:"completion_tokens"`

	// The number of prompt tokens used.
	PromptTokens int `json:"prompt_tokens"`

	// The total number of tokens used.
	TotalTokens int `json:"total_tokens"`
}

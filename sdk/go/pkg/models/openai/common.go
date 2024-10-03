/*
 * Copyright 2024 Hypermode, Inc.
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

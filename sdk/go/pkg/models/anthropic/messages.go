/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package anthropic

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/utils"
	"github.com/tidwall/sjson"
)

// Provides input and output types that conform to the Anthropic Messages API,
// as described in the [API Reference] docs.
//
// [API Reference]: https://docs.anthropic.com/en/api/messages
type MessagesModel struct {
	messagesModelBase
}

type messagesModelBase = models.ModelBase[MessagesModelInput, MessagesModelOutput]

// The input object for the Anthropic Messages API.
type MessagesModelInput struct {

	// The model that will complete your prompt.
	//
	// Must be the exact string expected by the model provider.
	// For example, "claude-3-5-sonnet-20240620".
	//
	// See [models](https://docs.anthropic.com/en/docs/models-overview) for additional
	// details and options.
	Model string `json:"model"`

	// Input messages.
	//
	// We do not currently support image content blocks, which are available starting with
	// Claude 3 models. This will be added in a future release.
	Messages []*Message `json:"messages"`

	// The maximum number of tokens to generate before stopping.
	//
	// Different models have different maximum values for this parameter. See
	// [models](https://docs.anthropic.com/en/docs/models-overview) for details.
	MaxTokens int `json:"max_tokens,omitempty"`

	// A `Metadata` object describing the request.
	Metadata *Metadata `json:"metadata,omitempty"`

	// Custom text sequences that will cause the model to stop generating.
	StopSequences []string `json:"stop_sequences,omitempty"`

	// System prompt.
	//
	// A system prompt is a way of providing context and instructions to Claude, such as
	// specifying a particular goal or role. See [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
	System string `json:"system,omitempty"`

	// A number between `0.0` and `1.0` that controls the randomness injected into the response.
	//
	// It is recommended to use `temperature` closer to `0.0`
	// for analytical / multiple choice, and closer to `1.0` for creative tasks.
	//
	// Note that even with `temperature` of `0.0`, the results will not be fully deterministic.
	//
	// The default value is 1.0.
	Temperature float64 `json:"temperature"` // should we make it *float64?

	// How the model should use the provided tools.
	//
	// Use either `ToolChoiceAuto`, `ToolChoiceAny`, or `ToolChoiceTool(name: string)`.
	ToolChoice *ToolChoice `json:"tool_choice,omitempty"`

	// Definitions of tools that the model may use.
	//
	// Tools can be used for workflows that include running client-side tools and functions,
	// or more generally whenever you want the model to produce a particular JSON structure of output.
	//
	// See Anthropic's [guide](https://docs.anthropic.com/en/docs/tool-use) for more details.
	Tools []Tool `json:"tools,omitempty"`

	// Only sample from the top K options for each subsequent token.
	//
	// Recommended for advanced use cases only. You usually only need to use `temperature`.
	TopK int `json:"top_k,omitempty"`

	// Use nucleus sampling.
	//
	// You should either alter `temperature` or `top_p`, but not both.
	//
	// Recommended for advanced use cases only. You usually only need to use `temperature`.
	TopP float64 `json:"top_p,omitempty"`
}

type Metadata struct {
	// An external identifier for the user who is associated with the request.
	UserId *string `json:"user_id,omitempty"`
}

type ToolChoice struct {

	// The type of tool to call.
	Type string `json:"type"`

	// The name of the tool to use.
	Name string `json:"name,omitempty"`

	// Whether to disable parallel tool use.
	//
	// Defaults to false. If set to true, the model will output at most one tool use.
	DisableParallelToolUse bool `json:"disable_parallel_tool_use,omitempty"`
}

var (
	// The model will automatically decide whether to use tools.
	ToolChoiceAuto ToolChoice = ToolChoice{Type: "auto"}

	// The model will use any available tools.
	ToolChoiceAny ToolChoice = ToolChoice{Type: "any"}

	// The model will use the specified tool.
	ToolChoiceTool = func(name string) ToolChoice {
		t := ToolChoice{Type: "tool"}
		t.Name = name
		return t
	}
)

// The output object for the Anthropic Messages API.
type MessagesModelOutput struct {

	// Unique object identifier.
	Id string `json:"id"`

	// Object type.
	//
	// For Messages, this is always "message".
	Type string `json:"type"`

	// Conversational role of the generated message.
	//
	// This will always be "assistant".
	Role string `json:"role"`

	// A list of chat completion choices. Can be more than one if n is greater than 1 in the input options.
	Content []ContentBlock `json:"content"`

	// The model that handled the request.
	Model string `json:"model"`

	// The reason that the model stopped.
	StopReason string `json:"stop_reason"`

	// Which custom stop sequence was generated, if any.
	StopSequence *string `json:"stop_sequence,omitempty"`

	// The usage statistics for the request.
	Usage Usage `json:"usage"`
}

type Usage struct {
	// The number of input tokens which were used.
	InputTokens int `json:"input_tokens"`

	// The number of output tokens which were used.
	OutputTokens int `json:"output_tokens"`
}

type Message struct {
	// The role of the author of this message.
	Role string `json:"role"`

	// The content of the message.
	Content Content `json:"content"`
}

type Content interface {
	isContent()
}

// Creates a new user message object.
func NewUserMessage(content Content) *Message {
	return &Message{
		Role:    "user",
		Content: content,
	}
}

type StringContent string

type TextContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type ToolUseContent struct {
	Type  string              `json:"type"`
	Id    string              `json:"id"`
	Name  string              `json:"name"`
	Input utils.RawJsonString `json:"input"`
}

type ToolResultContent struct {
	Type      string  `json:"type"`
	ToolUseId string  `json:"tool_use_id"`
	IsError   *bool   `json:"is_error,omitempty"`
	Content   *string `json:"content,omitempty"`
}

func (s StringContent) isContent()     {}
func (t TextContent) isContent()       {}
func (t ToolUseContent) isContent()    {}
func (t ToolResultContent) isContent() {}

type ContentBlock struct {
	// Either "text" or "tool_use"
	Type string `json:"type"`

	Text *string `json:"text,omitempty"`

	Id *string `json:"id,omitempty"`

	Name *string `json:"name,omitempty"`

	Input *utils.RawJsonString `json:"input,omitempty"`
}

// A tool object that the model may call.
type Tool struct {

	// Name of the tool.
	Name string `json:"name"`

	// [JSON schema](https://json-schema.org/) for this tool's input.
	//
	// This defines the shape of the `input` that your tool accepts and that the model will produce.
	InputSchema utils.RawJsonString `json:"input_schema"`

	// Optional, but strongly-recommended description of the tool.
	Description string `json:"description,omitempty"`
}

// Creates an input object for the Anthropic Messages API.
func (m *MessagesModel) CreateInput(messages ...*Message) (*MessagesModelInput, error) {
	return &MessagesModelInput{
		Model:       m.Info().FullName,
		Messages:    messages,
		Temperature: 1.0,
	}, nil
}

func (mi *MessagesModelInput) MarshalJSON() ([]byte, error) {

	type alias MessagesModelInput
	b, err := utils.JsonSerialize(alias(*mi))
	if err != nil {
		return nil, err
	}

	// omit default temperature
	if mi.Temperature == 1.0 {
		b, err = sjson.DeleteBytes(b, "temperature")
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

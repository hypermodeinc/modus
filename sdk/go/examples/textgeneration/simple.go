/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// In this example, we will generate text and other content using OpenAI chat models.
// See https://platform.openai.com/docs/api-reference/chat/create for more details
// about the options available on the model, which you can set on the input object.

// For more advanced examples see:
// - toolcalling.go for using tools with the model
// - media.go for using audio or image content with the model

// This function generates some text based on the instruction and prompt provided.
func GenerateText(instruction, prompt string) (string, error) {

	// The imported ChatModel type follows the OpenAI Chat completion model input format.
	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return "", err
	}

	// We'll start by creating an input object using the instruction and prompt provided.
	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
		// ... if we wanted to add more messages, we could do so here.
	)
	if err != nil {
		return "", err
	}

	// This is one of many optional parameters available for the OpenAI chat model.
	input.Temperature = 0.7

	// Here we invoke the model with the input we created.
	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	// The output is also specific to the ChatModel interface.
	// Here we return the trimmed content of the first choice.
	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

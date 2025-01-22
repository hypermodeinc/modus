/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/sdk/go/pkg/localtime"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
	"github.com/tidwall/gjson"
)

// This function generates some text from the model, which has already been given some tools to call,
// and instructions on how to use them.  The tools allow the model to answer time-related questions,
// such as "What time is it?", or "What time is it in New York?".
func GenerateTextWithTools(prompt string) (string, error) {

	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return "", err
	}

	instruction := `
	You are a helpful assistant that understands time in various parts of the world.
	Answer the user's question as directly as possible. If you need more information, ask for it.
	When stating a time zone to the user, either use a descriptive name such as "Pacific Time" or use the location's name.
	If asked for a time in a location, and the given location has multiple time zones, return the time in each of them.
	Make your answers as brief as possible.
	`

	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return "", err
	}

	// Let's turn the temperature down a bit, to make the model's responses more predictable.
	input.Temperature = 0.2

	// This is how you make tools available to the model.
	// Each tool is specified by the name of the function to call, and a description of what it does.
	//
	// If the function requires parameters, you can specify simple parameters individually with the WithParameter method,
	// or you can pass a JSON Schema string directly to the WithParametersSchema method.
	//
	// NOTE: A future release of Modus may further simplify this process.
	input.Tools = []openai.Tool{
		openai.NewToolForFunction("getUserTimeZone", "Returns the user's current IANA time zone."),
		openai.NewToolForFunction("getCurrentTimeInUserTimeZone", "Gets the current date and time in the user's time zone.  Returns the user's IANA time zone, and the time as an ISO 8601 string."),
		openai.NewToolForFunction("getCurrentTime", "Returns the current date and time in the specified IANA time zone as an ISO 8601 string.").
			WithParameter("tz", "string", `An IANA time zone identifier such as "America/New_York".`),
	}

	// To use the tools, you'll need to have a "conversation loop" like this one.
	// The model may be called multiple times, depending on the complexity of the conversation and the tools used.
	for {
		output, err := model.Invoke(input)
		if err != nil {
			return "", err
		}

		msg := output.Choices[0].Message

		// If the model requested tools to be called, we'll need to do that before continuing.
		if len(msg.ToolCalls) > 0 {

			// First, add the model's response to the conversation as an "Assistant Message".
			input.Messages = append(input.Messages, msg.ToAssistantMessage())

			// Then call the tools the model requested.  Add each result to the conversation as a "Tool Message".
			for _, tc := range msg.ToolCalls {

				// The model will tell us which tool to call, and what arguments to pass to it.
				// We'll need to parse the function name and arguments in order to call the appropriate function.
				//
				// If the function errors (for example, due to invalid input), we can either continue
				// the conversation by providing the model with the error message (as shown here),
				// or we can end the conversation by returning the error directly.
				//
				// NOTE: A future release of Modus may simplify this process.
				var toolMsg *openai.ToolMessage[string]
				switch tc.Function.Name {

				case "getCurrentTime":
					tz := gjson.Get(tc.Function.Arguments, "tz").Str
					if result, err := getCurrentTime(tz); err == nil {
						toolMsg = openai.NewToolMessage(result, tc.Id)
					} else {
						toolMsg = openai.NewToolMessage(err, tc.Id)
					}

				case "getUserTimeZone":
					timeZone := getUserTimeZone()
					toolMsg = openai.NewToolMessage(timeZone, tc.Id)

				case "getCurrentTimeInUserTimeZone":
					if result, err := getCurrentTimeInUserTimeZone(); err == nil {
						toolMsg = openai.NewToolMessage(result, tc.Id)
					} else {
						toolMsg = openai.NewToolMessage(err, tc.Id)
					}

				default:
					return "", fmt.Errorf("Unknown tool call: %s", tc.Function.Name)
				}

				// Add the tool's response to the conversation.
				input.Messages = append(input.Messages, toolMsg)
			}

		} else if msg.Content != "" {
			// return the model's final response to the user.
			return strings.TrimSpace(msg.Content), nil
		} else {
			// If the model didn't ask for tools, and didn't return any content, something went wrong.
			return "", errors.New("Invalid response from model.")
		}
	}
}

// The following functions are made available as "tools"
// for the model to call in the example above.

// This function will return the current time in a given time zone.
func getCurrentTime(tz string) (string, error) {
	if !localtime.IsValidTimeZone(tz) {
		return "", errors.New("Invalid time zone.")
	}

	now, err := localtime.NowInZone(tz)
	if err != nil {
		return "", err
	}
	return now.Format(time.RFC3339), nil
}

// This function will return the default time zone for the user.
func getUserTimeZone() string {
	tz := os.Getenv("TZ")
	if tz == "" {
		return "(unknown)"
	}
	return tz
}

// This function combines the functionality of getCurrentTime and getUserTimeZone,
// to reduce the number of tool calls required by the model.
func getCurrentTimeInUserTimeZone() (*struct{ Time, Zone string }, error) {
	tz := os.Getenv("TZ")
	if tz == "" {
		return nil, errors.New("Cannot determine user's time zone.")
	}
	time, err := getCurrentTime(tz)
	if err != nil {
		return nil, err
	}
	return &struct{ Time, Zone string }{Time: time, Zone: tz}, nil
}

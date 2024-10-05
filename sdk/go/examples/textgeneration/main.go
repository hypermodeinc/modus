/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// In this example, we will generate text using the OpenAI Chat model.
// See https://platform.openai.com/docs/api-reference/chat/create for more details
// about the options available on the model, which you can set on the input object.

// This model name should match the one defined in the hypermode.json manifest file.
const modelName = "text-generator"

// This function generates some text based on the instruction and prompt provided.
func GenerateText(instruction, prompt string) (string, error) {

	// The imported ChatModel type follows the OpenAI Chat completion model input format.
	model, err := models.GetModel[openai.ChatModel](modelName)
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

	// This is one of many optional parameters available for the OpenAI Chat model.
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

// This function generates a single product.
func GenerateProduct(category string) (*Product, error) {

	// We can get creative with the instruction and prompt to guide the model
	// in generating the desired output.  Here we provide a sample JSON of the
	// object we want the model to generate.
	instruction := "Generate a product for the category provided.\n" +
		"Only respond with valid JSON object in this format:\n" + sampleProductJson
	prompt := fmt.Sprintf(`The category is "%s".`, category)

	// Set up the input for the model, creating messages for the instruction and prompt.
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return nil, err
	}
	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return nil, err
	}

	// Let's increase the temperature to get more creative responses.
	// Be careful though, if the temperature is too high, the model may generate invalid JSON.
	input.Temperature = 1.2

	// This model also has a response format parameter that can be set to JSON,
	// Which, along with the instruction, can help guide the model in generating valid JSON output.
	input.ResponseFormat = openai.ResponseFormatJson

	// Here we invoke the model with the input we created.
	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	// The output should contain the JSON string we asked for.
	content := strings.TrimSpace(output.Choices[0].Message.Content)

	// We can now parse the JSON string as a Product object.
	var product Product
	if err := json.Unmarshal([]byte(content), &product); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &product, nil
}

// This function generates multiple product.
func GenerateProducts(category string, quantity int) ([]Product, error) {

	// Similar to the previous example above, we can tailor the instruction and prompt
	// to guide the model in generating the desired output.  Note that understanding the behavior
	// of the model is important to get the desired results.  In this case, we need the model
	// to return an _object_ containing an array, not an array of objects directly.
	// That's because the model will not reliably generate an array of objects directly.
	instruction := fmt.Sprintf("Generate %d products for the category provided.\n"+
		"Only respond with a valid JSON object containing a valid JSON array named 'list', in this format:\n"+
		`{"list":[%s]}`, quantity, sampleProductJson)
	prompt := fmt.Sprintf(`The category is "%s".`, category)

	// Set up the input for the model, creating messages for the instruction and prompt.
	model, err := models.GetModel[openai.ChatModel](modelName)
	if err != nil {
		return nil, err
	}
	input, err := model.CreateInput(
		openai.NewSystemMessage(instruction),
		openai.NewUserMessage(prompt),
	)
	if err != nil {
		return nil, err
	}

	// Adjust the model inputs, just like in the previous example.
	// Be careful, if the temperature is too high, the model may generate invalid JSON.
	input.Temperature = 1.2
	input.ResponseFormat = openai.ResponseFormatJson

	// Here we invoke the model with the input we created.
	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	// The output should contain the JSON string we asked for.
	content := strings.TrimSpace(output.Choices[0].Message.Content)

	// We can parse that JSON to a compatible object, to get the data we're looking for.
	var data map[string][]Product
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Now we can extract the list of products from the data.
	products, found := data["list"]
	if !found {
		return nil, fmt.Errorf("expected 'list' key in JSON object")
	}
	return products, nil
}

/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// In this example, we will generate text using the OpenAI Chat model.
// See https://platform.openai.com/docs/api-reference/chat/create for more details
// about the options available on the model, which you can set on the input object.

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
	model, err := models.GetModel[openai.ChatModel]("text-generator")
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
	model, err := models.GetModel[openai.ChatModel]("text-generator")
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

// This function generates text that describes the image at the provided url.
func DescribeImage(url string) (string, error) {

	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return "", err
	}

	input, err := model.CreateInput(
		openai.NewUserMessage([]openai.ContentPart{
			openai.NewContentPartText("Describe this image."),
			openai.NewContentPartImageFromUrl(url),
		}),
	)
	if err != nil {
		return "", err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(output.Choices[0].Message.Content), nil
}

// This type is used for image or audio data.
type Media struct {

	// The content type of the media.
	ContentType string

	// The binary data of the media.
	Data []byte

	// A text description or transcription of the media.
	Text string
}

// This function fetches a random image, and then generates text that describes it.
func DescribeRandomImage(width, height int) (*Media, error) {

	// Fetch a random image from the Picsum API.
	url := fmt.Sprintf("https://picsum.photos/%d/%d", width, height)
	response, err := http.Fetch(url)
	if err != nil {
		return nil, err
	}
	data := response.Body
	contentType := *response.Headers.Get("Content-Type")

	// Describe the image using the OpenAI Chat model.
	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(
		openai.NewUserMessage([]openai.ContentPart{
			openai.NewContentPartText("Describe this image."),
			openai.NewContentPartImage(data, contentType),
		}),
	)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	text := strings.TrimSpace(output.Choices[0].Message.Content)

	// Return the image and its description.
	image := &Media{
		ContentType: contentType,
		Data:        data,
		Text:        text,
	}

	return image, nil
}

// This function fetches a random "Harvard Sentences" speech file from OpenSpeech, and then generates a transcript from it.
// The sentences are from https://www.cs.columbia.edu/~hgs/audio/harvard.html
func TranscribeRandomSpeech() (*Media, error) {

	// Pick a random file number from the list of available here:
	// https://www.voiptroubleshooter.com/open_speech/american.html
	numbers := []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 34, 35, 36, 37, 38, 39, 40, 57, 58, 59, 60, 61}
	num := numbers[rand.IntN(len(numbers))]

	// Fetch the speech file corresponding to the number.
	url := fmt.Sprintf("https://www.voiptroubleshooter.com/open_speech/american/OSR_us_000_%04d_8k.wav", num)
	response, err := http.Fetch(url)
	if err != nil {
		return nil, err
	}
	data := response.Body

	// Transcribe the audio using the OpenAI Chat model.
	model, err := models.GetModel[openai.ChatModel]("audio-transcriber")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(
		openai.NewDeveloperMessage("Do not include any newlines or surrounding quotation marks in the response. Omit any explanation beyond the request."),
		openai.NewUserMessage([]openai.ContentPart{
			openai.NewContentPartText("Provide an exact transcription of the contents of this audio file."),
			openai.NewContentPartAudio(data, "wav"),
		}),
	)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	text := strings.TrimSpace(output.Choices[0].Message.Content)

	// Return the audio file and its transcript.
	image := &Media{
		ContentType: "audio/wav",
		Data:        data,
		Text:        text,
	}

	return image, nil
}

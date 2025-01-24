/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/hypermodeinc/modus/sdk/go/pkg/http"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
	"github.com/hypermodeinc/modus/sdk/go/pkg/models/openai"
)

// These examples demonstrate how to use audio or image data with OpenAI chat models.
// Currently, audio can be used for input or output, but images can be used only for input.

// This type is used in these examples to represent images or audio.
type Media struct {

	// The content type of the media.
	ContentType string

	// The binary data of the media.
	// This value will be base64 encoded when used in an API response.
	Data []byte

	// A text description or transcription of the media.
	Text string
}

// This function generates an audio response based on the instruction and prompt provided.
func GenerateAudio(instruction, prompt string) (*Media, error) {

	// Note, this is similar to GenerateText above, but with audio output requested.

	// We'll generate the audio using an audio-enabled OpenAI chat model.
	model, err := models.GetModel[openai.ChatModel]("audio-model")
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

	input.Temperature = 0.7

	// Request audio output from the model.
	// Note, this is a convenience method that requests audio modality and sets the voice and format.
	// You can also set these values manually on the input object, if you prefer.
	input.RequestAudioOutput("ash", "wav")

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	// Return the audio and its transcription.
	// Note that the message Content field will be empty for audio responses.
	// Instead, the text will be in the Message.Audio.Transcript field.
	audio := output.Choices[0].Message.Audio

	media := &Media{
		ContentType: "audio/wav",
		Data:        audio.Data,
		Text:        strings.TrimSpace(audio.Transcript),
	}

	return media, nil
}

// This function generates text that describes the image at the provided url.
// In this example the image url is passed to the model, and the model retrieves the image.
func DescribeImage(url string) (string, error) {

	// Note that because the model retrieves the image, any URL can be used.
	// However, this means that there is a risk of sending data to an unauthorized host, if the URL is not hardcoded or sanitized.
	// See the DescribeRandomImage function below for a safer approach.

	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return "", err
	}

	input, err := model.CreateInput(
		openai.NewUserMessageFromParts(
			openai.NewTextContentPart("Describe this image."),
			openai.NewImageContentPartFromUrl(url),
		),
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

// This function fetches a random image, and then generates text that describes it.
// In this example the image is retrieved by the function before passing it as data to the model.
func DescribeRandomImage() (*Media, error) {

	// Because this approach fetches the image directly, it is safer than the DescribeImage function above.
	// The host URL is allow-listed in the modus.json file, so we can trust the image source.

	// Fetch a random image from the Picsum API.  We'll just hardcode the size to make the demo simple to call.
	response, err := http.Fetch("https://picsum.photos/640/480")
	if err != nil {
		return nil, err
	}
	data := response.Body
	contentType := *response.Headers.Get("Content-Type")

	// Describe the image using the OpenAI chat model.
	model, err := models.GetModel[openai.ChatModel]("text-generator")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(
		openai.NewUserMessageFromParts(
			openai.NewTextContentPart("Describe this image."),
			openai.NewImageContentPartFromData(data, contentType),
		),
	)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	// Return the image and its generated description.
	text := strings.TrimSpace(output.Choices[0].Message.Content)
	media := &Media{
		ContentType: contentType,
		Data:        data,
		Text:        text,
	}

	return media, nil
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

	// Transcribe the audio using an audio-enabled OpenAI chat model.
	model, err := models.GetModel[openai.ChatModel]("audio-model")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(
		openai.NewDeveloperMessage("Do not include any newlines or surrounding quotation marks in the response. Omit any explanation beyond the request."),
		openai.NewUserMessageFromParts(
			openai.NewTextContentPart("Provide an exact transcription of the contents of this audio file."),
			openai.NewAudioContentPartFromData(data, "wav"),
		),
	)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	// Return the audio file and its transcript.
	text := strings.TrimSpace(output.Choices[0].Message.Content)
	media := &Media{
		ContentType: "audio/wav",
		Data:        data,
		Text:        text,
	}

	return media, nil
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package openai

import (
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/models"
)

// Provides input and output types that conform to the OpenAI Embeddings API,
// as described in the [API Reference] docs.
//
// [API Reference]: https://platform.openai.com/docs/api-reference/embeddings
type EmbeddingsModel struct {
	embeddingsModelBase
}

type embeddingsModelBase = models.ModelBase[EmbeddingsModelInput, EmbeddingsModelOutput]

// The input object for the OpenAI Embeddings API.
type EmbeddingsModelInput struct {

	// The name of the model to use to generate the embeddings.
	//
	// Must be the exact string expected by the model provider.
	// For example, "text-embedding-3-small".
	Model string `json:"model"`

	// The input content to vectorize.
	Input any `json:"input"`

	// The format for the output embeddings.
	// The default ("") is equivalent to [EncodingFormatFloat], which is currently the only supported format.
	EncodingFormat EncodingFormat `json:"encoding_format,omitempty"`

	// The maximum number of dimensions for the output embeddings.
	// The default (0) indicates that the model's default number of dimensions will be used.
	Dimensions int `json:"dimensions,omitempty"`

	// The user ID to associate with the request, as described in the [documentation].
	// If not specified, the request will be anonymous.
	//
	// [documentation]: https://platform.openai.com/docs/guides/safety-best-practices/end-user-ids
	User string `json:"user,omitempty"`
}

// The output object for the OpenAI Embeddings API.
type EmbeddingsModelOutput struct {

	// The name of the output object type returned by the API.
	// This will always be "list".
	Object string `json:"object"`

	// The name of the model used to generate the embeddings.
	// In most cases, this will match the requested model field in the input.
	Model string `json:"model"`

	// The usage statistics for the request.
	Usage Usage `json:"usage"`

	// The output vector embeddings data.
	Data []Embedding `json:"data"`
}

// The encoding format for the output embeddings.
type EncodingFormat string

const (
	// The output embeddings are encoded as an array of floating-point numbers.
	EncodingFormatFloat EncodingFormat = "float"

	// The output embeddings are encoded as a base64-encoded string,
	// containing an binary representation of an array of floating-point numbers.
	//
	// NOTE: This format is not currently supported.
	EncodingFormatBase64 EncodingFormat = "base64"
)

// The output vector embeddings data.
type Embedding struct {

	// The name of the output object type returned by the API.
	// This will always be "embedding".
	Object string `json:"object"`

	// The index of the input text that corresponds to this embedding.
	// Used when requesting embeddings for multiple texts.
	Index int `json:"index"`

	// The vector embedding of the input text.
	Embedding []float32 `json:"embedding"`
}

// Creates an input object for the OpenAI Embeddings API.
//
// The content parameter can be any of:
//   - A string representing the text to vectorize.
//   - A slice of strings representing multiple texts to vectorize.
//   - A slice of integers representing pre-tokenized text to vectorize.
//   - A slice of slices of integers representing multiple pre-tokenized texts to vectorize.
//
// NOTE: The input content must not exceed the maximum token limit of the model.
func (m *EmbeddingsModel) CreateInput(content any) (*EmbeddingsModelInput, error) {

	switch content.(type) {
	case string, []string,
		[]uint, []uint8, []uint16, []uint32, []uint64,
		[]int, []int8, []int16, []int32, []int64,
		[][]uint, [][]uint8, [][]uint16, [][]uint32, [][]uint64,
		[][]int, [][]int8, [][]int16, [][]int32, [][]int64:

		return &EmbeddingsModelInput{
			Model: m.Info().FullName,
			Input: content,
		}, nil
	}

	return nil, fmt.Errorf("invalid embedding content type: %T", content)
}

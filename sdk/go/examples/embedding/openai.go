/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"github.com/hypermodeAI/functions-go/pkg/models"
	"github.com/hypermodeAI/functions-go/pkg/models/openai"
)

func GetEmbeddingsForTextWithOpenAI(text string) ([]float32, error) {
	results, err := GetEmbeddingsForTextsWithOpenAI(text)
	if err != nil {
		return nil, err
	}

	return results[0], nil
}

func GetEmbeddingsForTextsWithOpenAI(texts ...string) ([][]float32, error) {
	model, err := models.GetModel[openai.EmbeddingsModel]("openai-embeddings")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(texts)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	results := make([][]float32, len(output.Data))
	for i, d := range output.Data {
		results[i] = d.Embedding
	}

	return results, nil
}

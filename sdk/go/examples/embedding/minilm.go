package main

import (
	"github.com/hypermodeAI/functions-go/pkg/models"
	"github.com/hypermodeAI/functions-go/pkg/models/experimental"
)

func GetEmbeddingsForTextWithMiniLM(text string) ([]float32, error) {
	results, err := GetEmbeddingsForTextsWithMiniLM(text)
	if err != nil {
		return nil, err
	}

	return results[0], nil
}

func GetEmbeddingsForTextsWithMiniLM(texts ...string) ([][]float32, error) {
	model, err := models.GetModel[experimental.EmbeddingsModel]("minilm")
	if err != nil {
		return nil, err
	}

	input, err := model.CreateInput(texts...)
	if err != nil {
		return nil, err
	}

	output, err := model.Invoke(input)
	if err != nil {
		return nil, err
	}

	return output.Predictions, nil
}

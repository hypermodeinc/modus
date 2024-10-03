/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"errors"
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/models"
	"github.com/hypermodeAI/functions-go/pkg/models/experimental"
)

const modelName = "my-classifier"

// This function takes input text and a probability threshold, and returns the
// classification label determined by the model, if the confidence is above the
// threshold. Otherwise, it returns an empty string.
func ClassifyText(text string, threshold float32) (string, error) {
	predictions, err := classify(text)
	if err != nil {
		return "", err
	}

	prediction := predictions[0]
	if prediction.Confidence < threshold {
		return "", nil
	}

	return prediction.Label, nil
}

// This function takes input text and returns the classification labels and their
// corresponding probabilities, as determined by the model.
func GetClassificationLabels(text string) (map[string]float32, error) {
	predictions, err := classify(text)
	if err != nil {
		return nil, err
	}

	return getLabels(predictions[0]), nil
}

func getLabels(prediction experimental.ClassifierResult) map[string]float32 {
	labels := make(map[string]float32, len(prediction.Probabilities))
	for _, p := range prediction.Probabilities {
		labels[p.Label] = p.Probability
	}
	return labels
}

// This function is similar to the previous, but allows multiple items to be classified at a time.
func GetMultipleClassificationLabels(ids []string, texts []string) (map[string]map[string]float32, error) {
	if len(ids) != len(texts) {
		return nil, errors.New("number of IDs does not match number of texts")
	}

	predictions, err := classify(texts...)
	if err != nil {
		return nil, err
	}

	labels := make(map[string]map[string]float32, len(predictions))
	for i, prediction := range predictions {
		labels[ids[i]] = getLabels(prediction)
	}

	return labels, nil
}

func classify(texts ...string) ([]experimental.ClassifierResult, error) {
	model, err := models.GetModel[experimental.ClassificationModel](modelName)
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

	if len(output.Predictions) != len(texts) {
		word := "prediction"
		if len(texts) > 1 {
			word += "s"
		}
		return nil, fmt.Errorf("expected %d %s, got %d", len(texts), word, len(output.Predictions))
	}

	return output.Predictions, nil
}

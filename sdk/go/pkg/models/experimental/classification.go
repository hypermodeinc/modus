/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package experimental

import (
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/models"
)

// A model that returns classification results for a list of text strings.
//
// NOTE: This model interface is experimental and may change in the future.
// It is primarily intended for use with with classification models hosted on Hypermode.
type ClassificationModel struct {
	classificationModelBase
}

type classificationModelBase = models.ModelBase[ClassificationModelInput, ClassificationModelOutput]

// The input object for the classification model.
type ClassificationModelInput struct {

	// A list of one or more text strings to classify.
	Instances []string `json:"instances"`
}

// The output object for the classification model.
type ClassificationModelOutput struct {

	// A list of prediction results that correspond to each input text string.
	Predictions []ClassifierResult `json:"predictions"`
}

// A classification result for a single text string.
type ClassifierResult struct {

	// The classification label with the highest confidence.
	Label string `json:"label"`

	// The confidence score for the classification label.
	Confidence float32 `json:"confidence"`

	// The list of all classification labels with their corresponding probabilities.
	Probabilities []ClassifierLabel `json:"probabilities"`
}

// A classification label with its corresponding probability.
type ClassifierLabel struct {

	// The classification label.
	Label string `json:"label"`

	// The probability value.
	Probability float32 `json:"probability"`
}

// Creates an input object for the classification model.
//
// The content parameter is a list of one or more text strings to classify.
func (m *ClassificationModel) CreateInput(content ...string) (*ClassificationModelInput, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("at least one text string must be provided")
	}

	return &ClassificationModelInput{
		Instances: content,
	}, nil
}

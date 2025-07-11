/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package experimental

import (
	"fmt"

	"github.com/hypermodeinc/modus/sdk/go/pkg/models"
)

// A model that returns embeddings for a list of text strings.
//
// NOTE: This model interface is experimental and may change in the future.
// It is primarily intended for use with with embedding models hosted on Hypermode.
type EmbeddingsModel struct {
	embeddingsModelBase
}

type embeddingsModelBase = models.ModelBase[EmbeddingsModelInput, EmbeddingsModelOutput]

// The input object for the embeddings model.
type EmbeddingsModelInput struct {

	// A list of one or more text strings to create vector embeddings for.
	Instances []string `json:"instances"`
}

// The output object for the embeddings model.
type EmbeddingsModelOutput struct {

	// A list of vector embeddings that correspond to each input text string.
	Predictions [][]float32 `json:"predictions"`
}

// Creates an input object for the embeddings model.
//
// The content parameter is a list of one or more text strings to create vector embeddings for.
func (m *EmbeddingsModel) CreateInput(content ...string) (*EmbeddingsModelInput, error) {
	if len(content) == 0 {
		return nil, fmt.Errorf("at least one text string must be provided")
	}

	return &EmbeddingsModelInput{
		Instances: content,
	}, nil
}

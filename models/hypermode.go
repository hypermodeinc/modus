/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"
	"hmruntime/manifest"
)

type hypermode struct {
	model manifest.Model
}

func (inf *hypermode) InvokeClassifier(ctx context.Context, input []string) (map[string]float64, error) {
	return nil, fmt.Errorf("invokeClassifier not implemented for Hypermode model")
}

func (inf *hypermode) ComputeEmbedding(ctx context.Context, sentenceMap map[string]string) (map[string][]float64, error) {

	result, err := PostToModelEndpoint[[]float64](ctx, sentenceMap, inf.model)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (inf *hypermode) ChatCompletion(ctx context.Context, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error) {

	return ChatResponse{}, fmt.Errorf("ChatCompletion not implemented for Hypermode models")

}

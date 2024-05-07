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

func (inf *hypermode) InvokeClassifier(ctx context.Context, sentenceMap map[string]string) (map[string]map[string]float32, error) {
	if inf.model.Task != manifest.ClassificationTask {
		return nil, fmt.Errorf("InvokeClassifier not supported for model %s with task %s", inf.model.Name, inf.model.Task)
	}
	result, err := PostToModelEndpoint[classifierResult](ctx, sentenceMap, inf.model)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("Empty result returned from model.")
	}

	resultMap := make(map[string]map[string]float32)
	for k, v := range result {
		resultMap[k] = make(map[string]float32)
		for _, p := range v.Probabilities {
			resultMap[k][p.Label] = p.Probability
		}
	}
	return resultMap, nil

}

func (inf *hypermode) ComputeEmbedding(ctx context.Context, sentenceMap map[string]string) (map[string][]float64, error) {
	if inf.model.Task != manifest.EmbeddingTask {
		return nil, fmt.Errorf("ComputeEmbedding not supported for model %s with task %s", inf.model.Name, inf.model.Task)
	}
	result, err := PostToModelEndpoint[[]float64](ctx, sentenceMap, inf.model)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (inf *hypermode) ChatCompletion(ctx context.Context, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error) {

	return ChatResponse{}, fmt.Errorf("ChatCompletion not implemented for Hypermode models")

}

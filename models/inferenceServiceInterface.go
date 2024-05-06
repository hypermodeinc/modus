package models

import (
	"context"
)

type ChatContext struct {
	Model          string         `json:"model"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Messages       []ChatMessage  `json:"messages"`
}

type ResponseFormat struct {
	Type string `json:"type"`
}

type EmbeddingRequest struct {
	Model          string   `json:"model"`
	Input          []string `json:"input"`
	EncodingFormat string   `json:"encoding_format"`
}
type EmbeddingResponse struct {
	Id     string          `json:"id"`
	Object string          `json:"object"`
	Data   []EmbeddingData `json:"data"`
}
type EmbeddingData struct {
	Index     int       `json:"index"`
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
}

type inferenceService interface {
	ChatCompletion(ctx context.Context, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error)
	ComputeEmbedding(ctx context.Context, sentenceMap map[string]string) (map[string][]float64, error)
	InvokeClassifier(ctx context.Context, input []string) (map[string]float64, error)
}

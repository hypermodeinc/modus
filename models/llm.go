package models

import (
	"context"
	"hmruntime/manifest"
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

type llmService interface {
	ChatCompletion(ctx context.Context, model manifest.Model, host manifest.Host, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error)
	Embedding(ctx context.Context, sentenceMap map[string]string, model manifest.Model, host manifest.Host) (map[string][]float64, error)
}

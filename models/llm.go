package models

import (
	"context"
	"hmruntime/appdata"
)

type ChatContext struct {
	Model          string         `json:"model"`
	ResponseFormat ResponseFormat `json:"response_format"`
	Messages       []ChatMessage  `json:"messages"`
}
type ResponseFormat struct {
	Type string `json:"type"`
}

type llmService interface {
	ChatCompletion(ctx context.Context, model appdata.Model, instruction string, sentence string, outputFormat OutputFormat) (ChatResponse, error)
}

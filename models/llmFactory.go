package models

import (
	"fmt"
)

const OpenAIHost string = "openai"
const MistralHost string = "mistral"

// Factory method return the LLM service given a host name
func CreateLlmService(host string) (llmService, error) {
	switch host {
	case OpenAIHost:
		return &openai{}, nil
	case MistralHost:
		return &mistral{}, nil

	default:
		return nil, fmt.Errorf("LLM model host  '%s' not in supported list ['%s','%s'].", host, OpenAIHost, MistralHost)
	}

}

package models

import (
	"fmt"
	"hmruntime/hosts"
	"hmruntime/manifest"
)

// Factory method return the LLM service given a host name
func CreateInferenceService(modelName string) (inferenceService, error) {
	model, err := GetModel(modelName, manifest.GenerationTask)
	if err != nil {
		return nil, err
	}

	// hypermode host is for huggingface models
	if hosts.HypermodeHost == model.Host {
		return &hypermode{model: model}, nil
	} else {
		var host manifest.Host
		if model.Host != hosts.HypermodeHost {
			host, err = hosts.GetHost(model.Host)
			if err != nil {
				return nil, err
			}
		}
		// for non hypermode hosts, check the provider name
		switch model.Provider {
		case OpenAIProvider:
			return &openai{model: model, host: host}, nil
		case MistralProvider:
			return &mistral{model: model, host: host}, nil

		default:
			return nil, fmt.Errorf("Provider '%s' not in supported list ['%s','%s'].", host, OpenAIProvider, MistralProvider)
		}
	}

}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"encoding/json"
	"fmt"

	"hmruntime/functions/assemblyscript"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/manifest"
	"hmruntime/models"
	"hmruntime/models/openai"

	wasm "github.com/tetratelabs/wazero/api"
)

func hostInvokeTextGenerator(ctx context.Context, mod wasm.Module, pModelName uint32, pInstruction uint32, pSentence uint32, pFormat uint32) uint32 {

	var modelName, instruction, sentence, format string
	err := readParams4(ctx, mod, pModelName, pInstruction, pSentence, pFormat, &modelName, &instruction, &sentence, &format)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}
	outputFormat := models.OutputFormat(format)

	model, err := models.GetModel(modelName, manifest.GenerationTask)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	host, err := hosts.GetHost(model.Host)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model host.")
		return 0
	}

	if models.OutputFormatText != outputFormat && models.OutputFormatJson != outputFormat {
		logger.Err(ctx, err).Msg("Unsupported output format.")
		return 0
	}

	var result models.ChatResponse
	switch model.Host {
	case hosts.OpenAIHost:
		result, err = openai.ChatCompletion(ctx, model, host, instruction, sentence, outputFormat)
		if err != nil {
			logger.Err(ctx, err).Msg("Error posting to OpenAI.")
			return 0
		}
		if result.Error.Message != "" {
			err := fmt.Errorf(result.Error.Message)
			logger.Err(ctx, err).Msg("Error returned from OpenAI.")
			return 0
		}
	default:
		err := fmt.Errorf("unsupported model host: %s", model.Host)
		logger.Err(ctx, err).Msg("Unsupported model host.")
		return 0
	}

	if models.OutputFormatJson == outputFormat {
		// safeguard: test is the output is a valid json
		// test every Choices.Message.Content
		for _, choice := range result.Choices {
			_, err := json.Marshal(choice.Message.Content)
			if err != nil {
				logger.Err(ctx, err).Msg("One of the generated message is not a valid JSON.")
				return 0
			}
		}
	}

	// return the first chat response
	if len(result.Choices) == 0 {
		logger.Err(ctx, err).Msg("Empty result returned from OpenAI.")
		return 0
	}
	firstMsgContent := result.Choices[0].Message.Content

	resBytes, err := json.Marshal(firstMsgContent)
	if err != nil {
		logger.Err(ctx, err).Msg("Error marshalling result.")
		return 0
	}

	offset, err := assemblyscript.WriteString(ctx, mod, string(resBytes))
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

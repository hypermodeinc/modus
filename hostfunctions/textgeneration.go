/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/functions/assemblyscript"
	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/models/openai"
	"hmruntime/utils"

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
	if models.OutputFormatText != outputFormat && models.OutputFormatJson != outputFormat {
		logger.Error(ctx).Msg("Unsupported output format.")
		return 0
	}

	llm, err := models.CreateInferenceService(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error instanciating LLM")
		return 0
	}

	result, err := llm.ChatCompletion(ctx, instruction, sentence, outputFormat)
	if err != nil {
		logger.Err(ctx, err)
		return 0
	}
	if result.Error.Message != "" {
		err := fmt.Errorf(result.Error.Message)
		logger.Err(ctx, err)
		return 0
	}

	if models.OutputFormatJson == outputFormat {
		// safeguard: test is the output is a valid json
		// test every Choices.Message.Content
		for _, choice := range result.Choices {
			_, err := utils.JsonSerialize(choice.Message.Content)
			if err != nil {
				logger.Err(ctx, err).Msg("One of the generated message is not a valid JSON.")
				return 0
			}
		}
	}

	// return the first chat response
	if len(result.Choices) == 0 {
		logger.Error(ctx).Msg("Empty result returned from OpenAI.")
		return 0
	}
	content := result.Choices[0].Message.Content

	offset, err := assemblyscript.WriteString(ctx, mod, content)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}

	return offset
}

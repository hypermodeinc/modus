/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	hyp_aws "hmruntime/aws"
	"hmruntime/db"
	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/plugins"
	"hmruntime/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hypermodeAI/manifest"
	wasm "github.com/tetratelabs/wazero/api"
)

type modelInfo struct {
	Name     string
	FullName string
}

func (m *modelInfo) GetTypeInfo() plugins.TypeInfo {
	return plugins.TypeInfo{
		Name: "ModelInfo",
		Path: "~lib/@hypermode/models-as/index/ModelInfo",
	}
}

func hostLookupModel(ctx context.Context, mod wasm.Module, pModelName uint32) (pModelInfo uint32) {
	var modelName string
	err := readParam(ctx, mod, pModelName, &modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	info := modelInfo{
		Name:     model.Name,
		FullName: model.SourceModel,
	}

	offset, err := writeResult(ctx, mod, info)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}

func hostInvokeModel(ctx context.Context, mod wasm.Module, pModelName uint32, pInput uint32) (pOutput uint32) {

	var modelName, input string
	err := readParams2(ctx, mod, pModelName, pInput, &modelName, &input)
	if err != nil {
		logger.Err(ctx, err).Msg("Error reading input parameters.")
		return 0
	}

	model, err := models.GetModel(modelName)
	if err != nil {
		logger.Err(ctx, err).Msg("Error getting model.")
		return 0
	}

	// TODO: use the provider pattern instead of branching
	var output string
	if model.Host == "aws-bedrock" {
		output, err = invokeAwsBedrockModel(ctx, model, input)
	} else {
		output, err = models.PostToModelEndpoint[string](ctx, model, input)
	}

	if err != nil {
		logger.Err(ctx, err).Msg("Error posting to model endpoint.")
		return 0
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}

func invokeAwsBedrockModel(ctx context.Context, model manifest.ModelInfo, input string) (output string, err error) {

	// NOTE: Bedrock support is experimental, and not advertised to users.
	// It currently uses the same AWS credentials as the Runtime.
	// In the future, we will support user-provided credentials for Bedrock.

	cfg := hyp_aws.GetAwsConfig() // TODO, connect to AWS using the user-provided credentials
	client := bedrockruntime.NewFromConfig(cfg)

	modelId := fmt.Sprintf("%s.%s", model.Provider, model.SourceModel)

	startTime := utils.GetTime()
	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelId,
		ContentType: aws.String("application/json"),
		Body:        []byte(input),
	})
	endTime := utils.GetTime()

	output = string(result.Body)

	db.WriteInferenceHistory(ctx, model, input, output, startTime, endTime)

	return output, err
}

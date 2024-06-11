/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	hyp_aws "hmruntime/aws"
	"hmruntime/hosts"
	"hmruntime/logger"
	"hmruntime/models"
	"hmruntime/plugins"

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

	// TODO: use the provider pattern, not a switch

	var output string
	switch model.Host {
	case "hypermode":
		// not yet implemented
		logger.Error(ctx).Msg("Hypermode model host not yet implemented.")
		return 0
	case "aws-bedrock":
		output, err = invokeAwsBedrockModel(ctx, model, input)
		if err != nil {
			logger.Err(ctx, err).Msg("Error invoking AWS Bedrock model.")
			return 0
		}
	default:
		host, err := hosts.GetHost(model.Host)
		if err != nil {
			logger.Err(ctx, err).Msg("Error getting model host.")
			return 0
		}

		result, err := hosts.PostToHostEndpoint[string](ctx, host, input)
		if err != nil {
			logger.Err(ctx, err).Msg("Error posting to model endpoint.")
			return 0
		}

		output = result.Data
	}

	offset, err := writeResult(ctx, mod, output)
	if err != nil {
		logger.Err(ctx, err).Msg("Error writing result to wasm memory.")
		return 0
	}
	return offset
}

func invokeAwsBedrockModel(ctx context.Context, model manifest.ModelInfo, input string) (output string, err error) {

	cfg := hyp_aws.GetAwsConfig() // TODO, connect to AWS using the correct credentials
	client := bedrockruntime.NewFromConfig(cfg)

	modelId := fmt.Sprintf("%s.%s", model.Provider, model.SourceModel)

	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     &modelId,
		ContentType: aws.String("application/json"),
		Body:        []byte(input),
	})

	return string(result.Body), err
}

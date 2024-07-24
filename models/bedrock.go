/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import (
	"context"
	"fmt"

	hyp_aws "hmruntime/aws"
	"hmruntime/db"
	"hmruntime/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/hypermodeAI/manifest"
)

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

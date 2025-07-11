/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

// import (
// 	"context"
// 	"fmt"

// 	"github.com/hypermodeinc/modus/lib/manifest"
// 	hyp_aws "github.com/hypermodeinc/modus/runtime/aws"
// 	"github.com/hypermodeinc/modus/runtime/db"
// 	"github.com/hypermodeinc/modus/runtime/utils"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
// )

// func invokeAwsBedrockModel(ctx context.Context, model *manifest.ModelInfo, input string) (output string, err error) {

// 	span, ctx := sentryutils.NewSentrySpanForCurrentFunc(ctx)
// 	defer span.Finish()

// 	// NOTE: Bedrock support is experimental, and not advertised to users.
// 	// It currently uses the same AWS credentials as the Runtime.
// 	// In the future, we will support user-provided credentials for Bedrock.

// 	cfg := hyp_aws.GetAwsConfig() // TODO, connect to AWS using the user-provided credentials
// 	client := bedrockruntime.NewFromConfig(cfg)

// 	modelId := fmt.Sprintf("%s.%s", model.Provider, model.SourceModel)

// 	startTime := utils.GetTime()
// 	result, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
// 		ModelId:     &modelId,
// 		ContentType: aws.String("application/json"),
// 		Body:        []byte(input),
// 	})
// 	endTime := utils.GetTime()

// 	output = string(result.Body)

// 	db.WriteInferenceHistory(ctx, model, input, output, startTime, endTime)

// 	return output, err
// }

/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import (
	"context"
	"fmt"

	"hmruntime/utils"

	"github.com/tetratelabs/wazero"
)

const hostModuleName string = "hypermode"

func Instantiate(ctx context.Context, runtime *wazero.Runtime) error {
	span := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	b := (*runtime).NewHostModuleBuilder(hostModuleName)

	// Each host function should get a line here:
	b.NewFunctionBuilder().WithFunc(hostExecuteGQL).Export("executeGQL")
	b.NewFunctionBuilder().WithFunc(hostInvokeClassifier).Export("invokeClassifier")
	b.NewFunctionBuilder().WithFunc(hostComputeEmbedding).Export("computeEmbedding")
	b.NewFunctionBuilder().WithFunc(hostInvokeTextGenerator).Export("invokeTextGenerator")

	_, err := b.Instantiate(ctx)
	if err != nil {
		return fmt.Errorf("failed to instantiate the %s module: %w", hostModuleName, err)
	}

	return fmt.Errorf("this is a test error")

	return nil
}

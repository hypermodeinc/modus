/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"

	"hmruntime/logger"
	"hmruntime/plugins"
	"hmruntime/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
)

// Global runtime instance for the WASM modules
var RuntimeInstance wazero.Runtime

// Gets a module instance for the given plugin, used for a single invocation.
func GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers *utils.OutputBuffers) (wasm.Module, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Capture stdout/stderr both to logs, and to provided writers.
	wOut := io.MultiWriter(&buffers.StdOut, wInfoLog)
	wErr := io.MultiWriter(&buffers.StdErr, wErrorLog)

	// Configure the module instance.
	cfg := wazero.NewModuleConfig().
		WithName(plugin.Name() + "_" + uuid.NewString()).
		WithSysWalltime().WithSysNanotime().
		WithRandSource(rand.Reader).
		WithStdout(wOut).WithStderr(wErr)

	// Instantiate the plugin as a module.
	// NOTE: This will also invoke the plugin's `_start` function,
	// which will call any top-level code in the plugin.
	mod, err := RuntimeInstance.InstantiateModule(ctx, *plugin.Module, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, nil
}

/*
 * Copyright 2024 Hypermode, Inc.
 */

package wasmhost

import (
	"context"
	"crypto/rand"
	"fmt"
	"hmruntime/logger"
	"hmruntime/plugins"
	"io"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

// Global runtime instance for the WASM modules
var RuntimeInstance wazero.Runtime

// Channel used to signal that registration is needed
var RegistrationRequest chan bool = make(chan bool)

// Global, thread-safe registry of all plugins loaded by the host
var Plugins = plugins.NewPluginRegistry()

type OutputBuffers struct {
	Stdout *strings.Builder
	Stderr *strings.Builder
}

// Gets a module instance for the given plugin, used for a single invocation.
func GetModuleInstance(ctx context.Context, plugin *plugins.Plugin) (api.Module, OutputBuffers, error) {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Create string buffers to capture stdout and stderr.
	// Still write to the log, but also capture the output in the buffers.
	buf := OutputBuffers{&strings.Builder{}, &strings.Builder{}}
	wOut := io.MultiWriter(buf.Stdout, wInfoLog)
	wErr := io.MultiWriter(buf.Stderr, wErrorLog)

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
		return nil, buf, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, buf, nil
}

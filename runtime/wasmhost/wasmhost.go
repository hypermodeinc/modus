/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package wasmhost

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/hypermodeinc/modus/runtime/functions"
	"github.com/hypermodeinc/modus/runtime/logger"
	"github.com/hypermodeinc/modus/runtime/middleware"
	"github.com/hypermodeinc/modus/runtime/plugins"
	"github.com/hypermodeinc/modus/runtime/timezones"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/rs/zerolog"
	"github.com/tetratelabs/wazero"
	wasm "github.com/tetratelabs/wazero/api"
	wasi "github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type WasmHost interface {
	RegisterHostFunction(modName, funcName string, fn any, opts ...HostFunctionOption) error
	CallFunction(ctx context.Context, fnInfo functions.FunctionInfo, parameters map[string]any) (ExecutionInfo, error)
	CallFunctionByName(ctx context.Context, fnName string, paramValues ...any) (ExecutionInfo, error)
	CallFunctionInModule(ctx context.Context, mod wasm.Module, buffers utils.OutputBuffers, fnInfo functions.FunctionInfo, parameters map[string]any) (ExecutionInfo, error)
	Close(ctx context.Context)
	CompileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error)
	GetFunctionInfo(fnName string) (functions.FunctionInfo, error)
	GetFunctionRegistry() functions.FunctionRegistry
	GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers utils.OutputBuffers) (wasm.Module, error)
}

type wasmHost struct {
	runtime       wazero.Runtime
	fnRegistry    functions.FunctionRegistry
	hostFunctions []*hostFunction
}

func NewWasmHost(ctx context.Context, registrations ...func(WasmHost) error) WasmHost {
	cfg := wazero.NewRuntimeConfig().WithCloseOnContextDone(true)
	runtime := wazero.NewRuntimeWithConfig(ctx, cfg)
	wasi.MustInstantiate(ctx, runtime)

	if err := instantiateEnvHostFunctions(ctx, runtime); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to instantiate env host functions.")
		return nil
	}

	host := &wasmHost{
		runtime:    runtime,
		fnRegistry: functions.NewFunctionRegistry(),
	}

	for _, reg := range registrations {
		if err := reg(host); err != nil {
			logger.Fatal(ctx).Err(err).Msg("Failed to apply a registration to the WASM host.")
			return nil
		}
	}

	if err := host.instantiateHostFunctions(ctx); err != nil {
		logger.Fatal(ctx).Err(err).Msg("Failed to instantiate host functions.")
		return nil
	}

	return host
}

func GetWasmHost(ctx context.Context) WasmHost {
	host, ok := ctx.Value(utils.WasmHostContextKey).(WasmHost)
	if !ok {
		logger.Fatal(ctx).Msg("WASM Host not found in context.")
		return nil
	}
	return host
}

func (host *wasmHost) Close(ctx context.Context) {
	if err := host.runtime.Close(ctx); err != nil {
		logger.Err(ctx, err).Msg("Failed to cleanly close the WASM runtime.")
	}
}

func (host *wasmHost) GetFunctionInfo(fnName string) (functions.FunctionInfo, error) {
	return host.fnRegistry.GetFunctionInfo(fnName)
}

func (host *wasmHost) GetFunctionRegistry() functions.FunctionRegistry {
	return host.fnRegistry
}

// Gets a module instance for the given plugin, used for a single invocation.
func (host *wasmHost) GetModuleInstance(ctx context.Context, plugin *plugins.Plugin, buffers utils.OutputBuffers) (wasm.Module, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cfg := getModuleConfig(ctx, buffers, plugin)
	mod, err := host.runtime.InstantiateModule(ctx, plugin.Module, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate the plugin module: %w", err)
	}

	return mod, nil
}

func (host *wasmHost) CompileModule(ctx context.Context, bytes []byte) (wazero.CompiledModule, error) {
	span, ctx := utils.NewSentrySpanForCurrentFunc(ctx)
	defer span.Finish()

	cm, err := host.runtime.CompileModule(ctx, bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to compile the plugin: %w", err)
	}

	return cm, nil
}

func getModuleConfig(ctx context.Context, buffers utils.OutputBuffers, plugin *plugins.Plugin) wazero.ModuleConfig {

	// Get the logger and writers for the plugin's stdout and stderr.
	log := logger.Get(ctx).With().Bool("user_visible", true).Logger()
	wInfoLog := logger.NewLogWriter(&log, zerolog.InfoLevel)
	wErrorLog := logger.NewLogWriter(&log, zerolog.ErrorLevel)

	// Capture stdout/stderr both to logs, and to provided writers.
	wOut := io.MultiWriter(buffers.StdOut(), wInfoLog)
	wErr := io.MultiWriter(buffers.StdErr(), wErrorLog)

	// Get the time zone to pass to the module instance.
	var timeZone string
	if tz, ok := ctx.Value(utils.TimeZoneContextKey).(string); ok {
		timeZone = tz
	} else {
		timeZone = timezones.GetLocalTimeZone()
	}

	// Get the JWT claims to pass to the module instance.
	jwtClaims := middleware.GetJWTClaims(ctx)

	// Note, we use an anonymous module name (empty string) here,
	// for concurrency and performance reasons.
	// See https://github.com/tetratelabs/wazero/pull/2275
	// And https://gophers.slack.com/archives/C040AKTNTE0/p1719587772724619?thread_ts=1719522663.531579&cid=C040AKTNTE0
	cfg := wazero.NewModuleConfig().
		WithName("").
		WithStartFunctions(plugin.StartFunction).
		WithSysWalltime().WithSysNanotime().
		WithRandSource(rand.Reader).
		WithStdout(wOut).WithStderr(wErr).
		WithEnv("TZ", timeZone).
		WithEnv("CLAIMS", jwtClaims)

	for _, env := range os.Environ() {
		split := strings.SplitN(env, "=", 2)
		key, val := split[0], split[1]
		if strings.HasPrefix(key, "MODUS_") {
			// Remove the MODUS_ prefix
			cfg = cfg.WithEnv(key[6:], val)
		}
	}

	return cfg
}

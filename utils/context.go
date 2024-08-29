/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

type contextKey string

const ExecutionIdContextKey contextKey = "execution_id"
const PluginContextKey contextKey = "plugin"
const MetadataContextKey contextKey = "metadata"
const WasmAdapterContextKey contextKey = "wasm_adapter"
const FunctionNameContextKey contextKey = "function_name"
const FunctionOutputContextKey contextKey = "function_output"
const FunctionMessagesContextKey contextKey = "function_messages"
const CustomTypesContextKey contextKey = "custom_types"

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

type contextKey string

const WasmHostContextKey contextKey = "wasm_host"
const ExecutionIdContextKey contextKey = "execution_id"
const PluginContextKey contextKey = "plugin"
const MetadataContextKey contextKey = "metadata"
const WasmAdapterContextKey contextKey = "wasm_adapter"
const FunctionNameContextKey contextKey = "function_name"
const FunctionOutputContextKey contextKey = "function_output"
const FunctionMessagesContextKey contextKey = "function_messages"
const CustomTypesContextKey contextKey = "custom_types"

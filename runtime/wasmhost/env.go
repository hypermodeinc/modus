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
	"time"

	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
)

func instantiateEnvHostFunctions(ctx context.Context, r wazero.Runtime) error {

	// AssemblyScript date import uses:
	//
	//   @external("env", "Date.now")
	//   export function now(): f64;
	//
	// We use this instead of relying on the date override in the AssemblyScript wasi shim,
	// because the shim's approach leads to problems such as:
	//   ERROR TS2322: Type '~lib/date/Date' is not assignable to type '~lib/wasi_date/wasi_Date'.
	//
	dateNow := func(ctx context.Context, stack []uint64) {
		now := time.Now().UnixMilli()
		stack[0] = api.EncodeF64(float64(now))
	}

	_, err := r.NewHostModuleBuilder("env").
		NewFunctionBuilder().WithGoFunction(
		api.GoFunc(dateNow), nil, []api.ValueType{api.ValueTypeF64}).Export("Date.now").
		Instantiate(ctx)

	return err
}

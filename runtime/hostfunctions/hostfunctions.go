/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import "hypruntime/wasmhost"

var registrations []func(wasmhost.WasmHost) error

func GetRegistrations() []func(wasmhost.WasmHost) error {
	return registrations
}

func registerHostFunction(modName, funcName string, fn any, opts ...wasmhost.HostFunctionOption) {
	registrations = append(registrations, func(host wasmhost.WasmHost) error {
		return host.RegisterHostFunction(modName, funcName, fn, opts...)
	})
}

func withStartingMessage(text string) wasmhost.HostFunctionOption {
	return wasmhost.WithStartingMessage(text)
}

func withCompletedMessage(text string) wasmhost.HostFunctionOption {
	return wasmhost.WithCompletedMessage(text)
}

func withCancelledMessage(text string) wasmhost.HostFunctionOption {
	return wasmhost.WithCancelledMessage(text)
}

func withErrorMessage(text string) wasmhost.HostFunctionOption {
	return wasmhost.WithErrorMessage(text)
}

func withMessageDetail(fn any) wasmhost.HostFunctionOption {
	return wasmhost.WithMessageDetail(fn)
}

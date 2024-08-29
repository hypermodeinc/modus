/*
 * Copyright 2024 Hypermode, Inc.
 */

package hostfunctions

import "hypruntime/wasmhost"

var registrations []func(*wasmhost.WasmHost) error

func GetRegistrations() []func(*wasmhost.WasmHost) error {
	return registrations
}

func registerHostFunction(modName, funcName string, fn any, opts ...func(*wasmhost.HostFunction)) {
	registrations = append(registrations, func(host *wasmhost.WasmHost) error {
		return host.RegisterHostFunction(modName, funcName, fn, opts...)
	})
}

func withStartingMessage(text string) func(*wasmhost.HostFunction) {
	return wasmhost.WithStartingMessage(text)
}

func withCompletedMessage(text string) func(*wasmhost.HostFunction) {
	return wasmhost.WithCompletedMessage(text)
}

func withCancelledMessage(text string) func(*wasmhost.HostFunction) {
	return wasmhost.WithCancelledMessage(text)
}

func withErrorMessage(text string) func(*wasmhost.HostFunction) {
	return wasmhost.WithErrorMessage(text)
}

func withMessageDetail(fn any) func(*wasmhost.HostFunction) {
	return wasmhost.WithMessageDetail(fn)
}

//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import "unsafe"

//go:noescape
//go:wasmimport modus_secrets getSecretValue
func _hostGetSecretValue(name unsafe.Pointer) unsafe.Pointer

//modus:import modus_secrets getSecretValue
func hostGetSecretValue(name *string) *string {
	return (*string)(_hostGetSecretValue(unsafe.Pointer(name)))
}

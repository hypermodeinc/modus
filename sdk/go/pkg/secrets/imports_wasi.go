//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
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

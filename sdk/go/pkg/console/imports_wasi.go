//go:build wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console

//go:noescape
//go:wasmimport modus_system logMessage
func _hostLogMessage(level, message *string)

func hostLogMessage(level, message string) {
	_hostLogMessage(&level, &message)
}

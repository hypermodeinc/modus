//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package console

//go:noescape
//go:wasmimport modus_system logMessage
func _hostLogMessage(level, message *string)

func hostLogMessage(level, message string) {
	_hostLogMessage(&level, &message)
}

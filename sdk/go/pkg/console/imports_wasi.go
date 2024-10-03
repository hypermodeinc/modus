//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package console

//go:noescape
//go:wasmimport hypermode log
func _log(level, message *string)

func log(level, message string) {
	_log(&level, &message)
}

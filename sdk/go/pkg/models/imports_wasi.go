//go:build wasip1

/*
 * Copyright 2024 Hypermode, Inc.
 */

package models

import "unsafe"

//go:noescape
//go:wasmimport hypermode lookupModel
func _lookupModel(modelName *string) unsafe.Pointer

//hypermode:import hypermode lookupModel
func lookupModel(modelName *string) *ModelInfo {
	info := _lookupModel(modelName)
	if info == nil {
		return nil
	}
	return (*ModelInfo)(info)
}

//go:noescape
//go:wasmimport hypermode invokeModel
func invokeModel(modelName *string, input *string) *string

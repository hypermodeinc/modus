//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import "unsafe"

//go:noescape
//go:wasmimport modus_models getModelInfo
func _hostGetModelInfo(modelName *string) unsafe.Pointer

//modus:import modus_models getModelInfo
func hostGetModelInfo(modelName *string) *ModelInfo {
	info := _hostGetModelInfo(modelName)
	if info == nil {
		return nil
	}
	return (*ModelInfo)(info)
}

//go:noescape
//go:wasmimport modus_models invokeModel
func hostInvokeModel(modelName *string, input *string) *string

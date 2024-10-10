//go:build wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package models

import "unsafe"

//go:noescape
//go:wasmimport modus_models lookupModel
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
//go:wasmimport modus_models invokeModel
func invokeModel(modelName *string, input *string) *string

/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package schemagen

import "github.com/hypermodeinc/modus/runtime/manifestdata"

func getFnFilter() func(*FunctionSignature) bool {
	embedders := make(map[string]bool)
	for _, collection := range manifestdata.GetManifest().Collections {
		for _, searchMethod := range collection.SearchMethods {
			embedders[searchMethod.Embedder] = true
		}
	}

	return func(f *FunctionSignature) bool {
		return !embedders[f.Name]
	}
}

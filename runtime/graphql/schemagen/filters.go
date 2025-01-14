/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package schemagen

import "github.com/hypermodeinc/modus/runtime/manifestdata"

func getFieldFilter() func(*FieldDefinition) bool {
	embedders := getEmbedderFields()
	return func(f *FieldDefinition) bool {
		return !embedders[f.Name]
	}
}

func getEmbedderFields() map[string]bool {
	embedders := make(map[string]bool)
	for _, collection := range manifestdata.GetManifest().Collections {
		for _, searchMethod := range collection.SearchMethods {
			embedders[getFieldName(searchMethod.Embedder)] = true
		}
	}
	return embedders
}

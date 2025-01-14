/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

type CollectionInfo struct {
	Name          string                      `json:"-"`
	SearchMethods map[string]SearchMethodInfo `json:"searchMethods"`
}

type SearchMethodInfo struct {
	Embedder string    `json:"embedder"`
	Index    IndexInfo `json:"index"`
}

type IndexInfo struct {
	Type    string      `json:"type"`
	Options OptionsInfo `json:"options"`
}

type OptionsInfo struct {
	EfConstruction int `json:"efConstruction"`
	MaxLevels      int `json:"maxLevels"`
}

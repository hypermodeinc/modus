/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package manifest

type ModelInfo struct {
	Name        string `json:"-"`
	SourceModel string `json:"sourceModel"`
	Provider    string `json:"provider"`
	Connection  string `json:"connection"`
	Path        string `json:"path"`
}

func (m ModelInfo) Hash() string {
	return computeHash(m.Name, m.SourceModel, m.Provider, m.Connection, m.Path)
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package app

import (
	"strings"
	"testing"
)

func TestVersionNumber(t *testing.T) {
	original := version
	t.Cleanup(func() {
		version = original
	})

	tests := []struct {
		name    string
		version string
	}{
		{"Empty version", ""},
		{"Non-empty version", "1.2.3"},
		{"Non-empty version with 'v'", "v1.2.3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			version = tt.version
			adjustVersion()

			v := VersionNumber()
			if v == "" {
				t.Errorf("Expected version number to be non-empty")
			}

			if v[0] != 'v' {
				t.Errorf("Expected version number to start with 'v'")
			}

			if len(v) < 6 {
				t.Errorf("Expected version number to be at least 6 characters long")
			}

			if strings.Count(v, ".") < 2 {
				t.Errorf("Expected version number to have at least two dots")
			}
		})
	}
}

func TestProductVersion(t *testing.T) {
	v := VersionNumber()
	pv := ProductVersion()

	if !strings.HasSuffix(pv, " "+v) {
		t.Errorf("Expected product version to end with version number")
	}

	if len(pv) < len(v)+2 {
		t.Errorf("Expected product version to be at least %d characters long", len(v)+2)
	}
}

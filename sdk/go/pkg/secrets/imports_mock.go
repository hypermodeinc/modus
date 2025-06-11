//go:build !wasip1

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

// hostGetSecretValue is a mock implementation for tests.
func hostGetSecretValue(name *string) *string {
	if name == nil {
		return nil
	}
	value := "value"
	return &value
}

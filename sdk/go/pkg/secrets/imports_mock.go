//go:build !wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

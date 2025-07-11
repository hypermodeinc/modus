/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package secrets

import "errors"

// GetSecretValue retrieves a secret value from the host environment.
func GetSecretValue(name string) (string, error) {
	resultPtr := hostGetSecretValue(&name)
	if resultPtr == nil {
		return "", errors.New("secret not found")
	}
	return *resultPtr, nil
}

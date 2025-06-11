/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
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

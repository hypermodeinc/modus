/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package utils

import (
	"errors"
	"fmt"
)

var ErrUserError = errors.New("")

func NewUserError(err error) error {
	return fmt.Errorf("%w%w", ErrUserError, err)
}

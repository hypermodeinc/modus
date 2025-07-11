/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

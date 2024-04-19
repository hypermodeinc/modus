/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"errors"
)

var errFunctionException = errors.New("function raised an exception")

func ShouldReturnErrorToResponse(err error) bool {
	return err != errFunctionException
}

func transformError(err error) error {
	if err.Error() == "module closed with exit_code(255)" {
		return errFunctionException
	}

	return err
}

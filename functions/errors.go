/*
 * Copyright 2024 Hypermode, Inc.
 */

package functions

import (
	"errors"

	"github.com/tetratelabs/wazero/sys"
)

var errModuleClosed = sys.NewExitError(255)
var errFunctionException = errors.New("function raised an exception")

func ShouldReturnErrorToResponse(err error) bool {
	return err != errFunctionException
}

func transformError(err error) error {
	if errors.Is(err, errModuleClosed) {
		return errFunctionException
	}

	return err
}

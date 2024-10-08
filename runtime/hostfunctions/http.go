/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/httpclient"
)

func init() {
	registerHostFunction("hypermode", "httpFetch", httpclient.HttpFetch,
		withStartingMessage("Starting HTTP request."),
		withCompletedMessage("Completed HTTP request."),
		withCancelledMessage("Cancelled HTTP request."),
		withErrorMessage("Error making HTTP request."),
		withMessageDetail(func(request *httpclient.HttpRequest) string {
			return fmt.Sprintf("%s %s", request.Method, request.Url)
		}))
}

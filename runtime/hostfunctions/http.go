/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/httpclient"
)

func init() {
	const module_name = "modus_http_client"

	registerHostFunction(module_name, "fetch", httpclient.Fetch,
		withStartingMessage("Starting HTTP request."),
		withCompletedMessage("Completed HTTP request."),
		withCancelledMessage("Cancelled HTTP request."),
		withErrorMessage("Error making HTTP request."),
		withMessageDetail(func(request *httpclient.HttpRequest) string {
			return fmt.Sprintf("%s %s", request.Method, request.Url)
		}))
}

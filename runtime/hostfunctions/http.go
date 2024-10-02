/*
 * Copyright 2024 Hypermode, Inc.
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

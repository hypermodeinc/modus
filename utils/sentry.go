/*
 * Copyright 2024 Hypermode, Inc.
 */

package utils

import (
	"context"
	"runtime"

	"github.com/getsentry/sentry-go"
)

func NewSentrySpanForCurrentFunc(ctx context.Context) *sentry.Span {
	span := sentry.StartSpan(ctx, "function")
	span.Description = GetFuncName(2)
	return span
}

func GetCurrentFuncName() string {
	return GetFuncName(1)
}

func GetFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	return fn.Name()
}

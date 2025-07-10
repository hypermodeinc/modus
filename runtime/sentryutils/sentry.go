/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package sentryutils

import (
	"context"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/hypermodeinc/modus/runtime/app"
	"github.com/hypermodeinc/modus/runtime/utils"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog/log"
)

const max_breadcrumbs = 100

var rootSourcePath = func() string {
	pc, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	callerName := runtime.FuncForPC(pc).Name()
	depth := strings.Count(callerName, "/") + 1
	s := filename
	for range depth {
		s = path.Dir(s)
	}
	return s + "/"
}()

var thisPackagePath = func() string {
	pc, _, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}

	callerName := runtime.FuncForPC(pc).Name()
	i := max(strings.LastIndexByte(callerName, '/'), 0)
	j := strings.IndexByte(callerName[i:], '.')
	return callerName[0 : i+j]
}()

func InitializeSentry() {

	// Don't initialize Sentry when running in debug mode.
	if utils.DebugModeEnabled() {
		return
	}

	// ONLY report errors to Sentry when the SENTRY_DSN environment variable is set.
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		return
	}

	// Allow the Sentry environment to be overridden by the SENTRY_ENVIRONMENT environment variable,
	// but default to the environment name from MODUS_ENV.
	environment := os.Getenv("SENTRY_ENVIRONMENT")
	if environment == "" {
		environment = app.Config().Environment()
	}

	// Allow the Sentry release to be overridden by the SENTRY_RELEASE environment variable,
	// but default to the Modus version number.
	release := os.Getenv("SENTRY_RELEASE")
	if release == "" {
		release = app.VersionNumber()
	}

	err := sentry.Init(sentry.ClientOptions{
		Dsn:                   dsn,
		Environment:           environment,
		Release:               release,
		BeforeSend:            sentryBeforeSend,
		BeforeSendTransaction: sentryBeforeSendTransaction,

		// Note - We use Prometheus for _metrics_ (see hypruntime/metrics package).
		// However, we can still use Sentry for _tracing_ to allow us to improve performance of the Runtime.
		// We should only trace code that we expect to run in a roughly consistent amount of time.
		// For example, we should not trace a user's function execution, or the outer GraphQL request handling,
		// because these can take an arbitrary amount of time depending on what the user's code does.
		// Instead, we should trace the runtime's startup, storage, secrets retrieval, schema generation, etc.
		// That way we can trace performance issues in the runtime itself, and let Sentry correlate them with
		// any errors that may have occurred.
		EnableTracing:    true,
		TracesSampleRate: utils.GetFloatFromEnv("SENTRY_TRACES_SAMPLE_RATE", 0.2),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize Sentry.")
	}
}

func CloseSentry() {
	if err := recover(); err != nil {
		hub := sentry.CurrentHub()
		hub.Recover(err)
	}
	sentry.Flush(2 * time.Second)
}

func NewSpanForCallingFunc(ctx context.Context) (*sentry.Span, context.Context) {
	funcName := getFuncName(3)
	return NewSpan(ctx, funcName)
}

func NewSpanForCurrentFunc(ctx context.Context) (*sentry.Span, context.Context) {
	funcName := getFuncName(2)
	return NewSpan(ctx, funcName)
}

func NewSpan(ctx context.Context, funcName string) (*sentry.Span, context.Context) {
	if tx := sentry.TransactionFromContext(ctx); tx == nil {
		tx = sentry.StartTransaction(ctx, funcName, sentry.WithOpName("function"))
		return tx, tx.Context()
	}

	span := sentry.StartSpan(ctx, "function")
	span.Description = funcName
	return span, span.Context()
}

func getFuncName(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		return "?"
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "?"
	}

	name := fn.Name()
	return utils.TrimStringBefore(name, "/")
}

func sentryBeforeSend(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	for _, e := range event.Exception {
		if e.Stacktrace != nil {
			frames := make([]sentry.Frame, 0, len(e.Stacktrace.Frames))

			// Keep frames except those from the current package.
			for f := range e.Stacktrace.Frames {
				if e.Stacktrace.Frames[f].Module != thisPackagePath {
					frames = append(frames, e.Stacktrace.Frames[f])
				}
			}

			// Adjust frames to include relative paths to source files.
			for i, f := range frames {
				if f.Filename == "" && f.AbsPath != "" && strings.HasPrefix(f.AbsPath, rootSourcePath) {
					f.Filename = f.AbsPath[len(rootSourcePath):]
					frames[i] = f
				}
			}

			e.Stacktrace.Frames = frames
		}
	}

	addExtraData(event)
	return event
}

var ignoredTransactions = map[string]bool{
	"GET /health":  true,
	"GET /metrics": true,
}

func sentryBeforeSendTransaction(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
	if ignoredTransactions[event.Transaction] {
		return nil
	}

	addExtraData(event)
	return event
}

// Include any extra information that may be useful for debugging.
func addExtraData(event *sentry.Event) {
	if ns, ok := app.KubernetesNamespace(); ok {
		addExtraDataToEvent(event, "namespace", ns)
	}
}

func addExtraDataToEvent(event *sentry.Event, key string, value any) {
	if event.Extra == nil {
		event.Extra = make(map[string]any)
	}
	event.Extra[key] = value
}

func Recover(ctx context.Context, r any) {
	ActiveHub(ctx).RecoverWithContext(ctx, r)
}

func CaptureError(ctx context.Context, err error, msg string, opts ...func(event *sentry.Event)) {
	hub := ActiveHub(ctx)
	client := hub.Client()
	if client == nil {
		return
	}

	var event *sentry.Event
	if err == nil {
		event = client.EventFromMessage(msg, sentry.LevelError)
	} else {
		event = client.EventFromException(err, sentry.LevelError)
		event.Message = msg
	}
	for _, opt := range opts {
		opt(event)
	}
	hub.CaptureEvent(event)
}

func CaptureWarning(ctx context.Context, err error, msg string, opts ...func(event *sentry.Event)) {
	hub := ActiveHub(ctx)
	client := hub.Client()
	if client == nil {
		return
	}

	var event *sentry.Event
	if err == nil {
		event = client.EventFromMessage(msg, sentry.LevelWarning)
	} else {
		event = client.EventFromException(err, sentry.LevelWarning)
		event.Message = msg
	}
	for _, opt := range opts {
		opt(event)
	}
	hub.CaptureEvent(event)
}

func WithData(key string, value any) func(event *sentry.Event) {
	return func(event *sentry.Event) {
		addExtraDataToEvent(event, key, value)
	}
}

func ActiveHub(ctx context.Context) *sentry.Hub {
	if hub := sentry.GetHubFromContext(ctx); hub != nil {
		return hub
	}

	return sentry.CurrentHub()
}

func NewScope(ctx context.Context) (scope *sentry.Scope, done func()) {
	hub := ActiveHub(ctx)
	scope = hub.PushScope()
	return scope, func() {
		hub.PopScope()
	}
}

func AddBreadcrumb(ctx context.Context, breadcrumb *sentry.Breadcrumb) {
	ActiveHub(ctx).ConfigureScope(func(scope *sentry.Scope) {
		scope.AddBreadcrumb(breadcrumb, max_breadcrumbs)
	})
}

func AddTextBreadcrumb(ctx context.Context, message string) {
	AddBreadcrumb(ctx, &sentry.Breadcrumb{
		Message:   message,
		Timestamp: time.Now().UTC(),
	})
}

func AddBreadcrumbToScope(scope *sentry.Scope, breadcrumb *sentry.Breadcrumb) {
	scope.AddBreadcrumb(breadcrumb, max_breadcrumbs)
}

func AddTextBreadcrumbToScope(scope *sentry.Scope, message string) {
	AddBreadcrumbToScope(scope, &sentry.Breadcrumb{
		Message:   message,
		Timestamp: time.Now().UTC(),
	})
}

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package actors

import (
	"context"
	"fmt"
	"time"

	"github.com/hypermodeinc/modus/runtime/sentryutils"
	goakt "github.com/tochemey/goakt/v3/actor"

	"google.golang.org/protobuf/proto"
)

// Sends a message to an actor identified by its name.
// Uses either Tell or RemoteTell based on whether the actor is local or remote.
func tell(ctx context.Context, actorName string, message proto.Message) error {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	addr, pid, err := _actorSystem.ActorOf(ctx, actorName)
	if err != nil {
		return err
	} else if pid != nil {
		return goakt.Tell(ctx, pid, message)
	} else if addr != nil {
		return goakt.NoSender.RemoteTell(ctx, addr, message)
	}
	return fmt.Errorf("failed to get address or PID for actor %s", actorName)
}

// Sends a message to an actor identified by its name, then waits for a response within the timeout duration.
// Uses either Ask or RemoteAsk based on whether the actor is local or remote.
func ask(ctx context.Context, actorName string, message proto.Message, timeout time.Duration) (response proto.Message, err error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	addr, pid, err := _actorSystem.ActorOf(ctx, actorName)
	if err != nil {
		return nil, err
	} else if pid != nil {
		return goakt.Ask(ctx, pid, message, timeout)
	} else if addr != nil {
		response, err := goakt.NoSender.RemoteAsk(ctx, addr, message, timeout)
		if err != nil {
			return nil, err
		}
		return response.UnmarshalNew()
	}
	return nil, fmt.Errorf("failed to get address or PID for actor %s", actorName)
}

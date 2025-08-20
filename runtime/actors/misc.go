/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

	return goakt.NoSender.SendAsync(ctx, actorName, message)
}

// Sends a message to an actor identified by its name, then waits for a response within the timeout duration.
// Uses either Ask or RemoteAsk based on whether the actor is local or remote.
func ask(ctx context.Context, actorName string, message proto.Message, timeout time.Duration) (response proto.Message, err error) {
	span, ctx := sentryutils.NewSpanForCurrentFunc(ctx)
	defer span.Finish()

	
	return  goakt.NoSender.SendSync(ctx, actorName, message, timeout)
}

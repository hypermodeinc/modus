//go:build !wasip1

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package agents

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/testutils"
)

var SpawnAgentActorCallStack = testutils.NewCallStack()
var SendMessageCallStack = testutils.NewCallStack()
var TerminateAgentCallStack = testutils.NewCallStack()

func hostSpawnAgentActor(agentName *string) *AgentInfo {
	SpawnAgentActorCallStack.Push(agentName)

	return &AgentInfo{
		Id:     "abc123",
		Name:   *agentName,
		Status: AgentStatusStarting,
	}
}

func hostSendMessage(agentId, msgName, data *string, timeout int64) *MessageResponse {
	SendMessageCallStack.Push(agentId, msgName, data, timeout)

	if *agentId == "abc123" {
		return &MessageResponse{
			data: data,
		}
	}

	return nil
}

func hostTerminateAgent(agentId *string) bool {
	TerminateAgentCallStack.Push(agentId)

	return agentId != nil && *agentId == "abc123"
}

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

var StartAgentCallStack = testutils.NewCallStack()
var SendMessageCallStack = testutils.NewCallStack()
var StopAgentCallStack = testutils.NewCallStack()

func hostStartAgent(agentName *string) *AgentInfo {
	StartAgentCallStack.Push(agentName)

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

func hostStopAgent(agentId *string) bool {
	StopAgentCallStack.Push(agentId)

	return agentId != nil && *agentId == "abc123"
}

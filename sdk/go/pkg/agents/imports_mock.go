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
var GetAgentInfoCallStack = testutils.NewCallStack()
var ListAgentsCallStack = testutils.NewCallStack()
var PublishEventCallStack = testutils.NewCallStack()

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

func hostStopAgent(agentId *string) *AgentInfo {
	StopAgentCallStack.Push(agentId)

	if *agentId == "abc123" {
		return &AgentInfo{
			Id:     "abc123",
			Name:   "Counter",
			Status: AgentStatusStopping,
		}
	}
	return nil
}

func hostGetAgentInfo(agentId *string) *AgentInfo {
	GetAgentInfoCallStack.Push(agentId)

	if *agentId == "abc123" {
		return &AgentInfo{
			Id:     "abc123",
			Name:   "Counter",
			Status: AgentStatusRunning,
		}
	}

	return nil
}

func hostListAgents() *[]AgentInfo {
	ListAgentsCallStack.Push()

	return &[]AgentInfo{
		{Id: "abc123", Name: "Counter", Status: AgentStatusRunning},
		{Id: "def456", Name: "Logger", Status: AgentStatusRunning},
	}
}

func hostPublishEvent(agentId, eventName, eventData *string) {
	PublishEventCallStack.Push(agentId, eventName, eventData)
}

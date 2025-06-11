/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package hostfunctions

import (
	"fmt"

	"github.com/hypermodeinc/modus/runtime/actors"
)

func init() {
	const module_name = "modus_agents"

	registerHostFunction(module_name, "startAgent", actors.StartAgent,
		withErrorMessage("Error starting agent."),
		withMessageDetail(func(agentName string) string {
			return fmt.Sprintf("Name: %s", agentName)
		}))

	registerHostFunction(module_name, "stopAgent", actors.StopAgent,
		withErrorMessage("Error stopping agent."),
		withMessageDetail(func(agentId string) string {
			return fmt.Sprintf("AgentId: %s", agentId)
		}))

	registerHostFunction(module_name, "getAgentInfo", actors.GetAgentInfo,
		withErrorMessage("Error getting agent info."),
		withMessageDetail(func(agentId string) string {
			return fmt.Sprintf("AgentId: %s", agentId)
		}))

	registerHostFunction(module_name, "listAgents", actors.ListActiveAgents,
		withErrorMessage("Error listing agents."))

	registerHostFunction(module_name, "sendMessage", actors.SendAgentMessage,
		withErrorMessage("Error sending message to agent."),
		withMessageDetail(func(agentId string, msgName string, data *string, timeout int64) string {
			return fmt.Sprintf("AgentId: %s, MsgName: %s", agentId, msgName)
		}))

	registerHostFunction(module_name, "publishEvent", actors.PublishAgentEvent,
		withErrorMessage("Error publishing agent event."),
		withMessageDetail(func(agentId, eventName string, eventData *string, createdAt *string) string {
			return fmt.Sprintf("AgentId: %s, EventName: %s", agentId, eventName)
		}))
}

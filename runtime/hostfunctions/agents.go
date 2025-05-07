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

	registerHostFunction(module_name, "registerAgent", actors.RegisterAgent,
		withErrorMessage("Error registering agent."),
		withMessageDetail(func(agentId int32, name string) string {
			return fmt.Sprintf("Id: %d, Name: %s", agentId, name)
		}))

	registerHostFunction(module_name, "sendMessage", actors.SendAgentMessage,
		withErrorMessage("Error sending message to agent."),
		withMessageDetail(func(agentId int32, msgName string, data *string, timeout int64) string {
			return fmt.Sprintf("AgentId: %d, MsgName: %s", agentId, msgName)
		}))
}

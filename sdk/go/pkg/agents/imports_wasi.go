//go:build wasip1

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
	"unsafe"
)

//go:noescape
//go:wasmimport modus_agents spawnAgentActor
func _hostSpawnAgentActor(agentName *string) unsafe.Pointer

//modus:import modus_agents spawnAgentActor
func hostSpawnAgentActor(agentName *string) *AgentInfo {
	info := _hostSpawnAgentActor(agentName)
	if info == nil {
		return nil
	}
	return (*AgentInfo)(info)
}

//go:noescape
//go:wasmimport modus_agents sendMessage
func _hostSendMessage(agentId, msgName, data *string, timeout int64) unsafe.Pointer

//modus:import modus_agents sendMessage
func hostSendMessage(agentId, msgName, data *string, timeout int64) *MessageResponse {
	response := _hostSendMessage(agentId, msgName, data, timeout)
	if response == nil {
		return nil
	}
	return (*MessageResponse)(response)
}

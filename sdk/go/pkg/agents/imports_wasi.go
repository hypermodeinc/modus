//go:build wasip1

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

package agents

import (
	"unsafe"
)

//go:noescape
//go:wasmimport modus_agents startAgent
func _hostStartAgent(agentName *string) unsafe.Pointer

//modus:import modus_agents startAgent
func hostStartAgent(agentName *string) *AgentInfo {
	info := _hostStartAgent(agentName)
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

//go:noescape
//go:wasmimport modus_agents stopAgent
func _hostStopAgent(agentId *string) unsafe.Pointer

//modus:import modus_agents stopAgent
func hostStopAgent(agentId *string) *AgentInfo {
	info := _hostStopAgent(agentId)
	if info == nil {
		return nil
	}
	return (*AgentInfo)(info)
}

//go:noescape
//go:wasmimport modus_agents getAgentInfo
func _hostGetAgentInfo(agentId *string) unsafe.Pointer

//modus:import modus_agents getAgentInfo
func hostGetAgentInfo(agentId *string) *AgentInfo {
	info := _hostGetAgentInfo(agentId)
	if info == nil {
		return nil
	}
	return (*AgentInfo)(info)
}

//go:noescape
//go:wasmimport modus_agents listAgents
func _hostListAgents() unsafe.Pointer

//modus:import modus_agents listAgents
func hostListAgents() *[]AgentInfo {
	ptr := _hostListAgents()
	if ptr == nil {
		return nil
	}

	return (*[]AgentInfo)(ptr)
}

//go:noescape
//go:wasmimport modus_agents publishEvent
func hostPublishEvent(agentId, eventName, eventData *string)

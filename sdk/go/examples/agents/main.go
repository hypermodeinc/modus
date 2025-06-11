/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

package main

import (
	"strconv"

	"github.com/hypermodeinc/modus/sdk/go/pkg/agents"
)

// All agents must be registered before they can be used.
// This is done in an init function.
func init() {
	agents.Register(&CounterAgent{})
}

// The following are regular Modus functions.
// They are not part of the agent, but are used to start the agent and interact with it.
// Note that they cannot use an instance of the CounterAgent struct directly,
// but rather they will start an instance by name, and then send messages to it by ID.
// This is because the agent instance will actually be running in a different WASM instance,
// perhaps on a different process or even on a different machine.

// Starts a counter agent and returns info including its ID and status.
func StartCounterAgent() (agents.AgentInfo, error) {
	return agents.Start("Counter")
}

// Starts a counter agent with an initial count value,
// and returns info including its ID and status.
func StartCounterAgentWithData(initialCount int) (agents.AgentInfo, error) {
	info, err := agents.Start("Counter")
	if err != nil {
		return agents.AgentInfo{}, err
	}

	err = UpdateCountAsync(info.Id, initialCount)
	return info, err
}

// Stops the specified agent by ID, returning its status info.
// This will terminate the agent, and it cannot be resumed or restarted.
// However, a new agent with the same name can be started at any time.
func StopAgent(agentId string) (agents.AgentInfo, error) {
	return agents.Stop(agentId)
}

// Gets information about the specified agent.
func GetAgentInfo(agentId string) (agents.AgentInfo, error) {
	return agents.GetInfo(agentId)
}

// List all agents, except those that have been fully terminated.
func ListAgents() ([]agents.AgentInfo, error) {
	return agents.ListAll()
}

// Returns the current count of the specified agent.
func GetCount(agentId string) (int, error) {
	count, err := agents.SendMessage(agentId, "count")
	if err != nil {
		return 0, err
	}
	if count == nil {
		return 0, nil
	}
	return strconv.Atoi(*count)
}

// Increments the count of the specified agent by 1 and returns the new count.
func UpdateCount(agentId string) (int, error) {
	count, err := agents.SendMessage(agentId, "increment")
	if err != nil {
		return 0, err
	}
	if count == nil {
		return 0, nil
	}
	return strconv.Atoi(*count)
}

// Increments the count of the specified agent by the specified quantity.
// This is an asynchronous operation and does not return a value.
func UpdateCountAsync(agentId string, qty int) error {
	return agents.SendMessageAsync(agentId, "increment", agents.WithData(strconv.Itoa(qty)))
}

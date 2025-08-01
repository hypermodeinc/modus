/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { agents, AgentInfo } from "@hypermode/modus-sdk-as";
import { CounterAgent } from "./counterAgent";

// All agents must be registered before they can be used.
// Note that this is done outside of any function.
agents.register<CounterAgent>();

// The following are regular Modus functions.
// They are not part of the agent, but are used to start the agent and interact with it.
// Note that they cannot use an instance of the CounterAgent class directly,
// but rather they will start an instance by name, and then send messages to it by ID.
// This is because the agent instance will actually be running in a different WASM instance,
// perhaps on a different process or even on a different machine.

/**
 * Starts a counter agent and returns info including its ID and status.
 */
export function startCounterAgent(): AgentInfo {
  return agents.start("Counter");
}

/**
 * Starts a counter agent with an initial count value,
 * and returns info including its ID and status.
 */
export function startCounterAgentWithData(initialCount: i32): AgentInfo {
  const info = agents.start("Counter");
  updateCountAsync(info.id, initialCount);
  return info;
}

/**
 * Stops the specified agent by ID, returning its status info.
 * This will terminate the agent, and it cannot be resumed or restarted.
 * However, a new agent with the same name can be started at any time.
 */
export function stopAgent(agentId: string): AgentInfo {
  return agents.stop(agentId);
}

/**
 * Gets information about the specified agent.
 */
export function getAgentInfo(agentId: string): AgentInfo {
  return agents.getInfo(agentId);
}

/**
 * List all agents, except those that have been fully terminated.
 */
export function ListAgents(): AgentInfo[] {
  return agents.listAll();
}

/**
 * Returns the current count of the specified agent.
 */
export function getCount(agentId: string): i32 {
  const count = agents.sendMessage(agentId, "count");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

/**
 * Increments the count of the specified agent by 1 and returns the new count.
 */
export function updateCount(agentId: string): i32 {
  const count = agents.sendMessage(agentId, "increment");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

/**
 * Increments the count of the specified agent by the specified quantity.
 * This is an asynchronous operation and does not return a value.
 */
export function updateCountAsync(agentId: string, qty: i32): void {
  agents.sendMessageAsync(agentId, "increment", qty.toString());
}

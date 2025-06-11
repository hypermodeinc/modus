/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { AgentStatus } from "./enums";
import * as utils from "./utils";
import { JSON } from "json-as";

const agents = new Map<string, Agent>();
let activeAgent: Agent | null = null;
let activeAgentId: string | null = null;

/**
 * Base class for all agents.
 * Agents are long-running processes that can receive messages and maintain state.
 */
export abstract class Agent {
  /**
   * The unique ID of the agent instance.
   * This is set by the Modus Runtime when the agent is activated.
   */
  get id(): string {
    // Only active agents have an ID, and only one agent can be active at a time
    // in a single wasm module instance
    return activeAgentId || "";
  }

  /**
   * The name of the agent.
   * This should be unique across all agents in the module.
   * This is required to be implemented by the agent subclass.
   */
  abstract get name(): string;

  /**
   * Returns the serialized state of the agent.
   * This is required to be implemented by the agent subclass.
   */
  abstract getState(): string | null;

  /**
   * Sets the state of the agent from a serialized string.
   * This is required to be implemented by the agent subclass.
   */
  abstract setState(data: string | null): void;

  /**
   * Called when the agent is started.
   * Override this method to perform any initialization.
   */
  onInitialize(): void {}

  /**
   * Called when the agent is suspended.
   * Override this method if you want to do anything special when the agent is suspended,
   * such as sending a notification.
   * Note that you do not need to save the internal state of the agent here,
   * as that is handled automatically.
   */
  onSuspend(): void {}

  /**
   * Called when the agent is resumed.
   * Override this method if you want to do anything special when the agent is resumed,
   * such as sending a notification.
   * Note that you do not need to resume the internal state of the agent here,
   * as that is handled automatically.
   */
  onResume(): void {}

  /**
   * Called when the agent is terminated.
   * Override this method to send or save any final data.
   * Note that once an agent is terminated, it cannot be resumed.
   */
  onTerminate(): void {}

  /**
   * Called when the agent receives a message.
   * Override this method to handle incoming messages.
   */
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onReceiveMessage(msgName: string, data: string | null): string | null {
    return null;
  }

  /**
   * Publishes an event from this agent to any subscribers.
   */
  publishEvent(event: AgentEvent): void {
    const createdAt = new Date(Date.now()).toISOString();

    const data = JSON.stringify(event);
    hostPublishEvent(this.id, event.eventName, data, createdAt);
  }
}

/**
 * Registers an agent so it can be used.
 * This should be called from the module's entry point, not from any agent or function.
 */
export function registerAgent<T extends Agent>(): void {
  const agent = instantiate<T>();
  agents.set(agent.name, agent);
}

// @ts-expect-error: decorator
@external("modus_agents", "startAgent")
declare function hostStartAgent(agentName: string): AgentInfo;

// @ts-expect-error: decorator
@external("modus_agents", "stopAgent")
declare function hostStopAgent(agentId: string): AgentInfo;

// @ts-expect-error: decorator
@external("modus_agents", "getAgentInfo")
declare function hostGetAgentInfo(agentId: string): AgentInfo;

// @ts-expect-error: decorator
@external("modus_agents", "listAgents")
declare function hostListAgents(): AgentInfo[];

// @ts-expect-error: decorator
@external("modus_agents", "publishEvent")
declare function hostPublishEvent(
  agentId: string,
  eventName: string,
  eventData: string | null,
  createdAt: string,
): void;

/**
 * Starts an agent with the given name.
 * This can be called from any user code, such as function or another agent's methods.
 */
export function startAgent(name: string): AgentInfo {
  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  const info = hostStartAgent(name);
  if (utils.resultIsInvalid(info)) {
    throw new Error(`Failed to start agent ${name}.`);
  }
  return info;
}

/**
 * Stops an agent with the given ID.
 * This will terminate the agent, and it cannot be resumed or restarted.
 * This can be called from any user code, such as function or another agent's methods.
 */
export function stopAgent(agentId: string): AgentInfo {
  if (agentId == "") {
    throw new Error("Agent ID cannot be empty.");
  }
  const info = hostStopAgent(agentId);
  if (utils.resultIsInvalid(info)) {
    throw new Error(`Failed to stop agent ${agentId}.`);
  }
  return info;
}

/**
 * Gets information about an agent with the given ID.
 */
export function getAgentInfo(agentId: string): AgentInfo {
  if (agentId == "") {
    throw new Error("Agent ID cannot be empty.");
  }
  const info = hostGetAgentInfo(agentId);
  if (utils.resultIsInvalid(info)) {
    throw new Error(`Failed to get info for agent ${agentId}.`);
  }
  return info;
}

/**
 * Returns a list of all agents, except those that have been fully terminated.
 */
export function listAgents(): AgentInfo[] {
  const agents = hostListAgents();
  if (utils.resultIsInvalid(agents)) {
    throw new Error("Failed to list agents.");
  }
  return agents;
}

/**
 * The Modus Runtime will call this function to set the active agent.
 * This is called after the wasm module instance is loaded, regardless of whether
 * the agent is newly started or resumed from a suspended state.
 */
export function activateAgent(name: string, id: string): void {
  if (activeAgent != null) {
    throw new Error("Another agent is already active in this module instance.");
  }

  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  activeAgent = agents.get(name);
  activeAgentId = id;
}

/**
 * The Modus Runtime will call this function to get the state of an agent.
 * It is not intended to be called from user code.
 */
export function getAgentState(): string | null {
  if (!activeAgent) {
    throw new Error("No active agent to get state for.");
  }
  return activeAgent!.getState();
}

/**
 * The Modus Runtime will call this function to set the state of an agent.
 * It is not intended to be called from user code.
 */
export function setAgentState(data: string | null): void {
  if (!activeAgent) {
    throw new Error("No active agent to set state for.");
  }
  activeAgent!.setState(data);
}

/**
 * The Modus Runtime will call this function when an agent event occurs.
 * Events are part of the agent's lifecycle, managed by the Modus Runtime.
 */
export function handleEvent(action: string): void {
  if (!activeAgent) {
    throw new Error(`No active agent to process ${action} event action.`);
  }

  switch (action) {
    case "initialize":
      activeAgent!.onInitialize();
      break;

    case "suspend":
      activeAgent!.onSuspend();
      break;

    case "resume":
      activeAgent!.onResume();
      break;

    case "terminate":
      activeAgent!.onTerminate();
      break;

    default:
      throw new Error(`Unknown agent event action: ${action}`);
  }
}

/**
 * The Modus Runtime will call this function when an agent receives a message.
 * Messages are user-defined and can be sent from API functions or other agents.
 */
export function handleMessage(
  msgName: string,
  data: string | null,
): string | null {
  if (!activeAgent) {
    throw new Error("No active agent to handle message for.");
  }
  return activeAgent!.onReceiveMessage(msgName, data);
}

/**
 * Information about an agent instance.
 */
export class AgentInfo {
  /**
   * The unique ID of the agent instance.
   */
  readonly id: string;

  /**
   * The name of the agent.
   */
  readonly name: string;

  /**
   * The current status of the agent.
   */
  readonly status: AgentStatus;

  constructor(id: string, name: string, status: AgentStatus = "") {
    this.id = id;
    this.name = name;
    this.status = status;
  }
}

/**
 * Base class for agent events.
 * Custom agent events should extend this class.
 */
@json
export class AgentEvent {
  /**
   * The name of the event.
   */
  @omit
  readonly eventName: string;

  /**
   * Creates a new agent event.
   * @param eventName The name of the event.
   */
  constructor(eventName: string) {
    if (eventName == "") {
      throw new Error("Event name cannot be empty.");
    }
    this.eventName = eventName;
  }
}

/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { AgentStatus } from "./enums";

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
  onStart(): void {}

  /**
   * Called when the agent is stopped.
   * Override this method to perform any cleanup.
   */
  onStop(): void {}

  /**
   * Called when the agent is reloaded.
   * Override this method if you want to perform any actions when the agent is reloaded.
   * Reloading is done when the agent is updated, or when the agent is resumed after being suspended.
   */
  onReload(): void {}

  /**
   * Called when the agent receives a message.
   * Override this method to handle incoming messages.
   */
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onReceiveMessage(msgName: string, data: string | null): string | null {
    return null;
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
@external("modus_agents", "spawnAgentActor")
declare function hostSpawnAgentActor(agentName: string): AgentInfo;

/**
 * Starts an agent with the given name.
 * This can be called from any user code, such as function or another agent's methods.
 */
export function startAgent(name: string): AgentInfo {
  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  return hostSpawnAgentActor(name);
}

/**
 * The Modus Runtime will call this function to activate an agent.
 * It is not intended to be called from user code.
 */
export function activateAgent(name: string, id: string, reloading: bool): void {
  if (activeAgent != null) {
    throw new Error("Another agent is already active.");
  }

  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  activeAgent = agents.get(name);
  activeAgentId = id;

  if (reloading) {
    activeAgent!.onReload();
  } else {
    activeAgent!.onStart();
  }
}

/**
 * The Modus Runtime will call this function to shutdown an agent.
 * It is not intended to be called from user code.
 */
export function shutdownAgent(): void {
  if (!activeAgent) {
    throw new Error("No active agent to shut down.");
  }

  activeAgent!.onStop();
  activeAgent = null;
  activeAgentId = null;
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
 * The Modus Runtime will call this function when an agent receives a message.
 * It is not intended to be called from user code.
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

  constructor(
    id: string,
    name: string,
    status: AgentStatus = AgentStatus.Uninitialized,
  ) {
    this.id = id;
    this.name = name;
    this.status = status;
  }
}

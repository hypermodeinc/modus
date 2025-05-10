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

// @ts-expect-error: decorator
@external("modus_agents", "spawnAgentActor")
declare function hostSpawnAgentActor(agentName: string): AgentInfo;

export function registerAgent<T extends Agent>(): void {
  const agent = instantiate<T>();
  agents.set(agent.name, agent);
}

export function startAgent(name: string): AgentInfo {
  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  return hostSpawnAgentActor(name);
}

export function activateAgent(name: string, id: string, reloading: bool): void {
  if (!agents.has(name)) {
    throw new Error(`Agent ${name} not found.`);
  }

  if (activeAgent) {
    throw new Error("Another agent is already active.");
  }

  activeAgent = agents.get(name);
  activeAgentId = id;

  if (reloading) {
    activeAgent!.onReload();
  } else {
    activeAgent!.onStart();
  }
}

export function shutdownAgent(): void {
  if (!activeAgent) {
    throw new Error("No active agent to shut down.");
  }

  activeAgent!.onStop();
  activeAgent = null;
  activeAgentId = null;
}

export function getAgentState(): string | null {
  if (!activeAgent) {
    throw new Error("No active agent to get state for.");
  }
  return activeAgent!.getState();
}

export function setAgentState(data: string | null): void {
  if (!activeAgent) {
    throw new Error("No active agent to set state for.");
  }
  activeAgent!.setState(data);
}

export function handleMessage(
  msgName: string,
  data: string | null,
): string | null {
  if (!activeAgent) {
    throw new Error("No active agent to handle message for.");
  }
  return activeAgent!.onReceiveMessage(msgName, data);
}

export class AgentInfo {
  readonly id: string;
  readonly name: string;
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

export abstract class Agent {
  get id(): string {
    // Only active agents have an ID, and only one agent can be active at a time
    // in a single wasm module instance
    return activeAgentId || "";
  }

  abstract get name(): string;
  abstract getState(): string | null;
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

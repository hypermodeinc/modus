/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Duration } from "./enums";

const agents: Agent[] = [];

// @ts-expect-error: decorator
@external("modus_agents", "registerAgent")
declare function hostRegisterAgent(agentId: i32, name: string): void;

// @ts-expect-error: decorator
@external("modus_agents", "sendMessage")
declare function hostSendMessage(
  agentId: i32,
  msgName: string,
  data: string | null,
  timeout: i64,
): string | null;

export abstract class Agent {
  readonly id: i32;

  constructor() {
    this.id = agents.length + 1;
    agents.push(this);
  }

  /**
   * Sends a message to the agent, and waits for a response.
   */
  sendMessage(
    name: string,
    data: string | null = null,
    timeout: Duration = 10 * Duration.second,
  ): string | null {
    return hostSendMessage(this.id, name, data, timeout);
  }

  /**
   * Sends a message to the agent without waiting for a response.
   */
  sendMessageAsync(name: string, data: string | null = null): void {
    hostSendMessage(this.id, name, data, 0);
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
   * Called when the agent receives a message.
   * Override this method to handle incoming messages.
   */
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onReceiveMessage(name: string, data: string | null): string | null {
    return null;
  }
}

export function registerAgents(): void {
  for (let i = 0; i < agents.length; i++) {
    const agent = agents[i];
    hostRegisterAgent(agent.id, agent.name);
  }
}

export function startAgent(id: i32): void {
  const agent = getAgent(id);
  agent.onStart();
}

export function stopAgent(id: i32): void {
  const agent = getAgent(id);
  agent.onStop();
}

export function handleMessage(
  id: i32,
  name: string,
  data: string | null,
): string | null {
  const agent = getAgent(id);
  return agent.onReceiveMessage(name, data);
}

export function getAgentState(id: i32): string | null {
  const agent = getAgent(id);
  return agent.getState();
}

export function setAgentState(id: i32, data: string | null): void {
  const agent = getAgent(id);
  agent.setState(data);
}

function getAgent(id: i32): Agent {
  if (id < 1 || id > agents.length) {
    throw new Error(`Agent ${id} not found`);
  }
  return agents[id - 1];
}

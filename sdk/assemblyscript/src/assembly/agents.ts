/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Duration } from "./enums";

export {
  registerAgent as register,
  startAgent as start,
  stopAgent as stop,
  getAgentInfo as getInfo,
  listAgents as listAll,
} from "./agent";

// @ts-expect-error: decorator
@external("modus_agents", "sendMessage")
declare function hostSendMessage(
  agentId: string,
  msgName: string,
  data: string | null,
  timeout: i64,
): MessageResponse;

/**
 * Sends a message to the specified agent and waits for a response.
 * The data is optional and can be null.
 * The message is sent synchronously and the function will block until a response is received or the timeout is reached.
 * If there is no timeout set (or set to Duration.maxValue), it will wait indefinitely until a response is received,
 * or until calling function is cancelled.
 */
export function sendMessage(
  agentId: string,
  msgName: string,
  data: string | null = null,
  timeout: Duration = Duration.maxValue,
): string | null {
  if (timeout < 0) {
    throw new Error("Timeout must be a zero or positive value.");
  }
  if (agentId == "") {
    throw new Error("Agent ID cannot be empty.");
  }
  if (msgName == "") {
    throw new Error("Message name cannot be empty.");
  }

  const response = hostSendMessage(agentId, msgName, data, timeout);
  if (response == null) {
    throw new Error("Failed to send message to agent.");
  }
  if (response.error) {
    throw new Error(response.error!);
  }

  return response.data;
}

/**
 * Sends a message to the specified agent without waiting for a response.
 * The data is optional and can be null.
 * The message is sent asynchronously and the function will return immediately.
 */
export function sendMessageAsync(
  agentId: string,
  msgName: string,
  data: string | null = null,
): void {
  sendMessage(agentId, msgName, data, Duration.zero);
}

class MessageResponse {
  data: string | null = null;
  error: string | null = null;
}

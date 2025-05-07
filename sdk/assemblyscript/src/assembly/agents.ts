/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Duration } from "./enums";

export { registerAgent as register, startAgent as start } from "./agent";

// @ts-expect-error: decorator
@external("modus_agents", "sendMessage")
declare function hostSendMessage(
  agentId: string,
  msgName: string,
  data: string | null,
  timeout: i64,
): MessageResponse;

/**
 * Sends a message to an agent, and waits for a response.
 * The data is optional and can be null.
 * The timeout is in seconds and defaults to 10 seconds.
 */
export function sendMessage(
  agentId: string,
  msgName: string,
  data: string | null = null,
  timeout: Duration = 10 * Duration.second,
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
  return response.data;
}

/**
 * Sends a message to an agent without waiting for a response.
 * The data is optional and can be null.
 */
export function sendMessageAsync(
  agentId: string,
  msgName: string,
  data: string | null = null,
): void {
  sendMessage(agentId, msgName, data, Duration.zero);
}

class MessageResponse {
  data!: string | null;
}

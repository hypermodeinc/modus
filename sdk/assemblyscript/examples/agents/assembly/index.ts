/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { agents, AgentInfo } from "@hypermode/modus-sdk-as";
import { CounterAgent } from "./counterAgent";

agents.register<CounterAgent>();

export function startAgent(): AgentInfo {
  return agents.start("Counter");
}

export function updateCount(agentId: string): i32 {
  const count = agents.sendMessage(agentId, "increment");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

export function getCount(agentId: string): i32 {
  const count = agents.sendMessage(agentId, "count");
  if (count == null) {
    return 0;
  }
  return i32.parse(count);
}

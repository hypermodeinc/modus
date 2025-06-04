/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, mockImport, run } from "as-test";
import { JSON } from "json-as";
import { AgentStatus } from "../enums";
import { AgentInfo, listAgents } from "../agent";
import { agents } from "..";
import { TaskManagerAgent } from "./src/taskagent";

let start_agent_ret: AgentInfo | null = null;
let stop_agent_ret: AgentInfo | null = null;
let get_agent_info_ret: AgentInfo | null = null;
let list_agents_ret: AgentInfo[] = [];
let send_message_ret: MessageResponse | null = null;

mockImport("modus_agents.startAgent", (name: string): AgentInfo => {
  return start_agent_ret!;
});

mockImport("modus_agents.stopAgent", (name: string): AgentInfo => {
  return stop_agent_ret!;
});

mockImport("modus_agents.getAgentInfo", (name: string): AgentInfo => {
  return get_agent_info_ret!;
});

mockImport("modus_agents.listAgents", (): AgentInfo[] => {
  return list_agents_ret;
});

mockImport("modus_agents.sendMessage", (agentId: string, msgName: string, data: string | null, timeout: i64): MessageResponse | null => {
  return send_message_ret;
});

it("should serialize an AgentStatus using type aliases", () => {
  const status: AgentStatus = AgentStatus.Resuming;
  expect(JSON.stringify(status)).toBe('"' + AgentStatus.Resuming + '"');
});

it("should list current agents", () => {
  const agents = listAgents();
  expect(agents.length).toBe(0);
});

it("should register an agent", () => {
  agents.register<TaskManagerAgent>();
  // Not sure what to expect() here...
});

it("should start an agent", () => {
  start_agent_ret = new AgentInfo("d1086e837bkp4ltjm150", "TaskManager", "starting");
  const agent = agents.start("TaskManager");
  expect(agent.id).toBe("d1086e837bkp4ltjm150");
  expect(agent.name).toBe("TaskManager");
  expect(agent.status).toBe("starting");
});

it("should list current agents", () => {
  list_agents_ret = [new AgentInfo("d1086e837bkp4ltjm150", "TaskManager", "running")];
  const agents = listAgents();
  expect(agents.length).toBe(1);
  const agent = agents[0]!;
  expect(agent.id).toBe("d1086e837bkp4ltjm150");
  expect(agent.name).toBe("TaskManager");
  expect(agent.status).toBe("running");
});

it("should stop an agent", () => {
  stop_agent_ret = new AgentInfo("d1086e837bkp4ltjm150", "TaskManager", "terminated");
  const agent = agents.stop("TaskManager");
  expect(agent.id).toBe("d1086e837bkp4ltjm150");
  expect(agent.name).toBe("TaskManager");
  expect(agent.status).toBe("terminated");
});

it("should list current agents", () => {
  list_agents_ret = [];
  const agents = listAgents();
  expect(agents.length).toBe(0);
});

it("should add a task", () => {
  send_message_ret = new MessageResponse('Task added: "Sign up for Hypermode"');
  const res = agents.sendMessage("d1086e837bkp4ltjm150", "addTask", 'Sign up for Hypermode');
  expect(res).toBe('Task added: "Sign up for Hypermode"');
});

it("should complete a task", () => {
  send_message_ret = new MessageResponse('Task completed: "Sign up for Hypermode"');
  const res = agents.sendMessage("d1086e837bkp4ltjm150", "completeTask", 'Sign up for Hypermode');
  expect(res).toBe('Task completed: "Sign up for Hypermode"');
});

it("should list tasks", () => {
  send_message_ret = new MessageResponse('[x] Sign up for Hypermode');
  const res = agents.sendMessage("d1086e837bkp4ltjm150", "list");
  expect(res).toBe('[x] Sign up for Hypermode');
});

it ("should get statistics for tasks", () => {
  send_message_ret = new MessageResponse('Tasks: 1/1 completed');
  const res = agents.sendMessage("d1086e837bkp4ltjm150", "stats");
  expect(res).toBe('Tasks: 1/1 completed');
});

run();

@json
class MessageResponse {
  data: string | null;
  constructor(data: string | null) {
    this.data = data;
  }
}

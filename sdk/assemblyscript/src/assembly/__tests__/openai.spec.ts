/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, run } from "as-test";
import { JSON } from "json-as";

import {
  SystemMessage,
  UserMessage,
  AssistantMessage,
  RequestMessage,
  parseMessages,
} from "../../models/openai/chat";

it("should round-trip chat messages", () => {
  const msgs: RequestMessage[] = [
    new SystemMessage("You are a helpful assistant."),
    new UserMessage("What is the capital of France?"),
    new AssistantMessage("The capital of France is Paris."),
  ];

  const data = JSON.stringify(msgs);
  const parsedMsgs = parseMessages(data);

  expect(parsedMsgs.length).toBe(msgs.length);
  for (let i = 0; i < msgs.length; i++) {
    expect(msgs[i].role).toBe(parsedMsgs[i].role);
  }

  const roundTrip = JSON.stringify(parsedMsgs);
  expect(roundTrip).toBe(data);
});

it("should parse oddly formatted message", () => {
  const data = `
{
    "role"  :  \t\n "assistant",
  "content" : "The capital of France is Paris."
}`;

  const p1 = JSON.parse<JSON.Obj>(data);
  expect(p1.get("role")!.get<string>()).toBe("assistant");
  expect(p1.get("content")!.get<string>()).toBe(
    "The capital of France is Paris.",
  );

  const p2 = JSON.parse<AssistantMessage<string>>(data);
  expect(p2.role).toBe("assistant");
  expect(p2.content).toBe("The capital of France is Paris.");
});

it("should parse multiple messages of the same type", () => {
  const json = JSON.stringify([
    new UserMessage("First"),
    new UserMessage("Second"),
    new UserMessage("Third"),
  ]);

  const parsed = parseMessages(json);
  expect(parsed.length).toBe(3);
  for (let i = 0; i < parsed.length; i++) {
    expect(parsed[i].role).toBe("user");
  }
});

it("should correctly parse mixed message roles", () => {
  const data = JSON.stringify<RequestMessage[]>([
    new UserMessage("foo"),
    new AssistantMessage<string>("bar"),
    new SystemMessage("baz"),
  ]);
  const parsed = parseMessages(data);
  expect(parsed[0].role).toBe("user");
  expect(parsed[1].role).toBe("assistant");
  expect(parsed[2].role).toBe("system");
});

run();

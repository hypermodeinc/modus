/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { expect, it, log, run } from "as-test";
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

run();

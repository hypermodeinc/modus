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

import {
  SystemMessage,
  UserMessage,
  AssistantMessage,
  RequestMessage,
  parseMessages,
  OpenAIChatModel,
} from "../../models/openai/chat";
import { models } from "..";
import { Model, ModelInfo } from "../models";

let _get_model_name: string = "";
let _get_model_info_out: ModelInfo | null = null;
let _invoke_input: string = "";
let _invoke_model_out: string | null = null;

mockImport("modus_models.getModelInfo", (modelName: string): ModelInfo => {
  _get_model_name = modelName;
  return _get_model_info_out!;
});

const modelInvokerFn = (modelName: string, input: string): string | null => {
  _invoke_model_name = modelName;
  _invoke_input = input;

  return _invoke_model_out;
};

mockImport<(modelName: string, input: string) => string | null>(
  "modus_models.invokeModel",
  modelInvokerFn,
);

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

it("should invoke mocked OpenAI model", () => {
  const modelName = "meta-llama/llama-3.2-3b-instruct";
  _get_model_info_out = new ModelInfo("text-generator", modelName);
  const model = models.getModel<OpenAIChatModel>("text-generator");
  expect(_get_model_name).toBe("text-generator");

  const input = model.createInput([
    new SystemMessage<string>("Your name is bob"),
    new UserMessage<string>("What's your name?"),
  ]);

  expect(JSON.stringify(input)).toBe(
    '{"model":"' +
      modelName +
      '","messages":[{"role":"system","content":"Your name is bob"},{"role":"user","content":"What\'s your name?"}]}',
  );

  input.temperature = 0.6;

  Model.invoker = modelInvokerFn;
  _invoke_model_out =
    '{"id":"chatcmpl-88bcf38f-a5da-4465-a824-76ed161af95c","object":"chat.completion","created":1748914205,"model":"meta-llama/llama-3.2-3b-instruct","choices":[{"index":0,"message":{"role":"assistant","reasoning_content":null,"content":"My name is Bob.","tool_calls":[]},"logprobs":null,"finish_reason":"stop","stop_reason":null}],"usage":{"prompt_tokens":44,"total_tokens":50,"completion_tokens":6,"prompt_tokens_details":null},"prompt_logprobs":null}';
  const output = model.invoke(input);

  expect(_invoke_input).toBe(
    '{"temperature":0.6,"model":"meta-llama/llama-3.2-3b-instruct","messages":[{"role":"system","content":"Your name is bob"},{"role":"user","content":"What\'s your name?"}]}',
  );

  expect(output.choices[0].message.content).toBe("My name is Bob.");
});

// it("should create prompt message for quote summary generation", () => {
//   const quote = "Every person, all the events of your life are there because you have drawn them there. What you choose to do with them is up to you.";
//   const author = "Richard Bach";
//   const instruction = "Provide a brief, insightful summary that captures the essence and meaning of the quote in 1-2 sentences.";
//   const prompt = "Quote: " + quote + " - " + author;

// })

run();

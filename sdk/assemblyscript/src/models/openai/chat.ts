/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Model } from "../../assembly/models";
import { JSON } from "json-as";

/**
 * Provides input and output types that conform to the OpenAI Chat API.
 *
 * Reference: https://platform.openai.com/docs/api-reference/chat
 */
export class OpenAIChatModel extends Model<OpenAIChatInput, OpenAIChatOutput> {
  /**
   * Creates an input object for the OpenAI Chat API.
   *
   * @param messages: An array of messages to send to the chat model.
   * @returns An input object that can be passed to the `invoke` method.
   */
  createInput(messages: Message[]): OpenAIChatInput {
    const model = this.info.fullName;
    return <OpenAIChatInput>{ model, messages };
  }
}

/**
 * The input object for the OpenAI Chat API.
 */
@json
export class OpenAIChatInput {
  /**
   * The name of the model to use for the chat.
   * Must be the exact string expected by the model provider.
   * For example, "gpt-3.5-turbo".
   *
   * @remarks
   * This field is automatically set by the `createInput` method when creating this object.
   * It does not need to be set manually.
   */
  model!: string;

  /**
   * An array of messages to send to the chat model.
   */
  messages!: Message[];

  /**
   * Number between `-2.0` and `2.0`.
   *
   * Positive values penalize new tokens based on their existing frequency in the text so far,
   * decreasing the model's likelihood to repeat the same line verbatim.
   *
   * @default 0.0
   */
  @alias("frequency_penalty")
  @omitif("this.frequencyPenalty == 0.0")
  frequencyPenalty: f64 = 0.0;

  /**
   * Modifies the likelihood of specified tokens appearing in the completion.
   *
   * Accepts an object that maps tokens (specified by their token ID in the tokenizer) to an associated bias value from -100 to 100.
   * Mathematically, the bias is added to the logits generated by the model prior to sampling. The exact effect will vary per model,
   * but values between -1 and 1 should decrease or increase likelihood of selection; values like -100 or 100 should result in a ban
   * or exclusive selection of the relevant token.
   */
  @alias("logit_bias")
  @omitnull()
  logitBias: Map<string, f64> | null = null;

  /**
   * Whether to return log probabilities of the output tokens or not,
   *
   * If true, returns the log probabilities of each output token returned in the content of message.
   *
   * @default false
   */
  @omitif("this.logprobs == false")
  logprobs: bool = false;

  /**
   * An integer between 0 and 20 specifying the number of most likely tokens to return at each token position,
   * each with an associated log probability. `logprobs` must be set to `true` if this parameter is used.
   */
  @alias("top_logprobs")
  @omitif("this.logprobs == false")
  topLogprobs: i32 = 0;

  /**
   * The maximum number of tokens to generate in the chat completion.
   *
   * @default 4096
   */
  @alias("max_tokens")
  @omitif("this.maxTokens == 4096")
  maxTokens: i32 = 4096; // TODO: make this an `i32 | null` when supported

  /**
   * The number of completions to generate for each prompt.
   */
  @omitif("this.n == 1")
  n: i32 = 1;

  /**
   *
   * Number between `-2.0` and `2.0`.
   *
   * Positive values penalize new tokens based on whether they appear in the text so far,
   * increasing the model's likelihood to talk about new topics.
   *
   * @default 0.0
   */
  @alias("presence_penalty")
  @omitif("this.presencePenalty == 0.0")
  presencePenalty: f64 = 0.0;

  /**
   * Specifies the format for the response.
   *
   * If set to `ResponseFormat.Json`, the response will be a JSON object.
   *
   * If set to `ResponseFormat.JsonSchema`, the response will be a JSON object
   * that conforms to the provided JSON schema.
   *
   * @default ResponseFormat.Text
   */
  @alias("response_format")
  @omitif("this.responseFormat.type == 'text'")
  responseFormat: ResponseFormat = ResponseFormat.Text;

  /**
   * If specified, the model will make a best effort to sample deterministically, such that
   * repeated requests with the same seed and parameters should return the same result.
   * Determinism is not guaranteed, and you should use the `systemFingerprint` response
   * parameter to monitor changes in the backend.
   */
  @omitif("this.seed == -1")
  seed: i32 = -1; // TODO: make this an `i32 | null` when supported

  /**
   * Specifies the latency tier to use for processing the request.
   */
  @alias("service_tier")
  @omitnull()
  serviceTier: ServiceTier | null = null;

  /**
   * Up to 4 sequences where the API will stop generating further tokens.
   */
  @omitnull()
  stop: string[] | null = null;

  // stream: bool = false;

  // @omitif("this.stream == false")
  // @alias("stream_options")
  // streamOptions: StreamOptions | null = null;

  /**
   * A number between `0.0` and `2.0` that controls the sampling temperature.
   *
   * Higher values like `0.8` will make the output more random, while lower values like `0.2` will make
   * it more focused and deterministic.
   *
   * We generally recommend altering this or `topP` but not both.
   *
   * @default 1.0
   */
  @omitif("this.temperature == 1.0")
  temperature: f64 = 1.0;

  /**
   * An alternative to sampling with temperature, called nucleus sampling, where the model
   * considers the results of the tokens with `topP` probability mass.
   *
   * For example, `0.1` means only the tokens comprising the top 10% probability mass are considered.
   *
   * We generally recommend altering this or `temperature` but not both.
   */
  @alias("top_p")
  @omitif("this.topP == 1.0")
  topP: f64 = 1.0;

  /**
   * A list of tools the model may call. Currently, only functions are supported as a tool.
   * Use this to provide a list of functions the model may generate JSON inputs for.
   * A max of 128 functions are supported.
   */
  @omitnull()
  tools: Tool[] | null = null;

  /**
   * Controls which (if any) tool is called by the model.
   * - `"none"` means the model will not call any tool and instead generates a message.
   * - `"auto"` means the model can pick between generating a message or calling one or more tools.
   * - `"required"` means the model must call one or more tools.
   * - Specifying a particular tool via `{"type": "function", "function": {"name": "my_function"}}` forces the model to call that tool.
   *
   * `none` is the default when no tools are present. `auto` is the default if tools are present.
   */
  @alias("tool_choice")
  @omitnull()
  toolChoice: string | null = null; // TODO: Make this work with a custom tool object

  /**
   * Whether to enable parallel function calling during tool use.
   *
   * @default true
   */
  @alias("parallel_tool_calls")
  @omitif("this.parallelToolCalls == true || !this.tools || this.tools!.length == 0")
  parallelToolCalls: bool = true;

  /**
   * The user ID to associate with the request.
   * If not specified, the request will be anonymous.
   * See https://platform.openai.com/docs/guides/safety-best-practices/end-user-ids
   */
  @omitnull()
  user: string | null = null;
}

/**
 * The OpenAI service tier used to process the request.
 */
// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace ServiceTier {
  /**
   * The OpenAI system will utilize scale tier credits until they are exhausted.
   */
  export const Auto = "auto";

  /**
   * The request will be processed using the default OpenAI service tier with a lower
   * uptime SLA and no latency guarantee.
   */
  export const Default = "default";
}
export type ServiceTier = string;

/**
 * The output object for the OpenAI Chat API.
 */
@json
export class OpenAIChatOutput {
  /**
   * A unique identifier for the chat completion.
   */
  id!: string;

  /**
   * The name of the output object type returned by the API.
   * Always `"chat.completion"`.
   */
  object!: string;

  /**
   * A list of chat completion choices. Can be more than one if `n` is greater than 1 in the input options.
   */
  choices!: Choice[];

  /**
   * The Unix timestamp (in seconds) of when the chat completion was created.
   */
  created!: i32;

  /**
   * The model used for the chat completion.
   * In most cases, this will match the requested `model` field in the input.
   */
  model!: string;

  /**
   * The service tier used for processing the request.
   *
   * This field is only included if the `serviceTier` parameter is specified in the request.
   */
  @alias("service_tier")
  serviceTier: string | null = null;

  /**
   * This fingerprint represents the OpenAI backend configuration that the model runs with.
   *
   * Can be used in conjunction with the seed request parameter to understand when backend changes
   * have been made that might impact determinism.
   */
  @alias("system_fingerprint")
  systemFingerprint!: string;

  /**
   * The usage statistics for the request.
   */
  usage!: Usage;
}

type JsonSchemaFunction = (jsonSchema: string) => ResponseFormat;

/**
 * An object specifying the format that the model must output.
 */
@json
export class ResponseFormat {
  /**
   * The type of response format.  Must be one of `"text"` or `"json_object"`.
   */
  readonly type!: string;

  /**
   * The JSON schema to use for the response format.
   */
  @omitnull()
  @alias("json_schema")
  readonly jsonSchema: JSON.Raw | null = null;

  /**
   * Instructs the model to output the response as a JSON object.
   *
   * @remarks
   * You must also instruct the model to produce JSON yourself via a system or user message.
   *
   * Additionally, if you need an array you must ask for an object that wraps the array,
   * because the model will not reliably produce arrays directly (ie., there is no `json_array` option).
   */
  static Json: ResponseFormat = { type: "json_object", jsonSchema: null };

  /**
   * Enables Structured Outputs which guarantees the model will match your supplied JSON schema.
   *
   * See https://platform.openai.com/docs/guides/structured-outputs
   */
  static JsonSchema: JsonSchemaFunction = (
    jsonSchema: string,
  ): ResponseFormat => {
    return {
      type: "json_schema",
      jsonSchema: jsonSchema,
    };
  };

  /**
   * Instructs the model to output the response as a plain text string.
   *
   * @remarks
   * This is the default response format.
   */
  static Text: ResponseFormat = { type: "text", jsonSchema: null };
}

// @json
// export class StreamOptions {

//   @omitif("this.includeUsage == false")
//   @alias("include_usage")
//   includeUsage: bool = false;
// }

/**
 * A tool object that the model may call.
 */
@json
export class Tool {
  /**
   * The type of the tool. Currently, only `"function"` is supported.
   *
   * @default "function"
   */
  type: string = "function";

  /**
   * The definition of the function.
   */
  function!: FunctionDefinition;
}

/**
 * The definition of a function that can be called by the model.
 */
@json
export class FunctionDefinition {
  /**
   * The name of the function to be called.
   *
   * Must be a-z, A-Z, 0-9, or contain underscores and dashes, with a maximum length of 64.
   */
  name!: string;

  /**
   * An optional description of what the function does, used by the model to choose when and how to call the function.
   */
  @omitnull()
  description: string | null = null;

  /**
   * Whether to enable strict schema adherence when generating the function call.
   * If set to true, the model will follow the exact schema defined in the parameters field.
   *
   * See https://platform.openai.com/docs/guides/function-calling
   *
   * @remarks
   * In order to guarantee strict schema adherence, disable parallel function calls
   * by setting {@link OpenAIChatInput.parallelToolCalls}=false.
   *
   * See https://platform.openai.com/docs/guides/function-calling/parallel-function-calling-and-structured-outputs
   *
   * @default false
   */
  @omitif("this.strict == false")
  strict: bool = false;

  /**
   * The parameters the functions accepts, described as a JSON Schema object.
   *
   * See https://platform.openai.com/docs/guides/function-calling
   */
  @omitnull()
  parameters: JSON.Raw | null = null;
}

/**
 * A tool call object that the model may generate.
 */
@json
export class ToolCall {
  /**
   * The ID of the tool call.
   */
  id!: string;

  /**
   * The type of the tool. Currently, only `"function"` is supported.
   *
   * @default "function"
   */
  type: string = "function";

  /**
   * The function that the model called.
   */
  function!: FunctionCall;
}

/**
 * A function call object that the model may generate.
 */
@json
export class FunctionCall {
  /**
   * The name of the function to call.
   */
  name!: string;

  /**
   * The arguments to call the function with, as generated by the model in JSON format.
   *
   * @remarks
   * Note that the model does not always generate valid JSON, and may hallucinate parameters not
   * defined by your function schema. Validate the arguments in your code before calling your function.
   */
  arguments!: string;
}

/**
 * The usage statistics for the request.
 */
@json
export class Usage {
  /**
   * The number of completion tokens used in the response.
   */
  @alias("completion_tokens")
  completionTokens!: i32;

  /**
   * The number of prompt tokens used in the request.
   */
  @alias("prompt_tokens")
  promptTokens!: i32;

  /**
   * The total number of tokens used in the request and response.
   */
  @alias("total_tokens")
  totalTokens!: i32;
}

/**
 * A completion choice object returned in the response.
 */
@json
export class Choice {
  /**
   * The reason the model stopped generating tokens.
   *
   * Possible values are:
   * - `"stop"` if the model hit a natural stop point or a provided stop sequence
   * - `"length"` if the maximum number of tokens specified in the request was reached
   * - `"content_filter"` if content was omitted due to a flag from our content filters
   * - `"tool_calls"` if the model called a tool
   */
  @alias("finish_reason")
  finishReason!: string;

  /**
   * The index of the choice in the list of choices.
   */
  index!: i32;

  /**
   * A chat completion message generated by the model.
   */
  message!: CompletionMessage;

  /**
   * Log probability information for the choice.
   */
  logprobs!: Logprobs | null;
}

/**
 * Log probability information for a choice.
 */
@json
export class Logprobs {
  /**
   * A list of message content tokens with log probability information.
   */
  content: LogprobsContent[] | null = null;
}

/**
 * Log probability information for a message content token.
 */
@json
export class LogprobsContent {
  /**
   * The token.
   */
  token!: string;

  /**
   * The log probability of this token, if it is within the top 20 most likely tokens.
   * Otherwise, the value `-9999.0` is used to signify that the token is very unlikely.
   */
  logprob!: f64;

  /**
   * A list of integers representing the UTF-8 bytes representation of the token.
   *
   * Useful in instances where characters are represented by multiple tokens and their byte
   * representations must be combined to generate the correct text representation.
   * Can be null if there is no bytes representation for the token.
   */
  bytes!: u8[] | null; // TODO: verify this works

  /**
   * List of the most likely tokens and their log probability, at this token position.
   * In rare cases, there may be fewer than the number of requested `topLogprobs` returned.
   */
  @alias("top_logprobs")
  topLogprobs!: TopLogprobsContent[]; // TODO: verify this works
}

/**
 * Log probability information for the most likely tokens at a given position.
 */
@json
export class TopLogprobsContent {
  /**
   * The token.
   */
  token!: string;

  /**
   * The log probability of this token, if it is within the top 20 most likely tokens.
   * Otherwise, the value `-9999.0` is used to signify that the token is very unlikely.
   */
  logprob!: f64;

  /**
   * A list of integers representing the UTF-8 bytes representation of the token.
   * Useful in instances where characters are represented by multiple tokens and their byte
   * representations must be combined to generate the correct text representation.
   * Can be null if there is no bytes representation for the token.
   */
  bytes!: u8[] | null; // TODO: verify this works
}

/**
 * A message object that can be sent to the chat model.
 */
@json
abstract class Message {
  /**
   * Creates a new message object.
   *
   * @param role The role of the author of this message.
   * @param content The contents of the message.
   */
  constructor(role: string, content: string) {
    this._role = role;
    this.content = content;
  }


  @alias("role")
  protected _role: string;

  /**
   * The role of the author of this message.
   */
  get role(): string {
    return this._role;
  }

  /**
   * The contents of the message.
   */
  content: string;
}

/**
 * A system message.
 */
@json
export class SystemMessage extends Message {
  /**
   * Creates a new system message object.
   *
   * @param content The contents of the message.
   */
  constructor(content: string) {
    super("system", content);
  }

  /**
   * An optional name for the participant.
   * Provides the model information to differentiate between participants of the same role.
   */
  @omitnull()
  name: string | null = null;
}

/**
 * A user message.
 */
@json
export class UserMessage extends Message {
  /**
   * Creates a new user message object.
   *
   * @param content The contents of the message.
   */
  constructor(content: string) {
    super("user", content);
  }

  /**
   * An optional name for the participant.
   * Provides the model information to differentiate between participants of the same role.
   */
  @omitnull()
  name: string | null = null;
}

/**
 * An assistant message.
 */
@json
export class AssistantMessage extends Message {
  /**
   * Creates a new assistant message object.
   *
   * @param content The contents of the message.
   */
  constructor(content: string) {
    super("assistant", content);
  }

  /**
   * An optional name for the participant.
   * Provides the model information to differentiate between participants of the same role.
   */
  @omitnull()
  name: string | null = null;

  /**
   * The tool calls generated by the model, such as function calls.
   */
  @alias("tool_calls")
  @omitif("this.toolCalls.length == 0")
  toolCalls: ToolCall[] = [];
}

/**
 * A tool message.
 */
@json
export class ToolMessage extends Message {
  /**
   * Creates a new tool message object.
   *
   * @param content The contents of the message.
   */
  constructor(content: string, toolCallId: string) {
    super("tool", content);
    this.toolCallId = toolCallId;
  }

  /**
   * Tool call that this message is responding to.
   */
  @alias("tool_call_id")
  toolCallId!: string;
}

/**
 * A chat completion message generated by the model.
 */
@json
export class CompletionMessage extends Message {
  /**
   * Creates a new completion message object.
   *
   * @param role The role of the author of this message.
   * @param content The contents of the message.
   */
  constructor(role: string, content: string) {
    super(role, content);
  }

  /**
   * The refusal message generated by the model.
   */
  refusal: string | null = null;

  /**
   * The tool calls generated by the model, such as function calls.
   */
  @alias("tool_calls")
  toolCalls: ToolCall[] = [];
}

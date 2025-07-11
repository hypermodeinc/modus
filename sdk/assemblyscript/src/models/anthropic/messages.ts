/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Model } from "../../assembly/models";
import { JSON } from "json-as";

/**
 * Provides input and output types that conform to the Anthropic Messages API.
 *
 * Reference: https://docs.anthropic.com/en/api/messages
 */
export class AnthropicMessagesModel extends Model<
  AnthropicMessagesInput,
  AnthropicMessagesOutput
> {
  /**
   * Creates an input object for the Anthropic Messages API.
   *
   * @param messages: An array of messages to send to the model.
   * Note that if you want to include a system prompt, you can use the top-level system parameter —
   * there is no "system" role for input messages in the Messages API.
   * @returns An input object that can be passed to the `invoke` method.
   */
  createInput(messages: Message[]): AnthropicMessagesInput {
    const model = this.info.fullName;
    return <AnthropicMessagesInput>{ model, messages };
  }
}

/**
 * A message object that can be sent to the model.
 */
@json
export class Message {
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
   * For now it can only be a string, even though Anthropic supports more complex types.
   */
  content: string; // TODO: support more complex types
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
}

/**
 * The input object for the Anthropic Messages API.
 */
@json
export class AnthropicMessagesInput {
  /**
   * The model that will complete your prompt.
   * Must be the exact string expected by the model provider.
   * For example, "claude-3-5-sonnet-20240620".
   *
   * See [models](https://docs.anthropic.com/en/docs/models-overview) for additional
   * details and options.
   *
   * @remarks
   * This field is automatically set by the `createInput` method when creating this object.
   * It does not need to be set manually.
   */
  model!: string;

  /**
   * Input messages.
   *
   * We do not currently support image content blocks, which are available starting with
   * Claude 3 models. This will be added in a future release.
   *
   * Note that if you want to include a
   * [system prompt](https://docs.anthropic.com/en/docs/system-prompts), you can use
   * the top-level `system` parameter — there is no `"system"` role for input
   * messages in the Messages API.
   */
  messages!: Message[];

  /**
   * The maximum number of tokens to generate before stopping.
   *
   * Different models have different maximum values for this parameter. See
   * [models](https://docs.anthropic.com/en/docs/models-overview) for details.
   *
   * @default 4096
   */
  @alias("max_tokens")
  maxTokens: i32 = 4096;

  /**
   * A `Metadata` object describing the request.
   */
  @omitnull()
  metadata: Metadata | null = null;

  /**
   * Custom text sequences that will cause the model to stop generating.
   */
  @alias("stop_sequences")
  @omitnull()
  stopSequences: string[] | null = null;

  /**
   * Streaming is not currently supported.
   *
   * @default false
   */
  @alias("stream")
  @omitif((self: AnthropicMessagesInput) => self._stream == false)
  private _stream: boolean = false;

  /**
   * System prompt.
   *
   * A system prompt is a way of providing context and instructions to Claude, such
   * as specifying a particular goal or role. See [guide to system prompts](https://docs.anthropic.com/en/docs/system-prompts).
   */
  @omitnull()
  system: string | null = null;

  /**
   * A number between `0.0` and `1.0` that controls the randomness injected into the response.
   *
   * It is recommended to use `temperature` closer to `0.0`
   * for analytical / multiple choice, and closer to `1.0` for creative tasks.
   *
   * Note that even with `temperature` of `0.0`, the results will not be fully
   * deterministic.
   *
   * @default 1.0
   */
  @omitif((self: AnthropicMessagesInput) => self.temperature == 1.0)
  temperature: f64 = 1.0;

  /**
   * How the model should use the provided tools.
   *
   * Use either `ToolChoiceAuto`, `ToolChoiceAny`, or `ToolChoiceTool(name: string)`.
   */
  @alias("tool_choice")
  @omitnull()
  toolChoice: ToolChoice | null = null;

  /**
   * Definitions of tools that the model may use.
   *
   * Tools can be used for workflows that include running client-side tools and
   * functions, or more generally whenever you want the model to produce a particular
   * JSON structure of output.
   *
   * See Anthropic's [guide](https://docs.anthropic.com/en/docs/tool-use) for more details.
   */
  @omitnull()
  tools: Tool[] | null = null;

  /**
   * Only sample from the top K options for each subsequent token.
   *
   * Recommended for advanced use cases only. You usually only need to use
   * `temperature`.
   */
  @alias("top_k")
  @omitif((self: AnthropicMessagesInput) => self.topK == -1)
  topK: i64 = -1; // The default value of top_k is not specified in the API docs

  /**
   * Use nucleus sampling.
   *
   * You should either alter `temperature` or `top_p`, but not both.
   *
   * Recommended for advanced use cases only. You usually only need to use
   * `temperature`.
   */
  @alias("top_p")
  @omitif((self: AnthropicMessagesInput) => self.topP == 0.999)
  topP: f64 = 0.999;
}


@json
export class Metadata {
  /**
   * An external identifier for the user who is associated with the request.
   */
  @alias("user_id")
  userId: string | null = null;
}

/**
 * A tool object that the model may call.
 */
@json
export class Tool {
  /**
   * Name of the tool.
   */
  name!: string;

  /**
   * [JSON schema](https://json-schema.org/) for this tool's input.
   *
   * This defines the shape of the `input` that your tool accepts and that the model
   * will produce.
   */
  @alias("input_schema")
  inputSchema!: JSON.Raw;

  /**
   * Optional, but strongly-recommended description of the tool.
   */
  @omitnull()
  description: string | null = null;
}


@json
export class ToolChoice {
  constructor(type: string, name: string | null = null) {
    this._type = type;
    this._name = name;
  }


  @alias("type")
  protected _type: string;

  /**
   * The name of the tool to use.
   */
  @alias("name")
  @omitnull()
  protected _name: string | null = null;
}

/**
 * The model will automatically decide whether to use tools.
 */
export const ToolChoiceAuto = new ToolChoice("auto");

/**
 * The model will use any available tools.
 */
export const ToolChoiceAny = new ToolChoice("any");

/**
 * The model will use the specified tool.
 */
export const ToolChoiceTool = (name: string): ToolChoice =>
  new ToolChoice("tool", name);

/**
 * The output object for the Anthropic Messages API.
 */
@json
export class AnthropicMessagesOutput {
  /**
   * Unique object identifier.
   */
  id!: string;

  /**
   * Content generated by the model.
   *
   */
  content!: ContentBlock[];

  /**
   * The model that handled the request.
   */
  model!: string;

  /**
   * Conversational role of the generated message.
   *
   * This will always be `"assistant"`.
   */
  role!: "assistant";

  /**
   * The reason that the model stopped.
   */
  @alias("stop_reason")
  stopReason!: string;

  /**
   * Which custom stop sequence was generated, if any.
   *
   * This value will be a non-null string if one of your custom stop sequences was
   * generated.
   */
  @alias("stop_sequence")
  stopSequence: string | null = null;

  /**
   * Object type. This is always `"message"`.
   */
  type!: "message";

  /**
   * Billing and rate-limit usage.
   */
  usage!: Usage;
}


@json
export class ContentBlock {
  type!: string;

  // Text block
  @omitnull()
  text: string | null = null;

  // Tool use block
  @omitnull()
  id: string | null = null;


  @omitnull()
  input: JSON.Raw | null = null;


  @omitnull()
  name: string | null = null;
}


@json
export class Usage {
  /**
   * The number of input tokens which were used.
   */
  @alias("input_tokens")
  inputTokens!: number;

  /**
   * The number of output tokens which were used.
   */
  @alias("output_tokens")
  outputTokens!: number;
}

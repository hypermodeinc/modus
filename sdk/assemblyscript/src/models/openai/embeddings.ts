/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Model } from "../../assembly/models"

/**
 * Provides input and output types that conform to the OpenAI Embeddings API.
 *
 * Reference: https://platform.openai.com/docs/api-reference/embeddings
 */
export class OpenAIEmbeddingsModel extends Model<OpenAIEmbeddingsInput, OpenAIEmbeddingsOutput> {
  /**
   * Creates an input object for the OpenAI Embeddings API.
   *
   * @param content The input content to vectorize.  Can be any of:
   * - A string representing the text to vectorize.
   * - An array of strings representing multiple texts to vectorize.
   * - An array of integers representing pre-tokenized text to vectorize.
   * - An array of arrays of integers representing multiple pre-tokenized texts to vectorize.
   *
   * @returns An input object that can be passed to the `invoke` method.
   *
   * @remarks
   * The input content must not exceed the maximum token limit of the model.
   */
  createInput<T>(content: T): OpenAIEmbeddingsInput {
    const model = this.info.fullName

    switch (idof<T>()) {
      case idof<string>():
      case idof<string[]>():
      case idof<i64[]>():
      case idof<i32[]>():
      case idof<i16[]>():
      case idof<i8[]>():
      case idof<u64[]>():
      case idof<u32[]>():
      case idof<u16[]>():
      case idof<u8[]>():
      case idof<i64[][]>():
      case idof<i32[][]>():
      case idof<i16[][]>():
      case idof<i8[][]>():
      case idof<u64[][]>():
      case idof<u32[][]>():
      case idof<u16[][]>():
      case idof<u8[][]>():
        return <TypedEmbeddingsInput<T>>{ model, input: content }
    }

    throw new Error("Unsupported input content type.")
  }
}

/**
 * The input object for the OpenAI Embeddings API.
 */
@json
export class OpenAIEmbeddingsInput {
  /**
   * The name of the model to use for the embeddings.
   * Must be the exact string expected by the model provider.
   * For example, "text-embedding-3-small".
   *
   * @remarks
   * This field is automatically set by the `createInput` method when creating this object.
   * It does not need to be set manually.
   */
  model!: string

  /**
   * The encoding format for the output embeddings.
   *
   * @default EncodingFormat.Float
   *
   * @remarks
   * Currently only `EncodingFormat.Float` is supported.
   */
  @alias("encoding_format")
  @omitif("this.encodingFormat == 'float'")
  encodingFormat: string = EncodingFormat.Float

  /**
   * The maximum number of dimensions for the output embeddings.
   * If not specified, the model's default number of dimensions will be used.
   */
  @omitif("this.dimensions == -1")
  dimensions: i32 = -1 // TODO: make this an `i32 | null` when supported

  /**
   * The user ID to associate with the request.
   * If not specified, the request will be anonymous.
   * See https://platform.openai.com/docs/guides/safety-best-practices/end-user-ids
   */
  @omitnull()
  user: string | null = null
}

/**
 * The input object for the OpenAI Embeddings API.
 */
@json
export class TypedEmbeddingsInput<T> extends OpenAIEmbeddingsInput {
  /**
   * The input content to vectorize.
   */
  input!: T
}

/**
 * The output object for the OpenAI Embeddings API.
 */
@json
export class OpenAIEmbeddingsOutput {
  /**
   * The name of the output object type returned by the API.
   * Always `"list"`.
   */
  object!: string

  /**
   * The name of the model used to generate the embeddings.
   * In most cases, this will match the requested `model` field in the input.
   */
  model!: string

  /**
   * The usage statistics for the request.
   */
  usage!: Usage

  /**
   * The output vector embeddings data.
   */
  data!: Embedding[]
}

/**
 * The encoding format for the output embeddings.
 */
// eslint-disable-next-line @typescript-eslint/no-namespace
export namespace EncodingFormat {
  /**
   * The output embeddings are encoded as an array of floating-point numbers.
   */
  export const Float = "float"

  /**
   * The output embeddings are encoded as a base64-encoded string,
   * containing an binary representation of an array of floating-point numbers.
   *
   * @remarks
   * This format is currently not supported through this interface.
   */
  export const Base64 = "base64"
}
export type EncodingFormat = string

/**
 * The output vector embeddings data.
 */
@json
export class Embedding {
  /**
   * The name of the output object type returned by the API.
   * Always `"embedding"`.
   */
  object!: string

  /**
   * The index of the input text that corresponds to this embedding.
   * Used when requesting embeddings for multiple texts.
   */
  index!: i32

  /**
   * The vector embedding of the input text.
   */
  embedding!: f32[] // TODO: support `f32[] | string` based on input encoding format
}

/**
 * The usage statistics for the request.
 */
@json
export class Usage {
  /**
   * The number of prompt tokens used in the request.
   */
  @alias("prompt_tokens")
  promptTokens!: i32

  /**
   * The total number of tokens used in the request.
   */
  @alias("total_tokens")
  totalTokens!: i32
}

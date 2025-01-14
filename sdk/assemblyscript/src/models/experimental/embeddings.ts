/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Model } from "../../assembly/models";

/**
 * A model that returns embeddings for a list of text strings.
 *
 * @remarks
 * This model interface is experimental and may change in the future.
 * It is primarily intended for use with with embedding models hosted on Hypermode.
 */
export class EmbeddingsModel extends Model<EmbeddingsInput, EmbeddingsOutput> {
  /**
   * Creates an input object for the embeddings model.
   *
   * @param instances - A list of one or more text strings to create vector embeddings for.
   * @returns An input object that can be passed to the `invoke` method.
   */
  createInput(instances: string[]): EmbeddingsInput {
    return <EmbeddingsInput>{ instances };
  }
}

/**
 * An input object for the embeddings model.
 */
@json
export class EmbeddingsInput {
  /**
   * A list of one or more text strings to create vector embeddings for.
   */
  instances!: string[];
}

/**
 * An output object for the embeddings model.
 */
@json
export class EmbeddingsOutput {
  /**
   * A list of vector embeddings that correspond to each input text string.
   */
  predictions!: f32[][];
}

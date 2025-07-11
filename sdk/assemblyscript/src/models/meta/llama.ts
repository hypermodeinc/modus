/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Model } from "../../assembly/models";

export class TextGenerationModel extends Model<
  TextGenerationInput,
  TextGenerationOutput
> {
  /**
   * Creates a new input object for the model.
   * @param prompt The prompt text to pass to the model.
   * @returns A new input object.
   * @remarks Optional properties may be set on the returned input object to
   * control the behavior of the model.
   */
  createInput(prompt: string): TextGenerationInput {
    return <TextGenerationInput>{ prompt };
  }
}


@json
export class TextGenerationInput {
  /**
   * The prompt text to pass to the model.
   * May contain special tokens to control the behavior of the model.
   * See https://llama.meta.com/docs/model-cards-and-prompt-formats/meta-llama-2/
   */
  prompt!: string;

  /**
   * The temperature of the generated text, which controls the randomness of the output.
   * Higher values may lead to mode creative but less coherent outputs,
   * while lower values may lead to more conservative but more coherent outputs.
   *
   * Default: 0.6
   */
  @omitif((self: TextGenerationInput) => self.temperature == 0.6)
  temperature: f64 = 0.6;

  /**
   * The maximum probability threshold for generating tokens.
   * Use a lower value to ignore less probable options.
   *
   * Default: 0.9
   */
  @omitif((self: TextGenerationInput) => self.topP == 0.9)
  @alias("top_p")
  topP: f64 = 0.9;

  /**
   * The maximum number of tokens in the generated text.
   *
   * Default: 512
   */
  @omitif((self: TextGenerationInput) => self.maxGenLen == 512)
  @alias("max_gen_len")
  maxGenLen: i32 = 512;
}


@json
export class TextGenerationOutput {
  /**
   * The generated text.
   */
  generation!: string;

  /**
   * The number of tokens in the prompt text.
   */
  @alias("prompt_token_count")
  promptTokenCount!: i32;

  /**
   * The number of tokens in the generated text.
   */
  @alias("generation_token_count")
  generationTokenCount!: i32;

  /**
   * The reason why the response stopped generating text.
   * Possible values are:
   * - `"stop"` - The model has finished generating text.
   * - `"length"` - The response has been truncated due to reaching the maximum
   *   token length specified by `maxGenLen` in the input.
   */
  @alias("stop_reason")
  stopReason!: string;
}

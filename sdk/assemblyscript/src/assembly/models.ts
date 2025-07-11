/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import * as utils from "./utils";

type ModelInvoker = (modelName: string, inputJson: string) => string | null;

// @ts-expect-error: decorator
@external("modus_models", "getModelInfo")
declare function hostGetModelInfo(modelName: string): ModelInfo;

// @ts-expect-error: decorator
@external("modus_models", "invokeModel")
declare function hostInvokeModel(
  modelName: string,
  input: string,
): string | null;

class ModusModelFactory implements ModelFactory {
  constructor() {
    // Note, we assign this to a static property on the base Model class so that it can be accessed
    // from the `invoke` method.  It would be preferable to use an instance private or protected
    // property on the model instance, and that does compile in AssemblyScript, but ends up displaying
    // an error in VS Code when the model type is requested from `getModel` in the factory.
    Model.invoker = hostInvokeModel;
  }

  /**
   * Gets a model object instance.
   * @param modelName The name of the model, as defined in the manifest.
   * @returns An instance of the model object, which can be used to interact with the model.
   */
  getModel<T extends Model>(modelName: string): T {
    const info = hostGetModelInfo(modelName);
    if (utils.resultIsInvalid(info)) {
      throw new Error(`Model ${modelName} not found.`);
    }

    return instantiate<T>(info);
  }
}

export class ModelInfo {
  constructor(
    public readonly name: string,
    public readonly fullName: string = name,
  ) {}
}

export interface ModelFactory {
  getModel<T extends Model>(modelName: string): T;
}


@json
export abstract class ModelError {
  abstract toString(): string;
}

export abstract class Model<TInput = unknown, TOutput = unknown> {
  static invoker: ModelInvoker | null = null;
  protected constructor(public info: ModelInfo) {}

  debug: boolean = false;
  validator: ((data: string) => ModelError | null) | null = null;

  /**
   * Invokes the model with the given input.
   * @param input The input object to pass to the model.
   * @returns The output object from the model.
   */
  invoke(input: TInput): TOutput {
    if (!Model.invoker) {
      throw new Error("Model invoker is not set.");
    }

    const modelName = this.info.name;
    const inputJson = JSON.stringify(input);
    if (this.debug) {
      console.debug(`Invoking ${modelName} model with input: ${inputJson}`);
    }

    const outputJson = Model.invoker(modelName, inputJson);
    if (outputJson === null) {
      throw new Error(`Failed to invoke ${modelName} model.`);
    }

    if (this.debug) {
      console.debug(`Received output: ${outputJson}`);
    }

    if (this.validator) {
      const err = this.validator(outputJson);
      if (err !== null) {
        utils.throwUserError(`The chat model returned an error: ${err}`);
      }
    }

    return JSON.parse<TOutput>(outputJson);
  }
}

const factory = new ModusModelFactory();
export default factory;

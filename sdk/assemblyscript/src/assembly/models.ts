/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as utils from "./utils";
import { Model, ModelFactory, ModelInfo } from "@hypermode/models-as";

// @ts-expect-error: decorator
@external("hypermode", "lookupModel")
declare function hostLookupModel(modelName: string): ModelInfo;

// @ts-expect-error: decorator
@external("hypermode", "invokeModel")
declare function hostInvokeModel(
  modelName: string,
  input: string,
): string | null;

class HypermodeModelFactory implements ModelFactory {
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
    const info = hostLookupModel(modelName);
    if (utils.resultIsInvalid(info)) {
      throw new Error(`Model ${modelName} not found.`);
    }

    return instantiate<T>(info);
  }
}

const factory = new HypermodeModelFactory();
export default factory;

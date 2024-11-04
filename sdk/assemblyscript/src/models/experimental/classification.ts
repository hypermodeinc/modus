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
 * A model that returns classification results for a list of text strings.
 *
 * @remarks
 * This model interface is experimental and may change in the future.
 * It is primarily intended for use with with classification models hosted on Hypermode.
 */
export class ClassificationModel extends Model<ClassificationInput, ClassificationOutput> {
  /**
   * Creates an input object for the classification model.
   *
   * @param instances - A list of one or more text strings to classify.
   * @returns An input object that can be passed to the `invoke` method.
   */
  createInput(instances: string[]): ClassificationInput {
    return <ClassificationInput>{ instances }
  }
}

/**
 * An input object for the classification model.
 */
@json
export class ClassificationInput {
  /**
   * A list of one or more text strings of text to classify.
   */
  instances!: string[]
}

/**
 * An output object for the classification model.
 */
@json
export class ClassificationOutput {
  /**
   * A list of prediction results that correspond to each input text string.
   */
  predictions!: ClassifierResult[]
}

/**
 * A classification result for a single text string.
 */
@json
export class ClassifierResult {
  /**
   * The classification label with the highest confidence.
   */
  label!: string

  /**
   * The confidence score for the classification label.
   */
  confidence!: f32

  /**
   * The list of all classification labels with their corresponding probabilities.
   */
  probabilities!: ClassifierLabel[]
}

/**
 * A classification label with its corresponding probability.
 */
@json
export class ClassifierLabel {
  /**
   * The classification label.
   */
  label!: string

  /**
   * The probability value.
   */
  probability!: f32
}

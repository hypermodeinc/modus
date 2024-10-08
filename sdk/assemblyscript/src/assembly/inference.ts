/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as utils from "./utils";
import { JSON } from "json-as";

// NOTE: all functions in this file are deprecated and will be removed in the future.

// @ts-expect-error: decorator
@external("hypermode", "invokeClassifier")
declare function invokeClassifier(
  modelName: string,
  sentenceMap: Map<string, string>,
): Map<string, Map<string, f32>>;

// @ts-expect-error: decorator
@external("hypermode", "computeEmbedding")
declare function computeEmbedding(
  modelName: string,
  sentenceMap: Map<string, string>,
): Map<string, f64[]>;

// @ts-expect-error: decorator
@external("hypermode", "invokeTextGenerator")
declare function invokeTextGenerator(
  modelName: string,
  instruction: string,
  sentence: string,
  format: string,
): string;

/**
 * @deprecated
 * This class is deprecated and will be removed in the future.
 * Use `models.getModel` and the appropriate model interface instead.
 */
export abstract class inference {
  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a classification model instead.
   */
  public static getClassificationProbability(
    modelName: string,
    text: string,
    label: string,
  ): f32 {
    const labels = this.getClassificationLabelsForText(modelName, text);
    const keys = labels.keys();
    for (let i = 0; i < keys.length; i++) {
      const key = keys[i];
      if (key === label) {
        return labels.get(key);
      }
    }
    return 0.0;
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a classification model instead.
   */
  public static classifyText(
    modelName: string,
    text: string,
    threshold: f32,
  ): string {
    const labels = this.getClassificationLabelsForText(modelName, text);

    const keys = labels.keys();
    let max = labels.get(keys[0]);
    let result = keys[0];
    for (let i = 1; i < keys.length; i++) {
      const key = keys[i];
      const value = labels.get(key);
      if (value >= max) {
        max = value;
        result = key;
      }
    }

    if (max < threshold) {
      return "";
    }
    return result;
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a classification model instead.
   */
  public static getClassificationLabelsForText(
    modelName: string,
    text: string,
  ): Map<string, f32> {
    const textMap = new Map<string, string>();
    textMap.set("text", text);
    const res = this.getClassificationLabelsForTexts(modelName, textMap);
    return res.get("text");
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a classification model instead.
   */
  public static getClassificationLabelsForTexts(
    modelName: string,
    texts: Map<string, string>,
  ): Map<string, Map<string, f32>> {
    const result = invokeClassifier(modelName, texts);
    if (utils.resultIsInvalid(result)) {
      throw new Error("Unable to classify text.");
    }
    return result;
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with an embedding model instead.
   */
  public static getTextEmbedding(modelName: string, text: string): f64[] {
    const textMap = new Map<string, string>();
    textMap.set("text", text);
    const res = this.getTextEmbeddings(modelName, textMap);
    return res.get("text");
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with an embedding model instead.
   */
  public static getTextEmbeddings(
    modelName: string,
    texts: Map<string, string>,
  ): Map<string, f64[]> {
    const result = computeEmbedding(modelName, texts);
    if (utils.resultIsInvalid(result)) {
      throw new Error("Unable to compute embeddings.");
    }
    return result;
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a generative text model instead.
   */
  public static generateText(
    modelName: string,
    instruction: string,
    prompt: string,
  ): string {
    const result = invokeTextGenerator(modelName, instruction, prompt, "text");
    if (utils.resultIsInvalid(result)) {
      throw new Error("Unable to generate text.");
    }
    return result;
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a generative text model instead.
   */
  public static generateData<TData>(
    modelName: string,
    instruction: string,
    text: string,
    sample: TData,
  ): TData {
    // Prompt trick: ask for a simple JSON object.
    const modifiedInstruction =
      "Only respond with valid JSON object in this format:\n" +
      JSON.stringify(sample) +
      "\n" +
      instruction;

    const result = invokeTextGenerator(
      modelName,
      modifiedInstruction,
      text,
      "json_object",
    );

    if (utils.resultIsInvalid(result)) {
      throw new Error("Unable to generate data.");
    }

    return JSON.parse<TData>(result);
  }

  /**
   * @deprecated
   * This function is deprecated and will be removed in the future.
   * Use `models.getModel` with a generative text model instead.
   */
  public static generateList<TData>(
    modelName: string,
    instruction: string,
    text: string,
    sample: TData,
  ): TData[] {
    // Prompt trick: ask for a simple JSON object containing a list.
    // Note, OpenAI will not generate an array of objects directly.
    const modifiedInstruction =
      "Only respond with valid JSON object containing a valid JSON array named 'list', in this format:\n" +
      '{"list":[' +
      JSON.stringify(sample) +
      "]}\n" +
      instruction;

    const result = invokeTextGenerator(
      modelName,
      modifiedInstruction,
      text,
      "json_object",
    );

    if (utils.resultIsInvalid(result)) {
      throw new Error("Unable to generate data.");
    }

    const jsonList = JSON.parse<Map<string, TData[]>>(result);
    return jsonList.get("list");
  }
}

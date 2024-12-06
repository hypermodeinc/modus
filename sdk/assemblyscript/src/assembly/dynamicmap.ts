/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";

export class DynamicMap {
  private data: Map<string, JSON.Raw> = new Map<string, JSON.Raw>();

  public get size(): i32 {
    return this.data.size;
  }

  public has(key: string): bool {
    return this.data.has(key);
  }

  public get<T>(key: string): T {
    return JSON.parse<T>(this.data.get(key));
  }

  public set<T>(key: string, value: T): void {
    if (value == null) {
      this.data.set(key, "null");
    } else {
      this.data.set(key, JSON.stringify(value));
    }
  }

  public delete(key: string): bool {
    return this.data.delete(key);
  }

  public clear(): void {
    this.data.clear();
  }

  public keys(): string[] {
    return this.data.keys();
  }

  public values(): JSON.Raw[] {
    return this.data.values();
  }

  __INITIALIZE(): this {
    return this;
  }

  __SERIALIZE(): string {
    // This would be ideal, but doesn't work:
    //   return JSON.stringify(this.data);
    // So instead, we have to do construct the JSON manually.
    //
    // TODO: Update this to use JSON.stringify once it's supported in json-as.
    // https://github.com/JairusSW/as-json/issues/98

    const segments: string[] = [];
    const keys = this.data.keys();
    const values = this.data.values();

    for (let i = 0; i < this.data.size; i++) {
      const key = JSON.stringify(keys[i]);
      const value = values[i]; // already in JSON
      segments.push(`${key}:${value}`);
    }

    return `{${segments.join(",")}}`;
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    data: string,
    key_start: i32,
    key_end: i32,
    value_start: i32,
    value_end: i32,
  ): boolean {
    // This would be ideal, but doesn't work:
    //   this.data = JSON.parse<Map<string, JSON.Raw>>(data);
    // So instead, we have to parse the JSON manually.
    //
    // TODO: Update this to use JSON.parse once it's supported in json-as.
    // https://github.com/JairusSW/as-json/issues/98

    this.data = parseJsonToRawMap(data);

    return true;
  }
}

// Everything below this line is a workaround for JSON.parse not working with JSON.Raw values.
// (It was generated via AI and may not be optimized.)

class ParseResult {
  constructor(
    public result: string,
    public nextIndex: i32,
  ) {}
}

function parseJsonToRawMap(jsonString: string): Map<string, string> {
  const map = new Map<string, string>();
  let i: i32 = 0;
  const length: i32 = jsonString.length;

  function skipWhitespace(index: i32, input: string): i32 {
    while (
      index < input.length &&
      (input.charCodeAt(index) === 32 || // space
        input.charCodeAt(index) === 9 || // tab
        input.charCodeAt(index) === 10 || // newline
        input.charCodeAt(index) === 13) // carriage return
    ) {
      index++;
    }
    return index;
  }

  function readString(index: i32, input: string): ParseResult {
    if (input.charCodeAt(index) !== 34) {
      throw new Error("Expected '\"' at position " + index.toString());
    }
    index++;
    let result = "";
    while (index < input.length) {
      const charCode = input.charCodeAt(index);
      if (charCode === 34) {
        return new ParseResult(result, index + 1);
      }
      if (charCode === 92) {
        index++;
        if (index < input.length) {
          result += String.fromCharCode(input.charCodeAt(index));
        }
      } else {
        result += String.fromCharCode(charCode);
      }
      index++;
    }
    throw new Error("Unterminated string");
  }

  function readValue(index: i32, input: string): ParseResult {
    const start: i32 = index;
    let braceCount: i32 = 0;
    let bracketCount: i32 = 0;

    while (index < input.length) {
      const char = input.charAt(index);
      if (
        braceCount === 0 &&
        bracketCount === 0 &&
        (char === "," || char === "}")
      ) {
        break;
      }
      if (char === "{") braceCount++;
      else if (char === "}") braceCount--;
      else if (char === "[") bracketCount++;
      else if (char === "]") bracketCount--;
      index++;
    }

    const result = input.substring(start, index).trim();
    return new ParseResult(result, index);
  }

  if (jsonString.charCodeAt(i) !== 123) {
    throw new Error("Expected '{' at the beginning of JSON");
  }
  i++;

  while (i < length) {
    i = skipWhitespace(i, jsonString);

    const keyResult = readString(i, jsonString);
    const key = keyResult.result;
    i = keyResult.nextIndex;

    i = skipWhitespace(i, jsonString);
    if (jsonString.charCodeAt(i) !== 58) {
      throw new Error("Expected ':' after key at position " + i.toString());
    }
    i++;

    i = skipWhitespace(i, jsonString);

    const valueResult = readValue(i, jsonString);
    const value = valueResult.result;
    i = valueResult.nextIndex;

    map.set(key, value);

    i = skipWhitespace(i, jsonString);
    if (jsonString.charCodeAt(i) === 44) {
      i++;
    } else if (jsonString.charCodeAt(i) === 125) {
      i++;
      break;
    } else {
      throw new Error(
        "Unexpected character '" +
          jsonString.charAt(i) +
          "' at position " +
          i.toString(),
      );
    }
  }

  return map;
}

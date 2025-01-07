/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import {
  BACK_SLASH,
  BRACE_LEFT,
  BRACE_RIGHT,
  BRACKET_LEFT,
  BRACKET_RIGHT,
  CHAR_F,
  CHAR_N,
  CHAR_T,
  COLON,
  COMMA,
  QUOTE,
} from "json-as/assembly/custom/chars";
import { unsafeCharCodeAt } from "json-as/assembly/custom/util";
import { isSpace } from "util/string";

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

    this.data = deserializeRawMap(data);

    return true;
  }
}

function deserializeRawMap(
  src: string,
  dst: Map<string, string> = new Map<string, string>(),
): Map<string, string> {
  const srcPtr = changetype<usize>(src);
  let depth = 0;
  let isKey = false;
  let index = 0;
  let lastIndex = 0;
  let end = src.length - 1;
  let key: string | null = null;

  while (index < end && isSpace(unsafeCharCodeAt(src, index))) index++;
  while (end > index && isSpace(unsafeCharCodeAt(src, end))) end--;

  if (end - index <= 0)
    throw new Error(
      "Input string had zero length or was all whitespace at position " +
        index.toString(),
    );

  if (unsafeCharCodeAt(src, index++) !== BRACE_LEFT)
    throw new Error("Expected '{' at position " + index.toString());

  while (index < end) {
    while (isSpace(unsafeCharCodeAt(src, index)) && index < end) index++;
    const code = unsafeCharCodeAt(src, index);
    if (key == null) {
      if (code == QUOTE && unsafeCharCodeAt(src, index - 1) !== BACK_SLASH) {
        if (isKey) {
          if (lastIndex == index)
            throw new Error(
              "Found empty key '\"\"' at position " + index.toString(),
            );
          key = src.slice(lastIndex, index);
          // console.log("Key: " + key);
          while (isSpace(unsafeCharCodeAt(src, ++index))) {
            /* empty */
          }
          if (unsafeCharCodeAt(src, index) !== COLON)
            throw new Error(
              "Expected ':' after key at position " + index.toString(),
            );
          isKey = false;
        } else {
          isKey = true;
        }
        lastIndex = index + 1;
      }
    } else {
      if (code == BRACE_LEFT) {
        depth++;
        index++;
        while (index < end) {
          const code = unsafeCharCodeAt(src, index);
          if (code == BRACE_RIGHT) {
            if (--depth == 0) {
              while (isSpace(unsafeCharCodeAt(src, ++index)) && index < end) {
                /* empty */
              }
              const value = src.slice(lastIndex, index);
              dst.set(key, value);
              // console.log("Value (object): " + value);
              const last = index == end;
              if (!last && unsafeCharCodeAt(src, index) !== COMMA)
                throw new Error(
                  "Expected ',' after value at position " + index.toString(),
                );
              if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
                throw new Error(
                  "Expected '}' after value at position " + index.toString(),
                );
              key = null;
              break;
            }
          } else if (code == BRACE_LEFT) depth++;
          index++;
        }
      } else if (code == BRACKET_LEFT) {
        depth++;
        index++;
        while (index < end) {
          const code = unsafeCharCodeAt(src, index);
          if (code == BRACKET_RIGHT) {
            if (--depth == 0) {
              while (isSpace(unsafeCharCodeAt(src, ++index)) && index < end) {
                /* empty */
              }
              const value = src.slice(lastIndex, index);
              dst.set(key, value);
              // console.log("Value (object): " + value);
              const last = index == end;
              if (!last && unsafeCharCodeAt(src, index) !== COMMA)
                throw new Error(
                  "Expected ',' after value at position " + index.toString(),
                );
              if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
                throw new Error(
                  "Expected '}' after value at position " + index.toString(),
                );
              key = null;
              break;
            }
          } else if (code == BRACKET_LEFT) depth++;
          index++;
        }
      } else if (code == CHAR_T) {
        if (load<u64>(srcPtr + (index << 1)) == 28429475166421108) {
          const value = src.slice(lastIndex, (index += 4));
          dst.set(key, value);
          // console.log("Value (bool): " + value);
          while (isSpace(unsafeCharCodeAt(src, index + 1)) && index < end)
            index++;
          const last = index == end;
          if (!last && unsafeCharCodeAt(src, index) !== COMMA)
            throw new Error(
              "Expected ',' after value at position " + index.toString(),
            );
          if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
            throw new Error(
              "Expected '}' after value at position " + index.toString(),
            );
          key = null;
        } else {
          throw new Error(
            "Expected 'true' as value but found '" +
              key +
              "' instead at position " +
              index.toString(),
          );
        }
      } else if (code == CHAR_F) {
        // eslint-disable-next-line no-loss-of-precision
        if (load<u64>(srcPtr + (index << 1), 2) == 28429466576093281) {
          const value = src.slice(lastIndex, (index += 5));
          dst.set(key, value);
          // console.log("Value (bool): " + value);
          while (isSpace(unsafeCharCodeAt(src, index + 1)) && index < end)
            index++;
          const last = index == end;
          if (!last && unsafeCharCodeAt(src, index) !== COMMA)
            throw new Error(
              "Expected ',' after value at position " + index.toString(),
            );
          if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
            throw new Error(
              "Expected '}' after value at position " + index.toString(),
            );
          key = null;
        } else {
          throw new Error(
            "Expected 'false' as value but found '" +
              key +
              "' instead at position " +
              index.toString(),
          );
        }
      } else if (code == CHAR_N) {
        // eslint-disable-next-line no-loss-of-precision
        if (load<u64>(srcPtr + (index << 1)) == 30399761348886638) {
          const value = src.slice(lastIndex, (index += 4));
          dst.set(key, value);
          // console.log("Value (null): " + value);
          while (isSpace(unsafeCharCodeAt(src, index)) && index < end) index++;
          const last = index == end;
          if (!last && unsafeCharCodeAt(src, index) !== COMMA)
            throw new Error(
              "Expected ',' after value at position " + index.toString(),
            );
          if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
            throw new Error(
              "Expected '}' after value at position " + index.toString(),
            );
          key = null;
        } else {
          throw new Error(
            "Expected 'null' as value but found '" +
              key +
              "' instead at position " +
              index.toString(),
          );
        }
      } else if (code == QUOTE) {
        index++;
        while (index < end) {
          const code = unsafeCharCodeAt(src, index);
          if (
            code == QUOTE &&
            unsafeCharCodeAt(src, index - 1) !== BACK_SLASH
          ) {
            while (isSpace(unsafeCharCodeAt(src, ++index)) && index < end) {
              /* empty */
            }
            const value = src.slice(lastIndex, index);
            dst.set(key, value);
            // console.log("Value (string): " + value);
            const last = index == end;
            if (!last && unsafeCharCodeAt(src, index) !== COMMA)
              throw new Error(
                "Expected ',' after value at position " + index.toString(),
              );
            if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
              throw new Error(
                "Expected '}' after value at position " + index.toString(),
              );
            key = null;
            break;
          }
          index++;
        }
      } else if ((code >= 48 && code <= 57) || code == 45) {
        lastIndex = index++;
        while (index < end) {
          const code = unsafeCharCodeAt(src, index);
          if (code == COMMA || code == BRACE_RIGHT || isSpace(code)) {
            const value = src.slice(lastIndex, index);
            dst.set(key, value);
            // console.log("Value (number): " + value);
            while (isSpace(unsafeCharCodeAt(src, index + 1)) && index < end)
              index++;
            const last = index == end;
            if (!last && unsafeCharCodeAt(src, index) !== COMMA)
              throw new Error(
                "Expected ',' after value at position " + index.toString(),
              );
            if (last && unsafeCharCodeAt(src, index) !== BRACE_RIGHT)
              throw new Error(
                "Expected '}' after value at position " + index.toString(),
              );
            key = null;
            break;
          }
          index++;
        }
      } else {
        throw new Error(
          "Expected valid character after key but found'" +
            String.fromCharCode(code) +
            "' at position " +
            index.toString(),
        );
      }
    }
    index++;
  }

  if (isKey)
    throw new Error("Unterminated key at position " + lastIndex.toString());

  if (unsafeCharCodeAt(src, end) !== BRACE_RIGHT)
    throw new Error("Expected '{' at position " + end.toString());

  return dst;
}

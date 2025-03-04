/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import { ptrToStr } from "json-as/assembly/util/ptrToStr";

@json
export class DynamicMap {
  private data: Map<string, JSON.Raw> = new Map<string, JSON.Raw>();

  public get size(): i32 {
    return this.data.size;
  }

  public has(key: string): bool {
    return this.data.has(key);
  }

  public get<T>(key: string): T {
    return JSON.parse<T>(this.data.get(key).data);
  }

  public set<T>(key: string, value: T): void {
    if (isInteger<T>() && nameof<T>() == "usize" && value == 0) {
      this.data.set(key, JSON.Raw.from("null"));
    } else {
      this.data.set(key, JSON.Raw.from(JSON.stringify(value)));
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

  @serializer
  serialize(self: DynamicMap): string {
    // This would be ideal, but doesn't work:
    //   return JSON.stringify(this.data);
    // So instead, we have to do construct the JSON manually.
    //
    // TODO: Update this to use JSON.stringify once it's supported in json-as.
    // https://github.com/JairusSW/as-json/issues/98

    let out = "{";

    const keys = self.data.keys();
    const values = self.data.values();

    const end = self.data.size - 1;

    for (let i = 0; i < end; i++) {
      const key = unchecked(keys[i]);
      const value = unchecked(values[i]);
      out += JSON.stringify(key) + ":" + value.data + ",";
    }
    const key = unchecked(keys[end]);
    const value = unchecked(values[end]);
    out += JSON.stringify(key) + ":" + value.data;

    return out + "}";
  }

  /* eslint-disable @typescript-eslint/no-unused-vars */
  __DESERIALIZE(
    keyStart: usize,
    keyEnd: usize,
    valStart: usize,
    valEnd: usize,
    out: usize
  ): void {
    // This would be ideal, but doesn't work:
    //   this.data = JSON.parse<Map<string, JSON.Raw>>(data);
    // So instead, we have to parse the JSON manually.
    //
    // TODO: Update this to use JSON.parse once it's supported in json-as.
    // https://github.com/JairusSW/as-json/issues/98
    const key = ptrToStr(keyStart, keyEnd);
    const value = ptrToStr(valStart, valEnd);
    this.data.set(key, JSON.Raw.from(value));
  }
}
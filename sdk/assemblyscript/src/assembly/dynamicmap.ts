/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";


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
    this.data.set(key, JSON.Raw.from(JSON.stringify<T>(value)));
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
  private serialize(self: DynamicMap): string {
    return JSON.stringify(self.data);
  }


  @deserializer
  private deserialize(data: string): DynamicMap {
    const dm = new DynamicMap();
    dm.data = JSON.parse<Map<string, JSON.Raw>>(data);
    return dm;
  }
}

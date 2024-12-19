/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import { expect, it, run } from "as-test";
import { DynamicMap } from "../dynamicmap";


@json
class Obj {
  foo: string = "";
}

it("should set values", () => {
  const m = new DynamicMap();
  m.set("a", 42);
  m.set("b", "hello");
  m.set("c", [1, 2, 3]);
  m.set("d", true);
  m.set("e", null);
  m.set("f", 3.14);
  m.set("g", { foo: "bar" } as Obj);

  const json = JSON.stringify(m);
  expect(json).toBe(
    '{"a":42,"b":"hello","c":[1,2,3],"d":true,"e":null,"f":3.14,"g":{"foo":"bar"}}',
  );
});

it("should get values", () => {
  const m = JSON.parse<DynamicMap>(
    '{"a":42,"b":"hello","c":[1,2,3],"d":true,"e":null,"f":3.14,"g":{"foo":"bar"}}',
  );
  expect(m.get<i32>("a")).toBe(42);
  expect(m.get<string>("b")).toBe("hello");
  expect(m.get<i32[]>("c")).toBe([1, 2, 3]);
  expect(m.get<bool>("d")).toBe(true);
  expect(m.get<Obj | null>("e")).toBe(null);
  expect(m.get<f64>("f")).toBe(3.14);

  const obj = m.get<Obj>("g");
  expect(obj.foo).toBe("bar");
});

it("should get size", () => {
  const m = new DynamicMap();
  expect(m.size).toBe(0);
  m.set("a", 42);
  expect(m.size).toBe(1);
  m.set("b", "hello");
  expect(m.size).toBe(2);
});

it("should test existence of keys", () => {
  const m = new DynamicMap();
  expect(m.has("a")).toBe(false);
  m.set("a", 42);
  expect(m.has("a")).toBe(true);
  expect(m.has("b")).toBe(false);
  m.set("b", "hello");
  expect(m.has("b")).toBe(true);
});

it("should delete keys", () => {
  const m = new DynamicMap();
  m.set("a", 42);
  m.set("b", "hello");
  expect(m.size).toBe(2);
  expect(m.has("a")).toBe(true);
  expect(m.has("b")).toBe(true);
  m.delete("a");
  expect(m.size).toBe(1);
  expect(m.has("a")).toBe(false);
  expect(m.has("b")).toBe(true);
  m.delete("b");
  expect(m.size).toBe(0);
  expect(m.has("a")).toBe(false);
  expect(m.has("b")).toBe(false);
});

it("should clear", () => {
  const m = new DynamicMap();
  m.set("a", 42);
  m.set("b", "hello");
  expect(m.size).toBe(2);
  m.clear();
  expect(m.size).toBe(0);
  expect(m.has("a")).toBe(false);
  expect(m.has("b")).toBe(false);
});

it("should iterate keys", () => {
  const m = new DynamicMap();
  m.set("a", 42);
  m.set("b", "hello");
  m.set("c", [1, 2, 3]);
  const keys = m.keys();
  expect(keys).toBe(["a", "b", "c"]);
});

it("should iterate raw values", () => {
  const m = new DynamicMap();
  m.set("a", 42);
  m.set("b", "hello");
  m.set("c", [1, 2, 3]);
  const values = m.values();
  expect(values).toBe(["42", '"hello"', "[1,2,3]"]);
});

run();

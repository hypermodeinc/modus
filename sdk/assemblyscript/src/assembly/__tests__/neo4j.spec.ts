/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { JSON } from "json-as";
import { expect, it, run } from "as-test";
import { neo4j } from "..";

it("should stringify a simple record", () => {
  const record = new neo4j.Record();
  record.Keys = ["name"];
  record.Values = ['"Alice"'];

  expect(JSON.stringify(record)).toBe('{"name":"Alice"}');
});

it("should stringify different types of record values", () => {
  const record = new neo4j.Record();
  record.Keys = ["name", "age", "friends"];
  record.Values = ['"Alice"', '"42"', '["Bob","Peter","Anna"]'];

  expect(JSON.stringify(record)).toBe(
    '{"name":"Alice","age":"42","friends":["Bob","Peter","Anna"]}',
  );
});

run();

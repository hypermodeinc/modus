/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import test from "node:test";
import * as assert from "node:assert";
import { Docs, FunctionSignature } from "../lib/types.js";
import {
  CommentKind,
  Parser,
  Range,
} from "assemblyscript/dist/assemblyscript.js";
import ModusTransform from "../lib/index.js";
import { Extractor } from "../lib/extractor.js";

test("Docs.from creates Docs from comment nodes", () => {
  const nodes = [{ commentKind: 2, text: "/**\n * This is a test\n */" }];
  const docs = Docs.from(nodes);
  assert.ok(docs, "Docs should be created");
  assert.deepStrictEqual(docs?.lines, ["This is a test"]);
});

test("Docs.from only creates Docs from comment nodes", () => {
  const nodes = [
    { commentKind: CommentKind.Triple, text: "/// This is a triple comment" },
    {
      commentKind: CommentKind.Triple,
      text: "// This is a single line comment",
    },
    { commentKind: CommentKind.Block, text: "/**\n * This is a test\n */" },
  ];
  const docs = Docs.from(nodes);
  assert.ok(docs, "Docs should be created");
  assert.deepStrictEqual(docs?.lines.length, 1);
  assert.deepStrictEqual(docs?.lines, ["This is a test"]);
});

test("Comments at the start of a file are parsed correctly", () => {
  const parser = new Parser(null, []);
  parser.parseFile(
    `
    /**
     * This is a JSDoc comment describing a foo()
    **/
    export function foo(): void {}
    `,
    "./test.ts",
    true,
  );
  const extractor = new Extractor(new ModusTransform());

  const node = parser.currentSource.statements[0];
  const nodeIndex = parser.currentSource.statements.indexOf(node);
  const prevNode = parser.currentSource.statements[Math.max(nodeIndex - 1, 0)];

  const range = new Range(
    nodeIndex > 0 ? prevNode.range.end : 0,
    node.range.start,
  );
  range.source = parser.currentSource;
  const parsed = extractor["parseComments"](range);

  assert.equal(parsed.length, 1);
  assert.equal(
    parsed[0].text,
    "/**\n     * This is a JSDoc comment describing a foo()\n    **/",
  );
});

test("Comments in the middle of a file are parsed correctly", () => {
  const parser = new Parser(null, []);
  parser.parseFile(
    `
    export function bar(): void {}
    /**
     * This is a JSDoc comment describing a foo()
    **/
    export function foo(): void {}
    `,
    "./test.ts",
    true,
  );
  const extractor = new Extractor(new ModusTransform());

  const node = parser.currentSource.statements[1];
  const nodeIndex = parser.currentSource.statements.indexOf(node);
  const prevNode = parser.currentSource.statements[Math.max(nodeIndex - 1, 0)];

  const range = new Range(
    nodeIndex > 0 ? prevNode.range.end : 0,
    node.range.start,
  );
  range.source = parser.currentSource;
  const parsed = extractor["parseComments"](range);

  assert.equal(parsed.length, 1);
  assert.equal(
    parsed[0].text,
    "/**\n     * This is a JSDoc comment describing a foo()\n    **/",
  );
});

test("FunctionSignature.toString outputs function signature", () => {
  const signature = new FunctionSignature(
    "myFunction",
    [{ name: "param1", type: "i32" }],
    [{ name: "result", type: "i32" }],
  );
  assert.strictEqual(
    signature.toString(),
    "myFunction(param1: i32): i32",
    "Function signature should match",
  );
});

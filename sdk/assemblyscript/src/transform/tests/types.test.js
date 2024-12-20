import test from "node:test";
import * as assert from "node:assert";
import { FunctionSignature, TypeDefinition } from "../lib/types.js";

test("FunctionSignature.toJSON ignores empty or default fields", () => {
  const signature = new FunctionSignature(
    "myFunction",
    [{ name: "param1", type: "i32" }],
    [{ name: "result", type: "void" }],
  );
  const json = signature.toJSON();
  assert.deepStrictEqual(
    json,
    { parameters: [{ name: "param1", type: "i32" }] },
    "toJSON should omit default fields",
  );
});

test("TypeDefinition.toString formats type definitions correctly", () => {
  const typeDef = new TypeDefinition("MyType", 1, [
    { name: "field1", type: "i32" },
  ]);
  assert.strictEqual(
    typeDef.toString(),
    "MyType { field1: i32 }",
    "Type definition should match",
  );
});

test("TypeDefinition.isHidden returns true for internal types", () => {
  const hiddenType = new TypeDefinition("~lib/internalType", 1);
  assert.strictEqual(
    hiddenType.isHidden(),
    true,
    "Hidden type should return true",
  );

  const visibleType = new TypeDefinition("CustomType", 1);
  assert.strictEqual(
    visibleType.isHidden(),
    false,
    "Visible type should return false",
  );
});

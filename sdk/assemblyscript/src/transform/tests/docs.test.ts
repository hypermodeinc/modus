import test from "node:test"
import * as assert from "node:assert";
import { Docs, Field, FunctionSignature, TypeDefinition } from "../src/types.js";
import { CommentKind, FunctionDeclaration, Parser, Range } from "assemblyscript/dist/assemblyscript.js";
import { CommentNode } from "types:assemblyscript/src/ast";
import ModusTransform from "../src/index.js";
import { Extractor } from "../src/extractor.js";

test("Docs.from creates Docs from comment nodes", () => {
    const nodes = [
        { commentKind: 2, text: "/**\n * This is a test\n */" },
    ] as CommentNode[];
    const docs = Docs.from(nodes);
    assert.ok(docs, "Docs should be created");
    assert.deepStrictEqual(docs?.lines, ["This is a test"]);
});

test("Docs.from only creates Docs from comment nodes", () => {
    const nodes = [
        { commentKind: CommentKind.Triple, text: "/// This is a triple comment" },
        { commentKind: CommentKind.Triple, text: "// This is a single line comment" },
        { commentKind: CommentKind.Block, text: "/**\n * This is a test\n */" },
    ] as CommentNode[];
    const docs = Docs.from(nodes);
    assert.ok(docs, "Docs should be created");
    assert.deepStrictEqual(docs?.lines.length, 1);
    assert.deepStrictEqual(docs?.lines, ["This is a test"]);
});

test("Comments at the start of a file are parsed correctly", () => {
    const parser = new Parser(null, []);
    parser.parseFile(`
    /**
     * This is a JSDoc comment describing a foo()
    **/
    export function foo(): void {}
    `, "./test.ts", true);
    const extractor = new Extractor(new ModusTransform());

    const node = parser.currentSource.statements[0] as FunctionDeclaration;
    const nodeIndex = parser.currentSource.statements.indexOf(node);
    const prevNode = parser.currentSource.statements[Math.max(nodeIndex - 1, 0)];

    const range = new Range(nodeIndex > 0 ? prevNode.range.end : 0, node.range.start)
    range.source = parser.currentSource;
    const parsed = (extractor["parseComments"] as (range: Range) => CommentNode[])(range);

    assert.equal(parsed.length, 1)
    assert.equal(parsed[0].text, "/**\n     * This is a JSDoc comment describing a foo()\n    **/")
});

test("Comments in the middle of a file are parsed correctly", () => {
    const parser = new Parser(null, []);
    parser.parseFile(`
    export function bar(): void {}
    /**
     * This is a JSDoc comment describing a foo()
    **/
    export function foo(): void {}
    `, "./test.ts", true);
    const extractor = new Extractor(new ModusTransform());

    const node = parser.currentSource.statements[1] as FunctionDeclaration;
    const nodeIndex = parser.currentSource.statements.indexOf(node);
    const prevNode = parser.currentSource.statements[Math.max(nodeIndex - 1, 0)];

    const range = new Range(nodeIndex > 0 ? prevNode.range.end : 0, node.range.start)
    range.source = parser.currentSource;
    const parsed = (extractor["parseComments"] as (range: Range) => CommentNode[])(range);

    assert.equal(parsed.length, 1)
    assert.equal(parsed[0].text, "/**\n     * This is a JSDoc comment describing a foo()\n    **/")
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
        { name: "field1", type: "i32" } as Field,
    ]);
    assert.strictEqual(
        typeDef.toString(),
        "MyType { field1: i32 }",
        "Type definition should match",
    );
});

test("TypeDefinition.isHidden returns true for internal types", () => {
    const hiddenType = new TypeDefinition("~lib/internalType", 1);
    assert.strictEqual(hiddenType.isHidden(), true, "Hidden type should return true");

    const visibleType = new TypeDefinition("CustomType", 1);
    assert.strictEqual(visibleType.isHidden(), false, "Visible type should return false");
});

import test from "node:test"
import * as assert from "node:assert";
import { Metadata } from "../src/metadata.js";
import { existsSync, rmSync, writeFileSync } from "node:fs";
import * as path from "node:path";
import { fileURLToPath } from "node:url";

process.env.npm_package_version = process.env.npm_package_version || "v0.0.0";
process.env.npm_package_name = process.env.npm_package_name || "test";

test("Metadata.generate creates a new Metadata instance with required fields", () => {
    const filePath = path.join(
        path.dirname(fileURLToPath(import.meta.url)),
        "..",
        "..",
        "package.json",
    );
    const exists = existsSync(filePath)
    if (!exists) writeFileSync(filePath, JSON.stringify({ name: process.env.npm_package_name, version: process.env.npm_package_version }));

    const metadata = Metadata.generate();

    assert.ok(metadata.buildId, "buildId should be defined");
    assert.ok(metadata.buildTs, "buildTs should be defined");
    assert.ok(metadata.plugin, "plugin should be defined");
    assert.ok(metadata.sdk, "sdk should be defined");

    if (!exists) rmSync(filePath);
});

test("Metadata.addExportFn adds exported functions", () => {
    const metadata = new Metadata();
    const functions = [{ name: "foo", parameters: [], results: [] }];
    metadata.addExportFn(functions);
    assert.deepStrictEqual(metadata.fnExports["foo"], functions[0]);
});

// test("Metadata.addImportFn adds imported functions", () => {
//     const metadata = new Metadata();
//     const functions = [{ name: "bar", parameters: [], results: [] }];
//     metadata.addImportFn(functions);
//     assert.deepStrictEqual(metadata.fnImports["bar"], functions[0]);
// });

// test("Metadata.addTypes adds type definitions", () => {
//     const metadata = new Metadata();
//     const types = [{ name: "MyType", id: 1 }];
//     metadata.addTypes(types);
//     assert.deepStrictEqual(metadata.types["MyType"], types[0]);
// });

// test("Metadata.writeToModule adds custom sections to the WebAssembly module", () => {
//     const metadata = Metadata.generate();
//     const module = new binaryen.Module();

//     const addCustomSectionSpy = module.addCustomSection.bind(module);
//     module.addCustomSection = function (name, data) {
//         addCustomSectionSpy(name, data);
//         if (name === "modus_metadata_version") {
//             assert.deepStrictEqual(data, Uint8Array.from([2]));
//         }
//         if (name === "modus_metadata") {
//             assert.ok(data instanceof Uint8Array, "Metadata section should be Uint8Array");
//         }
//     };

//     metadata.writeToModule(module);
// });

// test("Metadata.logResults outputs metadata information to the console", () => {
//     const metadata = Metadata.generate();
//     const hook = console.log.bind(console);
//     let i = 0;

//     console.log = function (...args) {
//         i++;
//         hook(...args);
//     };

//     metadata.logResults();
//     assert.ok(i > 0, "logResults should produce console output");
//     console.log = hook;
// });
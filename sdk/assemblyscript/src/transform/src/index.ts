/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import assemblyscript from "assemblyscript/dist/assemblyscript.js";
import { Transform } from "assemblyscript/dist/transform.js";
import { Metadata } from "./metadata.js";
import { Extractor } from "./extractor.js";
import binaryen from "assemblyscript/lib/binaryen.js";
import { Parser } from "types:assemblyscript/src/parser";
import { Program } from "types:assemblyscript/src/program";

export default class ModusTransform extends Transform {
  private extractor = new Extractor(this);

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  afterParse(parser: Parser): void | Promise<void> {
    // This forces the compiler to include function exports from the Modus SDK without any user code.
    const stmt = `export * from "@hypermode/modus-sdk-as/assembly/exports";`;
    assemblyscript.parse(this.program, stmt, "", true);
  }

  afterInitialize(program: Program): void | Promise<void> {
    this.extractor.initHook(program);
  }

  afterCompile(module: binaryen.Module) {
    this.extractor.compileHook(module);

    const info = this.extractor.getProgramInfo();

    const m = Metadata.generate();
    m.addExportFn(info.exportFns);
    m.addImportFn(info.importFns);
    m.addTypes(info.types);
    m.writeToModule(module);
    m.logResults();
  }
}

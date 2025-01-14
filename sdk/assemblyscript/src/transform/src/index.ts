/*
 * Copyright 2025 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2025 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Transform } from "assemblyscript/dist/transform.js";
import { Metadata } from "./metadata.js";
import { Extractor } from "./extractor.js";
import binaryen from "assemblyscript/lib/binaryen.js";
import { Program } from "types:assemblyscript/src/program";

export default class ModusTransform extends Transform {
  private extractor = new Extractor(this);
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

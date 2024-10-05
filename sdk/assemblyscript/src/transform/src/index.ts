/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Transform } from "assemblyscript/dist/transform.js";
import { createWriteStream } from "fs";
import { Metadata } from "./metadata.js";
import { Extractor } from "./extractor.js";
import binaryen from "assemblyscript/lib/binaryen.js";

export default class ModusTransform extends Transform {
  afterCompile(module: binaryen.Module) {
    const extractor = new Extractor(this, module);
    const info = extractor.getProgramInfo();

    const m = Metadata.generate();
    m.addExportFn(info.exportFns);
    m.addImportFn(info.importFns);
    m.addTypes(info.types);
    m.writeToModule(module);

    // Write to stdout
    m.logToStream(process.stdout);

    // If running in GitHub Actions, also write to the step summary
    if (process.env.GITHUB_ACTIONS && process.env.GITHUB_STEP_SUMMARY) {
      const stream = createWriteStream(process.env.GITHUB_STEP_SUMMARY, {
        flags: "a",
      });
      m.logToStream(stream, true);
    }
  }
}

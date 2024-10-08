/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

/* ESLint configuration for AssemblyScript */

import * as ts from "typescript";
import * as parser from "@typescript-eslint/parser";
import utils from "../node_modules/@typescript-eslint/typescript-estree/dist/node-utils.js";

// In AssemblyScript, functions and variables can be decorated
const nodeCanBeDecorated = utils.nodeCanBeDecorated;
utils.nodeCanBeDecorated = function (node) {
  switch (node.kind) {
    case ts.SyntaxKind.FunctionDeclaration:
    case ts.SyntaxKind.VariableStatement:
      return true;
    default:
      return nodeCanBeDecorated(node);
  }
};

const config = {
  files: ["assembly/**/*.ts"],
  languageOptions: { parser: parser },
};

export default { config };

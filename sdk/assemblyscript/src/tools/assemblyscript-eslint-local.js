/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

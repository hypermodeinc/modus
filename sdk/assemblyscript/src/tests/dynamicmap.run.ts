/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { readFileSync } from "fs";
import { instantiate } from "../build/dynamicmap.spec.js";
const binary = readFileSync("./build/dynamicmap.spec.wasm");
const module = new WebAssembly.Module(binary);
instantiate(module, {
  env: {},
});

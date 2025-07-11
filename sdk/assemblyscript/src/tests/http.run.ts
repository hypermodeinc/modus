/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { readFileSync } from "fs";
import { instantiate } from "../build/http.spec.js";
const binary = readFileSync("./build/http.spec.wasm");
const module = new WebAssembly.Module(binary);
instantiate(module, {
  env: {},
  modus_models: {},
});

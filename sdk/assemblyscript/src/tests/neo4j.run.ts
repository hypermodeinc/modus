/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { readFileSync } from "fs";
import { instantiate } from "../build/neo4j.spec.js";
const binary = readFileSync("./build/neo4j.spec.wasm");
const module = new WebAssembly.Module(binary);
instantiate(module, {
  env: {},
  modus_models: {},
});
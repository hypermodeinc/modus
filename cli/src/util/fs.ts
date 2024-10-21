/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import fs from "node:fs";
export * from "node:fs/promises";

export async function exists(path: string) {
  try {
    await fs.promises.stat(path);
    return true;
  } catch {
    return false;
  }
}

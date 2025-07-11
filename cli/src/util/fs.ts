/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
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

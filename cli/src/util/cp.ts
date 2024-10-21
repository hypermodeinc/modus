/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import util from "node:util";
import cp from "node:child_process";

/**
 * Promisified version of `child_process.execFile`.
 */
const execFile = util.promisify(cp.execFile);
export { execFile };

/**
 * Promisified version of `child_process.execFile`, but returns the exit code instead of throwing an error.
 */
export async function execFileWithExitCode(file: string, args?: string[], options?: cp.ExecFileOptions): Promise<{ stdout: string; stderr: string; exitCode: number }> {
  return new Promise((resolve, reject) => {
    cp.execFile(file, args, options, (error, stdout, stderr) => {
      if (typeof stdout !== "string") stdout = stdout.toString();
      if (typeof stderr !== "string") stderr = stderr.toString();
      if (error) {
        if (typeof error.code === "string") {
          reject(error);
        } else {
          resolve({ stdout, stderr, exitCode: error.code! });
        }
      } else {
        resolve({ stdout, stderr, exitCode: 0 });
      }
    });
  });
}

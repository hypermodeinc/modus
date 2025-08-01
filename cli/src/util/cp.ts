/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import util from "node:util";
import cp from "node:child_process";
import path from "node:path";

/**
 * Promisified version of `child_process.execFile`.
 */
const execFile = util.promisify(cp.execFile);
export { execFile };

/**
 * Promisified version of `child_process.exec`.
 */
const exec = util.promisify(cp.exec);
export { exec };

/**
 * Promisified version of `child_process.execFile`, but returns the exit code instead of throwing an error.
 */
export async function execFileWithExitCode(file: string, args?: string[], options?: cp.ExecFileOptions): Promise<{ stdout: string; stderr: string; exitCode: number }> {
  return new Promise((resolve, reject) => {
    if (options?.shell && path.isAbsolute(file)) {
      file = `"${file}"`;
    }

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

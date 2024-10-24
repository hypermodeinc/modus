/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import os from "node:os";
import path from "node:path";
import process from "node:process";
import { execFile } from "./cp.js";

export async function extract(archivePath: string, installDir: string, ...args: string[]): Promise<void> {
  let tarCmd = "tar";

  if (os.platform() === "win32") {
    // On Windows, make sure we're using the OS's tar command,
    // to avoid issues with the tar command that comes with Git Bash.
    tarCmd = path.join(process.env.SystemRoot!, "System32", "tar.exe");
  }

  await execFile(tarCmd, ["-xf", archivePath, "-C", installDir, ...args]);
}

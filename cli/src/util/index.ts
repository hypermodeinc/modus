/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import chalk from "chalk";
import ora, { Ora } from "ora";

import os from "node:os";
import path from "node:path";
import readline from "node:readline";
import { isatty } from "node:tty";
import { Readable } from "node:stream";
import { finished } from "node:stream/promises";
import { spawnSync } from "node:child_process";
import { createWriteStream } from "node:fs";
import * as fs from "./fs.js";

// Expand ~ to the user's home directory
export function expandHomeDir(filePath: string): string {
  if (filePath.startsWith("~")) {
    return path.normalize(path.join(os.homedir(), filePath.slice(1)));
  }

  return path.normalize(filePath);
}

export function isRunnable(cmd: string): boolean {
  const shell = spawnSync(cmd);
  if (!shell) return false;
  return true;
}

export async function withSpinner<T>(text: string, fn: (spinner: Ora) => Promise<T>): Promise<T> {
  const spinner = ora({
    color: "white",
    text: text,
  }).start();

  try {
    return await fn(spinner);
  } finally {
    spinner.stop();
  }
}

export async function ask(question: string, placeholder?: string): Promise<string> {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  try {
    return await new Promise<string>((res, _) => {
      rl.question(question + (placeholder ? " " + placeholder + " " : ""), (answer) => {
        res(answer);
      });
    });
  } finally {
    rl.close();
  }
}

export function clearLine(n: number = 1): void {
  const stream = process.stdout;
  if (isatty(stream.fd)) {
    for (let i = 0; i < n; i++) {
      stream.write("\u001B[1A");
      stream.write("\u001B[2K");
      stream.write("\u001B[2K");
    }
  } else {
    stream.write("\n");
  }
}

export async function downloadFile(url: string, dest: string): Promise<boolean> {
  const res = await fetch(url);
  if (!res.ok) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not download file."));
    console.log(chalk.dim("   url : " + url));
    console.log(chalk.dim(`result : ${res.status} ${res.statusText}`));
    return false;
  }

  const dir = path.dirname(dest);
  if (!(await fs.exists(dir))) {
    await fs.mkdir(dir, { recursive: true });
  }

  const fileStream = createWriteStream(dest);

  // @ts-ignore
  await finished(Readable.fromWeb(res.body).pipe(fileStream));
  return true;
}

let online: boolean | undefined;
export async function isOnline(): Promise<boolean> {
  // Cache this, as we only need to check once per use of any CLI command that requires it.
  if (online !== undefined) return online;
  const controller = new AbortController();
  const timeout = setTimeout(() => controller.abort(), 1000);
  try {
    const headers = getGitHubApiHeaders();
    const response = await fetch("https://api.github.com", { signal: controller.signal, headers });
    online = response.ok;
  } catch {
    online = false;
  } finally {
    clearTimeout(timeout);
  }
  return online;
}

export function getGitHubApiHeaders() {
  const headers = new Headers({
    Accept: "application/vnd.github.v3+json",
    "X-GitHub-Api-Version": "2022-11-28",
    "User-Agent": "Modus CLI",
  });
  if (process.env.GITHUB_TOKEN) {
    headers.append("Authorization", `Bearer ${process.env.GITHUB_TOKEN}`);
  }
  return headers;
}

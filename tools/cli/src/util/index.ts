/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command } from "@oclif/core";
import chalk from "chalk";

import { spawnSync } from "node:child_process";
import { createWriteStream, existsSync, mkdirSync } from "node:fs";
import { homedir } from "node:os";
import path from "node:path";
import readline from "node:readline";
import { Readable } from "node:stream";
import { finished } from "node:stream/promises";

export async function ensureDir(dir: string): Promise<void> {
  if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
}

// Expand ~ to the user's home directory
export function expandHomeDir(filePath: string): string {
  if (filePath.startsWith("~")) {
    return path.normalize(path.join(homedir(), filePath.slice(1)));
  }

  return path.normalize(filePath);
}

export function isRunnable(cmd: string): boolean {
  const shell = spawnSync(cmd);
  if (!shell) return false;
  return true;
}

export async function cloneRepo(url: string, pth: string): Promise<boolean> {
  // should download .zip curl https://github.com/<org>/<repo>/archive/refs/heads/<branch>.zip
  // or download a release
  return true;
}

export async function withReadline<T>(fn: (rl: readline.Interface) => Promise<T>): Promise<T> {
  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
  });

  try {
    return await fn(rl);
  } finally {
    rl.close();
  }
}

export function ask(question: string, rl: readline.Interface, placeholder?: string): Promise<string> {
  return new Promise<string>((res, _) => {
    rl.question(question + (placeholder ? " " + placeholder + " " : ""), (answer) => {
      res(answer);
    });
  });
}

export function clearLine(): void {
  process.stdout.write(`\u001B[1A`);
  process.stdout.write("\u001B[2K");
  process.stdout.write("\u001B[0G");
}

export function checkVersion(instance: Command) {
  const outdated = false;
  if (outdated) {
    console.log(chalk.bgBlueBright(" INFO ") + chalk.dim(": You are running an outdated version of the Modus SDK! Please set your sdk version to stable"));
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

  const fileStream = createWriteStream(dest);

  // @ts-ignore
  await finished(Readable.fromWeb(res.body).pipe(fileStream));
  return true;
}

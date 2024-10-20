/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { spawnSync } from "node:child_process";
import { createWriteStream, existsSync, mkdirSync } from "node:fs";
import path from "node:path";
import { Interface } from "node:readline";
import { CLI_VERSION } from "../custom/globals.js";
import { Command } from "@oclif/core";
import chalk from "chalk";
import { Readable } from "node:stream";
import { finished } from "node:stream/promises";
import { rm } from "node:fs/promises";

export async function ensureDir(dir: string): Promise<void> {
  if (!existsSync(dir)) mkdirSync(dir, { recursive: true });
}

// Expand ~ to the user's home directory
export function expandHomeDir(filePath: string): string {
  if (filePath.startsWith("~")) {
    return path.normalize(path.join(process.env.HOME || "", filePath.slice(1)));
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

export function ask(question: string, rl: Interface, placeholder?: string): Promise<string> {
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

export function getLatestCLI(): string {
  // implement logic later
  return CLI_VERSION;
}

export function checkVersion(instance: Command) {
  const outdated = false;
  if (outdated) {
    console.log(chalk.bgBlueBright(" INFO ") + chalk.dim(": You are running an outdated version of the Modus SDK! Please set your sdk version to stable"));
  }
}

export async function downloadFile(url: string, dest: string) {
  const res = await fetch(url);
  if (!res.ok) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not download latest! Please check your internet connection and try again."));
    process.exit(0);
  }
  if (existsSync(dest)) await rm(dest, { recursive: true });
  const fileStream = createWriteStream(dest, { flags: "wx" });
  // @ts-ignore
  await finished(Readable.fromWeb(res.body).pipe(fileStream));
}

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
import dns from "node:dns";
import path from "node:path";
import { Readable } from "node:stream";
import { finished } from "node:stream/promises";
import { createWriteStream } from "node:fs";
import * as http from "./http.js";
import * as fs from "./fs.js";

export async function withSpinner<T>(text: string, fn: (spinner: Ora) => Promise<T>): Promise<T> {
  // NOTE: Ora comes with "oraPromise", but it doesn't clear the original text on completion.
  // Thus, we use this custom async function to ensure the spinner is applied correctly.
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

export async function downloadFile(url: string, dest: string): Promise<boolean> {
  const res = await http.get(url, false);
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

  // @ts-expect-error necessary
  await finished(Readable.fromWeb(res.body).pipe(fileStream));
  return true;
}

let online: boolean | undefined;
export async function isOnline(): Promise<boolean> {
  // Cache this, as we only need to check once per use of any CLI command that requires it.
  if (online !== undefined) return online;

  try {
    await dns.promises.lookup("releases.hypermode.com");
    online = true;
  } catch {
    online = false;
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

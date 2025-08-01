/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import process from "node:process";
import { execFile } from "./cp.js";

const EXEC_OPTIONS = {
  shell: true,
  env: process.env,
};

export async function getGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("go", ["version"], EXEC_OPTIONS);
    const parts = result.stdout.split(" ");
    const str = parts.length > 2 ? parts[2] : undefined;
    if (str?.startsWith("go")) {
      const ver = str.slice(2);

      // if version is two parts, add a .0 to make it semver compatible
      return ver.split(".").length === 2 ? ver + ".0" : ver;
    }
  } catch {
    /* empty */
  }
}

export async function getTinyGoVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("tinygo", ["version"], EXEC_OPTIONS);
    const parts = result.stdout.split(" ");
    const ver = parts.length > 2 ? parts[2] : undefined;
    if (!ver) return undefined;

    // if version is two parts, add a .0 to make it semver compatible
    return ver.split(".").length === 2 ? ver + ".0" : ver;
  } catch {
    /* empty */
  }
}

export async function getBinaryenVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("wasm-opt", ["--version"], EXEC_OPTIONS);
    const parts = result.stdout.split(" ");
    return parts.length > 2 ? parts[2] : undefined;
  } catch {
    /* empty */
  }
}

export async function getNPMVersion(): Promise<string | undefined> {
  try {
    const result = await execFile("npm", ["--version"], EXEC_OPTIONS);
    const parts = result.stdout.split(" ");
    return parts.length > 0 ? parts[0].trim() : undefined;
  } catch {
    /* empty */
  }
}

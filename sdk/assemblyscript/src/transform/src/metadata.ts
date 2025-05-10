/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { readFileSync } from "fs";
import { fileURLToPath } from "url";
import { execSync } from "child_process";
import * as path from "path";
import { Xid } from "xid-ts";
import binaryen from "assemblyscript/lib/binaryen.js";
import chalk from "chalk";
import { FunctionSignature, TypeDefinition } from "./types.js";

const METADATA_VERSION = 2;

export class Metadata {
  public plugin: string;
  public module: string;
  public sdk: string;
  public buildId: string;
  public buildTs: string;
  public gitRepo?: string;
  public gitCommit?: string;
  public fnExports: { [key: string]: FunctionSignature } = {};
  public fnImports: { [key: string]: FunctionSignature } = {};
  public types: { [key: string]: TypeDefinition } = {};

  static generate(): Metadata {
    const m = new Metadata();

    m.buildId = new Xid().toString();
    m.buildTs = new Date().toISOString();
    m.plugin = getPluginInfo();
    m.sdk = getSdkInfo();

    if (isGitRepo()) {
      const gitRepo = getGitRepo();
      if (gitRepo) {
        m.gitRepo = getGitRepo();
      }

      const gitCommit = getGitCommit();
      if (gitCommit) {
        m.gitCommit = getGitCommit();
      }
    }

    return m;
  }

  addExportFn(functions: FunctionSignature[]) {
    for (const fn of functions) {
      const name = fn.name;
      this.fnExports[name] = fn;
    }
  }

  addImportFn(functions: FunctionSignature[]) {
    for (const fn of functions) {
      this.fnImports[fn.name] = fn;
    }
  }

  addTypes(types: TypeDefinition[]) {
    for (const t of types) {
      this.types[t.name] = t;
    }
  }

  writeToModule(module: binaryen.Module) {
    const encoder = new TextEncoder();

    const fnExports = this.fnExports;
    const fnImports = this.fnImports;

    const json = JSON.stringify(this);

    this.fnExports = fnExports;
    this.fnImports = fnImports;

    module.addCustomSection(
      "modus_metadata_version",
      Uint8Array.from([METADATA_VERSION]),
    );

    module.addCustomSection("modus_metadata", encoder.encode(json));
  }

  logResults() {
    const writeHeader = (text: string) => {
      console.log(chalk.bold.blue(text));
    };

    const writeItem = (text: string) => {
      console.log(`  ${chalk.cyan(text)}`);
    };

    const writeTable = (rows: string[][]) => {
      rows = rows.filter((r) => !!r);
      const pad = rows.reduce(
        (max, row) => [
          Math.max(max[0], row[0].length),
          Math.max(max[1], row[1].length),
        ],
        [0, 0],
      );

      rows.forEach((row) => {
        if (row) {
          const padding = " ".repeat(pad[0] - row[0].length);
          const key = chalk.cyan(row[0] + ":");
          const value = chalk.blue(row[1]);
          console.log(`  ${key}${padding} ${value}`);
        }
      });
    };

    writeHeader("Metadata:");
    writeTable([
      ["Plugin Name", this.plugin],
      ["Modus SDK", this.sdk],
      ["Build ID", this.buildId],
      ["Build Timestamp", this.buildTs],
      this.gitRepo ? ["Git Repository", this.gitRepo] : undefined,
      this.gitCommit ? ["Git Commit", this.gitCommit] : undefined,
    ]);
    console.log();

    const functions = Object.values(this.fnExports).filter(
      (f) => !f.isHidden(),
    );
    if (functions.length > 0) {
      writeHeader("Functions:");
      functions.forEach((f) => writeItem(f.toString()));
      console.log();
    }

    const types = Object.values(this.types).filter((t) => !t.isHidden());
    if (types.length > 0) {
      writeHeader("Custom Types:");
      types.forEach((t) => writeItem(t.toString()));
      console.log();
    }

    if (process.env.MODUS_DEBUG) {
      writeHeader("Metadata JSON:");
      console.log(JSON.stringify(this, undefined, 2) + "\n\n");
    }
  }
}

function getSdkInfo(): string {
  const filePath = path.join(
    path.dirname(fileURLToPath(import.meta.url)),
    "..",
    "..",
    "package.json",
  );
  const json = readFileSync(filePath).toString();
  const lib = JSON.parse(json);
  return `${lib.name.split("/")[1]}@${lib.version || "dev"}`;
}

function getPluginInfo(): string {
  const pluginName = process.env.npm_package_name;
  const pluginVersion = process.env.npm_package_version;

  if (!pluginName) {
    throw new Error("Missing name in package.json");
  }

  if (!pluginVersion) {
    // versionless plugins are allowed
    return pluginName;
  }

  return `${pluginName}@${pluginVersion}`;
}

function isGitRepo(): boolean {
  try {
    // This will throw if not in a git repo, or if git is not installed.
    execSync("git rev-parse --is-inside-work-tree", { stdio: "ignore" });
    return true;
  } catch {
    return false;
  }
}

function getGitRepo(): string | undefined {
  try {
    let url = execSync("git remote get-url origin", {
      stdio: ["ignore", "pipe", "ignore"],
    })
      .toString()
      .trim();

    // Convert ssh to https
    if (url.startsWith("git@")) {
      url = url.replace(":", "/").replace("git@", "https://");
    }

    // Remove the .git suffix
    if (url.endsWith(".git")) {
      url = url.slice(0, -4);
    }

    return url;
  } catch {
    return undefined;
  }
}

function getGitCommit(): string | undefined {
  try {
    return execSync("git rev-parse HEAD", {
      stdio: ["ignore", "pipe", "ignore"],
    })
      .toString()
      .trim();
  } catch {
    return undefined;
  }
}

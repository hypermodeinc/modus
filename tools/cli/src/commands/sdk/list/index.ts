/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import { existsSync, readdirSync } from "node:fs";
import { expandHomeDir } from "../../../util/index.js";

export default class SDKListCommand extends Command {
  static args = {};
  static description = "List installed SDK versions";
  static examples = ["modus sdk list"];
  static flags = {};

  async run(): Promise<void> {
    let versions: string[] = [];

    if (!existsSync(expandHomeDir("~/.modus/sdk/"))) {
      this.log("No versions installed!");
      process.exit(0);
    }

    try {
      versions = readdirSync(expandHomeDir("~/.modus/sdk")).reverse();
    } catch {
      versions = [];
    }

    if (!versions.length) {
      this.log("No versions installed!");
      return;
    }

    for (const version of versions) {
      this.log(version);
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

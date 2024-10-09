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
import os from "node:os";
import { expandHomeDir } from "../../../util/index.js";
import { existsSync, readdirSync, rmSync } from "node:fs";

export default class SDKRemoveCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to remove",
      hidden: false,
      required: false,
    }),
  };
  static description = "Remove a specific SDK version";
  static examples = ["modus sdk remove v0.0.0", "modus sdk remove all"];
  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKRemoveCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk remove <version>"), process.exit(0);
    const isDev = args.version.startsWith("dev-") || args.version.startsWith("link");
    let version = isDev ? args.version : args.version?.trim().toLowerCase().replace("v", "");
    if (!existsSync(expandHomeDir("~/.modus/sdk/"))) {
      this.log("No versions installed!");
      process.exit(0);
    }
    const versions = readdirSync(expandHomeDir("~/.modus/sdk/"));
    if (!versions.length) {
      this.log("No versions installed!");
      process.exit(0);
    }
    if (version === "all") {
      for (const version of versions) {
        rmSync(expandHomeDir("~/.modus/sdk/" + version), { recursive: true, force: true });
      }
      this.log("Removed all SDK versions");
      process.exit(0);
    } else {
      rmSync(expandHomeDir("~/.modus/sdk/" + version), { recursive: true, force: true });
      this.log("Removed Modus " + (isDev ? "" : "v") + version);
      process.exit(0);
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

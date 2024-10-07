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
import { readdirSync, rmSync } from "node:fs";

const versions = ["0.12.0", "0.12.1", "0.12.2", "0.12.3", "0.12.4", "0.12.5", "0.12.6"];
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
    let version = args.version?.trim().toLowerCase().replace("v", "");
    const platform = os.platform();
    const arch = os.arch();
    const file = "modus-runtime-v" + version + "-" + platform + "-" + arch + (platform === "win32" ? ".exe" : "");
    const versions = readdirSync(expandHomeDir("~/.hypermode/sdk/"));
    if (!versions.length) {
      this.log("No versions installed!");
      process.exit(0);
    }
    if (version === "all") {
      for (const version of versions) {
        rmSync(expandHomeDir("~/.hypermode/sdk" + version), { recursive: true, force: true });
      }
      this.log("Removed all SDK versions");
      process.exit(0);
      return;
    } else {
      rmSync(expandHomeDir("~/.hypermode/sdk" + version), { recursive: true, force: true });
      this.log("Removed Modus v" + version);
      process.exit(0);
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

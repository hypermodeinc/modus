/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import { createInterface } from "node:readline";
import path from "node:path";
import chalk from "chalk";

import * as fs from "../../../util/fs.js";
import * as vi from "../../../util/versioninfo.js";
import * as globals from "../../../custom/globals.js";
import { ask, clearLine } from "../../../util/index.js";

export default class SDKRemoveCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to remove",
      hidden: false,
      required: false,
    }),
  };

  static flags = {
    force: Flags.boolean({
      char: "f",
      description: "Remove without prompting",
    }),
  };

  static description = "Remove a Modus SDK";
  static examples = ["modus sdk remove v0.0.0", "modus sdk remove all"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKRemoveCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk remove <version>"), this.exit(0);

    const installedVersions = await vi.getInstalledVersions();
    if (installedVersions.length == 0) {
      this.log("No versions installed!");
      this.exit(0);
    }

    const rl = createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    if (args.version.toLowerCase() === "all") {
      if (!flags.force && !(await this.confirmAction(rl, "Really, remove all Modus SDK versions? [y/n]"))) {
        clearLine();
        return;
      }

      await fs.rm(globals.ModusHomeDir, { recursive: true, force: true });
      this.log("Removed all Modus SDK versions");
      this.exit(0);
    }

    if (!installedVersions.includes(args.version)) {
      this.logError("Specified version is not installed!");
      this.exit(1);
    }

    if (!flags.force && !(await this.confirmAction(rl, `Really, remove Modus SDK version ${args.version}? [y/n]`))) {
      clearLine();
      return;
    }

    const dir = path.join(globals.ModusHomeDir, "sdk", args.version);
    await fs.rm(dir, { recursive: true, force: true });
    this.log(`Removed Modus SDK version ${args.version}`);
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async confirmAction(rl: ReturnType<typeof createInterface>, message: string): Promise<boolean> {
    this.log(message);
    const cont = ((await ask(chalk.dim(" -> "), rl)) || "n").toLowerCase().trim();
    clearLine();
    return cont === "yes" || cont === "y";
  }
}

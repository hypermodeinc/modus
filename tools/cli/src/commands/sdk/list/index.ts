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
import * as vi from "../../../util/versioninfo.js";

export default class SDKListCommand extends Command {
  static args = {};
  static description = "List installed Modus SDKs";
  static examples = ["modus sdk list"];
  static flags = {};

  async run(): Promise<void> {
    const versions = await vi.getInstalledVersions();
    if (versions.length === 0) {
      this.log(chalk.yellow("No Modus SDK versions installed!"));
    }

    this.log(chalk.cyan(chalk.bold("Installed Modus SDK Versions:")));
    for (const version of versions) {
      this.log(version);
    }
  }
}

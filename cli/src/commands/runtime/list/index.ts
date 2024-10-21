/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command, Flags } from "@oclif/core";
import chalk from "chalk";
import * as vi from "../../../util/versioninfo.js";
import { getHeader } from "../../../custom/header.js";

export default class RuntimeListCommand extends Command {
  static args = {};
  static description = "List installed Modus runtimes";
  static examples = ["modus runtime list"];
  static flags = {
    nologo: Flags.boolean({ hidden: true }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(RuntimeListCommand);

    if (!flags.nologo) {
      this.log(getHeader(this.config.version));
    }

    await this.showInstalledRuntimes();
  }

  private async showInstalledRuntimes(): Promise<void> {
    const versions = await vi.getInstalledRuntimeVersions();
    if (versions.length > 0) {
      this.log(chalk.bold.cyan("Installed Runtimes:"));
      for (const version of versions) {
        this.log(`• Modus Runtime ${version}`);
      }
    } else {
      this.log(chalk.yellow("No Modus runtimes installed"));
    }

    this.log();
  }
}

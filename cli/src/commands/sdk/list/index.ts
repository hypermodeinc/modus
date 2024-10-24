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
import { SDK } from "../../../custom/globals.js";

export default class SDKListCommand extends Command {
  static args = {};
  static description = "List installed Modus SDKs";
  static examples = ["modus sdk list"];
  static flags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    "no-logo": Flags.boolean({
      aliases: ["nologo"],
      hidden: true,
    }),
  };

  async run(): Promise<void> {
    const { flags } = await this.parse(SDKListCommand);

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    await this.showInstalledSDKs();
  }

  private async showInstalledSDKs(): Promise<void> {
    let found = false;
    for (const sdk of Object.values(SDK)) {
      const versions = await vi.getInstalledSdkVersions(sdk);
      if (versions.length === 0) {
        continue;
      }
      if (!found) {
        this.log(chalk.bold.cyan("Installed SDKs:"));
        found = true;
      }

      for (const version of versions) {
        this.log(`• Modus ${sdk} SDK ${version}`);
      }
    }

    if (!found) {
      this.log(chalk.yellow("No Modus SDKs installed"));
    }

    this.log();
  }
}

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import chalk from "chalk";
import * as vi from "../../../util/versioninfo.js";
import { getHeader } from "../../../custom/header.js";
import { BaseCommand } from "../../../baseCommand.js";

export default class RuntimeListCommand extends BaseCommand {
  static args = {};
  static description = "List installed Modus runtimes";
  static examples = ["modus runtime list"];
  static flags = {};

  async run(): Promise<void> {
    const { flags } = await this.parse(RuntimeListCommand);

    if (!flags["no-logo"]) {
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

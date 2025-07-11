/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import chalk from "chalk";
import * as vi from "../../../util/versioninfo.js";
import { getHeader } from "../../../custom/header.js";
import { SDK } from "../../../custom/globals.js";
import { BaseCommand } from "../../../baseCommand.js";

export default class SDKListCommand extends BaseCommand {
  static args = {};
  static description = "List installed Modus SDKs";
  static examples = ["modus sdk list"];
  static flags = {};

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

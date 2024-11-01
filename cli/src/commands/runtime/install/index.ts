/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import chalk from "chalk";
import * as fs from "../../../util/fs.js";
import * as vi from "../../../util/versioninfo.js";
import * as installer from "../../../util/installer.js";
import { isOnline, withSpinner } from "../../../util/index.js";
import { getHeader } from "../../../custom/header.js";
import { BaseCommand } from "../../../baseCommand.js";

export default class RuntimeInstallCommand extends BaseCommand {
  static args = {
    version: Args.string({
      description: "Runtime version to install",
      default: "latest",
    }),
  };

  static description = "Install a Modus runtime";
  static examples = ["modus runtime install", "modus runtime install v1.2.3", "modus runtime install --prerelease"];

  static flags = {
    force: Flags.boolean({
      char: "f",
      default: false,
      description: "Force re-installation if version already exists",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      default: false,
      description: "Install a prerelease version (used with 'latest' version)",
    }),
  };

  async run(): Promise<void> {
    const { args, flags } = await this.parse(RuntimeInstallCommand);

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    if (!(await isOnline())) {
      this.logError("No internet connection.  You must be online to install a Modus runtime.");
      this.exit(1);
    }

    let runtimeVersion = args.version;
    if (runtimeVersion.toLowerCase() === "latest") {
      const version = await vi.getLatestRuntimeVersion(flags.prerelease);
      if (!version) {
        this.logError("Failed to fetch latest Modus Runtime version.");
        this.exit(1);
      }
      runtimeVersion = version;
    } else if (!runtimeVersion.startsWith("v")) {
      this.logError("Version must start with 'v'.");
      this.exit(1);
    } else if (!(await vi.runtimeReleaseExists(runtimeVersion))) {
      this.logError(`Modus Runtime ${runtimeVersion} does not exist!`);
      this.exit(1);
    }

    await this.installRuntime(runtimeVersion, flags.force);

    this.log();
    this.log(chalk.bold.cyanBright("Installation successful!"));
    this.log();
  }

  private async installRuntime(version: string, force: boolean): Promise<void> {
    const runtimeText = `Modus Runtime ${version}`;
    const installRuntime = await withSpinner(chalk.dim("Checking for runtime installation."), async (spinner) => {
      if (await vi.runtimeVersionIsInstalled(version)) {
        if (force) {
          spinner.warn(chalk.yellow(runtimeText + " was already installed. Reinstalling."));
          const dir = vi.getRuntimePath(version);
          await fs.rm(dir, { recursive: true, force: true });
          return true;
        } else {
          spinner.succeed(chalk.dim(runtimeText + " was already installed."));
          return false;
        }
      }
      return true;
    });

    if (installRuntime) {
      await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
        try {
          await installer.installRuntime(version);
        } catch (e) {
          spinner.fail(chalk.red(`Failed to download ${runtimeText}`));
          throw e;
        }
        spinner.succeed(chalk.dim(`Installed ${runtimeText}`));
      });
    }
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

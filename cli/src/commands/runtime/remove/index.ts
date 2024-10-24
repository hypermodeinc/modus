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
import { withSpinner } from "../../../util/index.js";
import * as inquirer from "@inquirer/prompts";

export default class RuntimeRemoveCommand extends Command {
  static args = {
    version: Args.string({
      description: "Runtime version to remove, or 'all' to remove all runtimes.",
      required: true,
    }),
  };

  static flags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    force: Flags.boolean({
      char: "f",
      default: false,
      description: "Remove without prompting",
    }),
  };

  static description = "Remove a Modus runtime";
  static examples = ["modus runtime remove v0.0.0", "modus runtime remove all"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(RuntimeRemoveCommand);
    if (!args.version) {
      this.logError(`No runtime version specified. Run ${chalk.whiteBright("modus runtime remove <version>")}, or ${chalk.whiteBright("modus runtime remove all")}`);
      return;
    }

    if (args.version.toLowerCase() === "all") {
      const versions = await vi.getInstalledRuntimeVersions();
      if (versions.length === 0) {
        this.log(chalk.yellow("No Modus runtimes are installed."));
        this.exit(1);
      } else if (!flags.force) {
        try {
          const confirmed = await inquirer.confirm({
            message: "Are you sure you want to remove all Modus runtimes?",
            default: false,
          });
          if (!confirmed) {
            this.abort();
          }
        } catch (err: any) {
          if (err.name === "ExitPromptError") {
            this.abort();
          }
        }
      }

      for (const version of versions) {
        await this.removeRuntime(version);
      }
    } else if (!args.version.startsWith("v")) {
      this.logError("Version must start with 'v'.");
      this.exit(1);
    } else {
      const runtimeText = `Modus Runtime ${args.version}`;
      const isInstalled = await vi.runtimeVersionIsInstalled(args.version);
      if (!isInstalled) {
        this.log(chalk.yellow(runtimeText + "is not installed."));
        this.exit(1);
      } else if (!flags.force) {
        try {
          const confirmed = await inquirer.confirm({
            message: `Are you sure you want to remove ${runtimeText}?`,
            default: false,
          });
          if (!confirmed) {
            this.abort();
          }
        } catch (err: any) {
          if (err.name === "ExitPromptError") {
            this.abort();
          }
        }
      }

      await this.removeRuntime(args.version);
    }
  }

  private async removeRuntime(version: string): Promise<void> {
    const runtimeText = `Modus Runtime ${version}`;
    await withSpinner(chalk.dim("Removing " + runtimeText), async (spinner) => {
      const dir = vi.getRuntimePath(version);
      try {
        await fs.rm(dir, { recursive: true, force: true });
        spinner.succeed(chalk.dim("Removed " + runtimeText));
      } catch (e) {
        spinner.fail(chalk.red("Failed to remove " + runtimeText));
        throw e;
      }
    });
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private abort() {
    this.log(chalk.dim("Aborted"));
    this.exit(1);
  }
}

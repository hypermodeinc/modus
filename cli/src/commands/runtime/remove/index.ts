/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command } from "@oclif/core";
import chalk from "chalk";

import * as fs from "../../../util/fs.js";
import * as vi from "../../../util/versioninfo.js";
import { withSpinner } from "../../../util/index.js";
import * as inquirer from "@inquirer/prompts";
import { isErrorWithName } from "../../../util/errors.js";

export default class RuntimeRemoveCommand extends Command {
  static args = {
    version: Args.string({
      description: "Runtime version to remove, or 'all' to remove all runtimes.",
    }),
  };

  static flags = {};

  static description = "Remove a Modus runtime";
  static examples = ["modus runtime remove v0.0.0", "modus runtime remove all"];

  async run(): Promise<void> {
    try {
      const { args, flags } = await this.parse(RuntimeRemoveCommand);

      if (!args.version) {
        const versions = await vi.getInstalledRuntimeVersions();

        if (versions.length === 0) {
          this.log(chalk.yellow("No Modus runtimes are installed."));
          this.exit(1);
        }

        const runtimeVersionToRemove = await inquirer.select({
          message: "Select a runtime version to remove",
          choices: [
            ...versions.map((v) => ({ name: v, value: v })),
            {
              name: chalk.bold("All versions"),
              value: "all",
            },
          ],
        });

        if (runtimeVersionToRemove === "all") {
          const confirmed = await inquirer.confirm({
            message: "Are you sure you want to remove all Modus runtimes?",
            default: false,
          });
          if (!confirmed) {
            this.abort();
          }

          for (const version of versions) {
            await this.removeRuntime(version);
          }
        } else {
          await this.removeRuntime(runtimeVersionToRemove);
        }

        return;
      }

      if (args.version.toLowerCase() === "all") {
        const versions = await vi.getInstalledRuntimeVersions();
        if (versions.length === 0) {
          this.log(chalk.yellow("No Modus runtimes are installed."));
          this.exit(1);
        } else if (!flags.force) {
          const confirmed = await inquirer.confirm({
            message: "Are you sure you want to remove all Modus runtimes?",
            default: false,
          });
          if (!confirmed) {
            this.abort();
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
          const confirmed = await inquirer.confirm({
            message: `Are you sure you want to remove ${runtimeText}?`,
            default: false,
          });
          if (!confirmed) {
            this.abort();
          }
        }

        await this.removeRuntime(args.version);
      }
    } catch (err) {
      if (isErrorWithName(err) && err.name === "ExitPromptError") {
        this.abort();
      } else {
        throw err;
      }
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

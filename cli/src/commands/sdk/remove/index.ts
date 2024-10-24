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
import { parseSDK, SDK } from "../../../custom/globals.js";
import { withSpinner } from "../../../util/index.js";
import * as inquirer from "@inquirer/prompts";

export default class SDKRemoveCommand extends Command {
  static args = {
    name: Args.string({
      description: "SDK name to remove (or 'all' to remove all SDKs)",
      required: true,
    }),
    version: Args.string({
      description: "SDK version to remove, if removing a specific SDK. Leave blank to remove all versions of the SDK.",
      default: "all",
    }),
  };

  static flags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    runtimes: Flags.boolean({
      char: "r",
      default: false,
      description: "Remove runtimes also. Only valid when removing 'all' SDKs.",
    }),
    force: Flags.boolean({
      char: "f",
      default: false,
      description: "Remove without prompting",
    }),
  };

  static description = "Remove a Modus SDK";
  static examples = ["modus sdk remove assemblyscript v0.0.0", "modus sdk remove all"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKRemoveCommand);
    if (!args.version) {
      this.logError(`No SDK specified! Run ${chalk.whiteBright("modus sdk remove <name> [version]")}, or ${chalk.whiteBright("modus sdk remove all")}`);
      return;
    }

    if (args.name.toLowerCase() === "all") {
      let found = false;
      for (const sdk of Object.values(SDK)) {
        const versions = await vi.getInstalledSdkVersions(sdk);
        if (versions.length > 0) {
          found = true;
          break;
        }
      }
      if (!found) {
        this.log(chalk.yellow("No Modus SDKs are installed."));
        this.exit(1);
      }

      if (!flags.force) {
        try {
          const confirmed = inquirer.confirm({
            message: "Are you sure you want to remove all Modus SDKs?",
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

      for (const sdk of Object.values(SDK)) {
        const versions = await vi.getInstalledSdkVersions(sdk);
        for (const version of versions) {
          await this.removeSDK(sdk, version);
        }
      }

      if (flags.runtimes) {
        const versions = await vi.getInstalledRuntimeVersions();
        for (const version of versions) {
          await this.removeRuntime(version);
        }
      }
    } else {
      const sdk = parseSDK(args.name);
      if (args.version.toLowerCase() === "all") {
        const versions = await vi.getInstalledSdkVersions(sdk);
        if (versions.length === 0) {
          this.log(chalk.yellow(`No Modus ${sdk} SDKs are installed.`));
          this.exit(1);
        } else if (!flags.force) {
          try {
            const confirmed = inquirer.confirm({
              message: `Are you sure you want to remove all Modus ${sdk} SDKs?`,
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
          await this.removeSDK(sdk, version);
        }
      } else if (!args.version.startsWith("v")) {
        this.logError("Version must start with 'v'.");
        this.exit(1);
      } else {
        const sdkText = `Modus ${sdk} SDK ${args.version}`;
        const isInstalled = await vi.sdkVersionIsInstalled(sdk, args.version);
        if (!isInstalled) {
          this.log(chalk.yellow(sdkText + "is not installed."));
          this.exit(1);
        } else if (!flags.force) {
          try {
            const confirmed = inquirer.confirm({
              message: `Are you sure you want to remove ${sdkText}?`,
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

        await this.removeSDK(sdk, args.version);
      }
    }
  }

  private async removeSDK(sdk: SDK, version: string): Promise<void> {
    const sdkText = `Modus ${sdk} SDK ${version}`;
    await withSpinner(chalk.dim("Removing " + sdkText), async (spinner) => {
      const dir = vi.getSdkPath(sdk, version);
      try {
        await fs.rm(dir, { recursive: true, force: true });
        spinner.succeed(chalk.dim("Removed " + sdkText));
      } catch (e) {
        spinner.fail(chalk.red("Failed to remove " + sdkText));
        throw e;
      }
    });
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

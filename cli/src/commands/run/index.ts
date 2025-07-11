/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Flags } from "@oclif/core";
import { BaseCommand } from "../../baseCommand.js";
import DevCommand from "../dev/index.js";

export default class RunCommand extends BaseCommand {
  static args = {
    path: Args.directory({
      description: "Path to app directory",
      default: ".",
      exists: true,
    }),
  };

  static flags = {
    runtime: Flags.string({
      char: "r",
      description: "Modus runtime version to use. If not provided, the latest runtime compatible with the app will be used.",
    }),
    prerelease: Flags.boolean({
      char: "p",
      aliases: ["pre"],
      description: "Use a prerelease version of the Modus runtime. Not needed if specifying a runtime version.",
    }),
    build: Flags.boolean({
      char: "b",
      description: "Build the app before running (or when watching for changes)",
      default: true,
    }),
    watch: Flags.boolean({
      char: "w",
      description: "Watch app code for changes",
      default: false,
    }),
    delay: Flags.integer({
      description: "Delay (in milliseconds) between file change detection and rebuild",
      default: 500,
    }),
  };

  static description = "Run a Modus app locally";

  static examples = ["modus run", "modus run ./my-app", "modus run ./my-app --no-watch"];

  async run(): Promise<void> {
    const { flags } = await this.parse(RunCommand);
    const argv = this.argv.filter((arg) => arg != "--watch" && arg != "-w" && arg != "--build" && arg != "-b");

    if (!flags.watch) argv.push("--no-watch");
    if (!flags.build) argv.push("--no-build");

    await DevCommand.run(argv);
  }
}

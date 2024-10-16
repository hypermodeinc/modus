/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import { quote } from "shell-quote";
import { execSync } from "node:child_process";
import path from "node:path";
import { SDK } from "../../custom/globals.js";
import { isRunnable } from "../../util/index.js";
import { fileURLToPath } from "node:url";
import { existsSync } from "node:fs";

const NPM_CMD = isRunnable("npm") ? "npm" : path.join(path.dirname(fileURLToPath(import.meta.url)), "../../../bin/node-bin/bin/npm");
export default class BuildCommand extends Command {
  static args = {
    path: Args.string({
      description: "./my-project-|-Directory to build",
      hidden: false,
      required: false,
    }),
  };

  static description = "Build a Modus project";

  static examples = ["modus build ./my-project"];

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(BuildCommand);

    const cwd = args.path ? path.join(process.cwd(), args.path) : process.cwd();
    const sdk = SDK.AssemblyScript;
    if (!isRunnable(NPM_CMD)) {
      this.logError("Could not locate NPM. Please install and try again!");
      process.exit(0);
    }

    if (!existsSync(path.join(cwd, "/node_modules"))) {
      this.logError("Dependencies are not installed! Please install dependencies with npm i");
      process.exit(0);
    }

    // should pass argv
    if (sdk === SDK.AssemblyScript) {
      execSync(quote([NPM_CMD, "run", "build"]), { cwd, stdio: "inherit" });
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

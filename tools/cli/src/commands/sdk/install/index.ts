/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import chalk from "chalk";
import { cpSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import { expandHomeDir } from "../../../util/index.js";
import { Metadata } from "../../../util/metadata.js";
import { execSync } from "node:child_process";

const versions = ["0.12.0", "0.12.1", "0.12.2", "0.12.3", "0.12.4", "0.12.5", "0.12.6"];
export default class SDKInstallCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to install",
      hidden: false,
      required: false,
    }),
  };

  static description = "Install a specific SDK version";

  static examples = ["modus sdk install v0.0.0", "modus sdk install latest"];

  static flags = {
    silent: Flags.boolean({
      description: "Suppress output logs",
      hidden: false,
      required: false
    })
  };

  async run(): Promise<void> {
    const { args, flags } = await this.parse(SDKInstallCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk install <version>");
    let version = args.version?.trim().toLowerCase().replace("v", "");
    const platform = os.platform();
    const arch = os.arch();
    const file = "modus-runtime-v" + version + "-" + platform + "-" + arch + (platform === "win32" ? ".exe" : "");

    const src = path.join(path.dirname(import.meta.url.replace("file:", "")), "../../../../runtime-bin/" + "modus-runtime-v0.12.6-" + platform + "-" + arch + (platform === "win32" ? ".exe" : ""));
    if (version === "all") {
      for (const version of versions) {
        cpSync(src, expandHomeDir("~/.hypermode/sdk/" + version + "/runtime" + (platform === "win32" ? ".exe" : "")));
      }
      if (!flags.silent) this.log("Installed versions 0.12.0-0.12.6");
      return;
    } else if (version === "latest") {
      version = (await Metadata.getLatestRuntime())!;
    }

    const runtimePath = expandHomeDir("~/.hypermode/sdk/" + version + "/runtime" + (platform === "win32" ? ".exe" : ""));
    cpSync(src, runtimePath);

    if (platform === "linux" || platform === "darwin") {
      execSync("chmod +x " + runtimePath, { stdio: "ignore" });
    }

    if (!flags.silent) this.log("Installed Modus v" + version);
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

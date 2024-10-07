/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import { expandHomeDir } from "../../util/index.js";
import { Metadata } from "../../util/metadata.js";
import BuildCommand from "../build/index.js";
import path from "path";
import { copyFileSync, existsSync, readFileSync, watch as watchFolder } from "fs";
import chalk from "chalk";
import { spawn } from "child_process";
import os from "node:os";

export default class Run extends Command {
  static args = {
    path: Args.string({
      description: "./my-project-|-Path to project directory",
      hidden: false,
      required: false
    }),
  };

  static flags = {
    watch: Flags.boolean({
      char: "w",
      description: "Watch project and rebuild continually",
      hidden: false,
      required: false,
    }),
    build: Flags.boolean({
      char: "b",
      description: "Build the latest before running",
      hidden: false,
      required: false,
    }),
    silent: Flags.boolean({
      char: "s",
      description: "Suppress output logs from cluttering terminal",
      hidden: false,
      required: false,
    }),
    verbose: Flags.boolean({
      char: "v",
      description: "Enable descriptive logging",
      hidden: false,
      required: false,
    }),
    freq: Flags.integer({
      char: "f",
      description: "Frequency to check for changes",
      hidden: false,
      required: false,
      default: 3000
    }),
  };

  static description = "Launch a Modus app to local development";

  static examples = [`<%= config.bin %> <%= command.id %> run ./project-path --watch`];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Run);
    const runtimePath = expandHomeDir("~/.hypermode/sdk/" + Metadata.runtime_version + "/runtime") + (os.platform() === "win32" ? ".exe" : "");

    const cwd = args.path ? path.join(process.cwd(), args.path) : process.cwd();
    const watch = flags.watch;

    if (!existsSync(path.join(cwd))) {
      this.logError("Could not target folder! Please try again");
      process.exit(0);
    }

    if (!existsSync(path.join(cwd, "/node_modules"))) {
      this.logError("Dependencies not installed! Please install dependencies by running `npm i` and try again");
      process.exit(0);
    }

    if (!existsSync(path.join(cwd, "/package.json"))) {
      this.logError("Could not locate package.json! Please try again");
      process.exit(0);
    }

    if (!existsSync(runtimePath)) {
      this.logError("Modus Runtime v" + Metadata.runtime_version + " not installed! Run `modus sdk install " + Metadata.runtime_version + "` and try again!");
      process.exit(0);
    }

    let project_name: string;
    try {
      project_name = JSON.parse(readFileSync(path.join(cwd, "/package.json")).toString()).name;
    } catch {
      this.logError("Could not read package.json! Please try again");
      process.exit(0);
    }

    try {
      if (flags.build) await BuildCommand.run(args.path ? [args.path] : []);
    } catch { }
    const build_wasm = path.join(cwd, "/build/" + project_name + ".wasm");
    const deploy_wasm = expandHomeDir("~/.hypermode/" + project_name + ".wasm");
    copyFileSync(build_wasm, deploy_wasm);

    spawn(runtimePath, {
      stdio: "inherit", env: {
        ...process.env,
        MODUS_ENV: "dev"
      }
    });

    if (watch) {
      const delay = flags.freq;
      let lastModified = 0;
      let lastBuild = 0;
      let paused = true;

      setInterval(async () => {
        if (paused) return;
        if (lastBuild > lastModified) {
          paused = true;
          return;
        }
        if (Date.now() - lastModified > delay * 2) paused = true;
        lastBuild = Date.now();
        try {
          await BuildCommand.run(args.path ? [args.path] : []);
        } catch { }
        copyFileSync(build_wasm, deploy_wasm);
      }, delay);

      watchFolder(path.join(cwd, "/assembly"), () => {
        lastModified = Date.now();
        paused = false;
      });
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Command, Flags } from "@oclif/core";
import { expandHomeDir, isRunnable } from "../../util/index.js";
import BuildCommand from "../build/index.js";
import path from "path";
import { copyFileSync, existsSync, readdirSync, readFileSync, watch as watchFolder } from "fs";
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
    runtime: Flags.string({
      char: "r",
      description: "Runtime to use",
      hidden: false,
      required: false,
      default: getLatestRuntime()
    })
  };

  static description = "Launch a Modus app to local development";

  static examples = [`<%= config.bin %> <%= command.id %> run ./project-path --watch`];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Run);
    const isDev = flags.runtime.startsWith("dev-") || flags.runtime.startsWith("link");
    const runtimePath = expandHomeDir("~/.modus/sdk/" + flags.runtime) + (isDev ? "" : "/runtime" + (os.platform() === "win32" ? ".exe" : ""));

    const cwd = args.path ? path.join(process.cwd(), args.path) : process.cwd();
    const watch = flags.watch;

    if (!existsSync(path.join(cwd))) {
      this.logError("Could not target folder! Please try again");
      process.exit(0);
    }

    // TODO: Check the type of SDK we are running

    if (!existsSync(path.join(cwd, "/node_modules"))) {
      this.logError("Dependencies not installed! Please install dependencies by running `npm i` and try again");
      process.exit(0);
    }

    if (!existsSync(path.join(cwd, "/package.json"))) {
      this.logError("Could not locate package.json! Please try again");
      process.exit(0);
    }

    let install_cmd = flags.runtime;
    if (isDev) {
      if (install_cmd.startsWith("link")) {
        install_cmd = "--link ./path-to-modus";
      } else if (install_cmd.startsWith("dev-")) {
        install_cmd = "--branch " + install_cmd.split("-")[1] + " --commit " + install_cmd.split("-")[2];
      }
    }

    if (!existsSync(runtimePath)) {
      this.logError("Modus Runtime " + (isDev ? "" : "v") + flags.runtime + " not installed!\n Run `modus sdk install " + install_cmd + "` and try again!");
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
    const deploy_wasm = expandHomeDir("~/.modus/" + project_name + ".wasm");
    copyFileSync(build_wasm, deploy_wasm);

    if (isDev) {
      if (!isRunnable("go")) {
        this.logError("Cannot find any valid versions of Go! Please install go")
      }
      spawn("go run .", {
        cwd: runtimePath,
        stdio: "inherit",
        env: {
          ...process.env,
          MODUS_ENV: "dev"
        }
      });
    } else {
      spawn(runtimePath, {
        stdio: "inherit",
        env: {
          ...process.env,
          MODUS_ENV: "dev"
        }
      });
    }

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

function getLatestRuntime(): string | undefined {
  let versions: string[] = [];
  try {
    versions = readdirSync(expandHomeDir("~/.modus/sdk")).reverse().filter(v => !v.startsWith("dev-") && !v.startsWith("link"));
  } catch { }
  if (!versions.length) return undefined;
  return versions[0];
}
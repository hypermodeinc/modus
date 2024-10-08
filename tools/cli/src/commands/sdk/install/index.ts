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
import { cpSync, existsSync, mkdirSync, readdirSync, rmdirSync, rmSync, statSync, symlinkSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import { ask, clearLine, downloadFile, expandHomeDir } from "../../../util/index.js";
import { execSync } from "node:child_process";
import { rm } from "node:fs/promises";
import { createInterface } from "node:readline";
import { ParserOutput } from "@oclif/core/interfaces";

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
    }),
    // These are meant for maintainers
    branch: Flags.string({
      hidden: true,
      char: "b",
      description: "Install runtime from branch on GitHub",
      required: false
    }),
    link: Flags.string({
      hidden: true,
      char: "l",
      description: "Link an in-development runtime to CLI"
    })
  };

  async run(): Promise<void> {
    const ctx = await this.parse(SDKInstallCommand);
    if (ctx.flags.branch) this.installBranch(ctx);
    else if (ctx.flags.link) this.linkDir(ctx);
  }

  async installBranch(ctx: ParserCtx) {
    const { args, flags } = ctx;

    this.log("[1/4] Getting latest commit");
    const { id, branch } = await get_branch_info("hypermodeinc", "modus", "main");
    const version = branch + "/" + id;
    clearLine();
    this.log("[1/4] Found latest commit");

    this.log("[2/4] Fetching Modus from latest commit" + " " + chalk.dim("(" + version + ")"));
    const downloadLink = "https://github.com/hypermodeinc/modus/archive/refs/heads/" + branch + ".zip";
    const archiveName = ("modus-" + version + ".zip").replaceAll("/", "-");
    const tempDir = expandHomeDir("~/.modus/.modus-temp");
    await downloadFile(downloadLink, archiveName);
    clearLine();
    this.log("[2/4] Fetched Modus");

    this.log("[3/4] Unpacking archive");
    mkdirSync(tempDir, { recursive: true });
    const unpackedDir = tempDir + "/" + archiveName.replace(".zip", "");
    await rm(unpackedDir, { recursive: true, force: true });
    if (os.platform() === "win32") {
      execSync("tar -xf " + archiveName + " -C " + unpackedDir);
    } else {
      execSync("unzip " + archiveName + " -d " + unpackedDir);
    }
    clearLine();
    this.log("[3/4] Unpacked archive");

    rmSync(archiveName);
    this.log("[4/4] Installing");
    const installDir = expandHomeDir("~/.modus/sdk/dev-" + branch + "-" + id) + "/";
    cpSync(unpackedDir + "/modus-" + branch + "/", installDir, { recursive: true, force: true });
    clearLine();
    this.log("[4/4] Successfully installed Modus " + version)
  }

  async linkDir(ctx: ParserCtx) {
    const { args, flags } = ctx;

    const srcDir = path.join(process.cwd(), flags.link!);
    if (!existsSync(srcDir)) {
      this.logError("Cannot link to invalid directory! Please try again");
      process.exit(0);
    }

    const installDir = expandHomeDir("~/.modus/sdk/link") + "/";
    this.log("[1/4] Linking directories " + srcDir + " <-> " + installDir);

    const rl = createInterface({
      input: process.stdin,
      output: process.stdout,
    });

    if (!(await this.confirmAction(rl, "[2/4] Continue? [y/n]"))) {
      process.exit(0);
    }

    linkDirectories(srcDir, installDir)
    clearLine();
    clearLine();

    this.log("\nSuccessfully linked directories!");
    process.exit(0);
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }

  private async confirmAction(rl: ReturnType<typeof createInterface>, message: string): Promise<boolean> {
    this.log(message);
    const cont = ((await ask(chalk.dim(" -> "), rl)) || "n").toLowerCase().trim();
    clearLine();
    return cont === "yes" || cont === "y";
  }
}

async function get_latest_release(): Promise<string> {
  const res = (await (await fetch("https://api.github.com/repos/hypermodeinc/modus/releases/latest")).json());
  if (!res) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not find latest release! Please check your internet connection and try again."));
    process.exit(0);
  }
  return res;
}

async function get_branch_info(owner: string, repo: string, branch: string): Promise<{
  sha: string,
  id: string,
  branch: string
}> {
  const res = (await (await fetch(`https://api.github.com/repos/${owner}/${repo}/branches/${branch}`)).json());
  if (!res) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not find branch! Please check your internet connection and try again."));
    process.exit(0);
  }
  return {
    sha: res["commit"]["sha"],
    id: res["commit"]["sha"].slice(0, 7),
    branch: res["name"]
  }
}

function linkDirectories(srcDir: string, destDir: string) {

  rmSync(destDir, { recursive: true, force: true });
  mkdirSync(destDir, { recursive: true });

  const items = readdirSync(srcDir);

  for (const item of items) {
    const srcItem = path.join(srcDir, item);
    const destItem = path.join(destDir, item);

    if (statSync(srcItem).isDirectory()) {
      if (!srcItem.endsWith(".git")) {
        symlinkSync(srcItem, destItem, "dir");
        linkDirectories(srcItem, destItem);
      }
    } else {
      symlinkSync(srcItem, destItem);
    }
  }
}

type ParserCtx = ParserOutput<{
  silent: boolean;
  branch: string | undefined;
  link: string | undefined;
}, {
  [flag: string]: any;
}, {
  version: string | undefined;
}>
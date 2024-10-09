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
      char: "s",
      description: "Suppress output logs",
      hidden: false,
      required: false
    }),
    // These are meant for developers working on modus itself
    // You can link a directory as your runtime
    // You can also run a specific branch
    // You can also run a specific branch and commit on that branch
    repo: Flags.string({
      hidden: true,
      char: "r",
      description: "Change the GitHub owner and repo eg. hypermodeinc/modus",
      required: false,
      default: "hypermodeinc/modus"
    }),
    branch: Flags.string({
      hidden: true,
      char: "b",
      description: "Install runtime from branch on GitHub",
      required: false
    }),
    commit: Flags.string({
      hidden: true,
      char: "c",
      description: "Install runtime from specific commit on GitHub",
      default: "latest",
      required: false
    }),
    link: Flags.string({
      hidden: true,
      char: "l",
      description: "Link an in-development runtime to CLI",
      required: false
    })
  };

  async run(): Promise<void> {
    const ctx = await this.parse(SDKInstallCommand);
    if (ctx.args.version) await this.installVersion(ctx);
    else if (ctx.flags.branch) await this.installCommit(ctx);
    else if (ctx.flags.link) await this.linkDir(ctx);
  }

  async installVersion(ctx: ParserCtx) {
    const { args, flags } = ctx;
    let version = args.version?.toLowerCase();
    const release_info = await fetchRelease(flags.repo, version!);
    if (version == "latest") {
      this.log("[1/4] Getting latest version");
      version = release_info["tag_name"];
      clearLine();
      this.log("[1/4] Latest version: " + chalk.dim(version));
    } else {
      this.log("[1/4] Getting release info")
    }

    version = version?.replace("v", "");

    this.log("[2/4] Downloading release");
    const extension = os.platform() == "win32" ? ".zip" : ".tar.gz";
    const filename = "v" + version + extension;
    const downloadLink = "https://github.com/" + flags.repo + "/archive/refs/tags/" + filename;

    const archiveName = ("modus-" + version + extension).replaceAll("/", "-");
    const tempDir = expandHomeDir("~/.modus/.modus-temp");
    const archivePath = path.join(tempDir, archiveName);

    mkdirSync(tempDir, { recursive: true });
    await downloadFile(downloadLink, archivePath);

    clearLine();
    this.log("[2/4] Downloaded release");

    const unpackedDir = tempDir + "/" + archiveName.replace(extension, "");
    this.log("[3/4] Unpacking release" + "tar -xf " + archivePath + " -C " + unpackedDir);
    await rm(unpackedDir, { recursive: true, force: true });
    mkdirSync(unpackedDir, { recursive: true })
    execSync("tar -xf " + archivePath + " -C " + unpackedDir);

    clearLine();
    this.log("[3/4] Unpacked release");

    this.log("[4/4] Building");
    const installDir = expandHomeDir("~/.modus/sdk/" + version) + "/";
    cpSync(unpackedDir + "/modus-" + version?.replace("v", "") + "/", installDir, { recursive: true, force: true });
    clearLine();

    await rm(tempDir, { recursive: true, force: true });
    this.log("[4/4] Successfully installed Modus " + version)
  }

  async installCommit(ctx: ParserCtx) {
    const { flags } = ctx;

    flags.commit = flags.commit.slice(0, 7);

    let sha = "";

    const commit_info = await fetchCommit(flags.repo, flags.commit);
    sha = commit_info.sha;
    this.log("[1/4] Getting commit " + flags.commit);
    let { id, branch } = await fetchBranch(flags.repo, "main");
    if (sha) id = sha.slice(0, 7);
    const version = branch + "/" + id;
    clearLine();
    this.log("[1/4] Found latest commit");

    this.log("[2/4] Fetching Modus from latest commit" + " " + chalk.dim("(" + version + ")"));
    const downloadLink = flags.commit == "latest" ? "https://github.com/" + flags.repo + "/archive/refs/heads/" + branch + ".zip" : "https://github.com/" + flags.repo + "/archive/" + sha + ".zip";
    const archiveName = ("modus-" + version + ".zip").replaceAll("/", "-");
    const tempDir = expandHomeDir("~/.modus/.modus-temp");
    const archivePath = path.join(tempDir, archiveName);
    mkdirSync(tempDir, { recursive: true });

    await downloadFile(downloadLink, archivePath);
    clearLine();
    this.log("[2/4] Fetched Modus");

    this.log("[3/4] Unpacking archive");
    const unpackedDir = tempDir + "/" + archiveName.replace(".zip", "");
    await rm(unpackedDir, { recursive: true, force: true });
    mkdirSync(unpackedDir, { recursive: true });
    if (os.platform() === "win32") {
      execSync("tar -xf " + archivePath + " -C " + unpackedDir);
    } else {
      execSync("unzip " + archivePath + " -d " + unpackedDir);
    }
    clearLine();
    this.log("[3/4] Unpacked archive");

    this.log("[4/4] Installing");
    const installDir = expandHomeDir("~/.modus/sdk/dev-" + branch + "-" + id) + "/";
    cpSync(unpackedDir + "/modus-" + (flags.commit == "latest" ? branch : sha) + "/", installDir, { recursive: true, force: true });
    clearLine();
    await rm(tempDir, { recursive: true, force: true });
    this.log("[4/4] Successfully installed Modus " + version)
  }

  async linkDir(ctx: ParserCtx) {
    const { flags } = ctx;

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

async function fetchRelease(repo: string, version: string): Promise<{
  tag_name: string
}> {
  const res = (await (await fetch("https://api.github.com/repos/" + repo + "/releases/" + version)).json());
  if (!res) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not find latest release! Please check your internet connection and try again."));
    process.exit(0);
  }
  return res;
}

async function fetchBranch(repo: string, branch: string): Promise<{
  sha: string,
  id: string,
  branch: string
}> {
  const res = (await (await fetch(`https://api.github.com/repos/${repo}/branches/${branch}`)).json());
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

async function fetchCommit(repo: string, commitSha: string): Promise<{
  sha: string
}> {
  const response = await fetch(`https://api.github.com/repos/${repo}/commits/${commitSha}`);
  if (!response.ok) {
    console.log(chalk.red(" ERROR ") + chalk.dim(": Could not find commit! Please check your internet connection and try again."));
    process.exit(0);
  }
  return response.json();
}

function linkDirectories(srcDir: string, destDir: string) {

  rmSync(destDir, { recursive: true, force: true });
  mkdirSync(destDir, { recursive: true });

  const items = readdirSync(srcDir);

  for (const item of items) {
    const srcItem = path.join(srcDir, item);
    const destItem = path.join(destDir, item);

    if (statSync(srcItem).isDirectory()) {
      symlinkSync(srcItem, destItem, "dir");
      linkDirectories(srcItem, destItem);
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
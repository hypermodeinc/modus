/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Args, Flags } from "@oclif/core";
import { spawn } from "node:child_process";
import { Readable } from "node:stream";
import path from "node:path";
import os from "node:os";
import chalk from "chalk";
import pm from "picomatch";
import chokidar from "chokidar";

import * as fs from "../../util/fs.js";
import * as vi from "../../util/versioninfo.js";
import * as installer from "../../util/installer.js";
import { SDK } from "../../custom/globals.js";
import { getHeader } from "../../custom/header.js";
import { getAppInfo } from "../../util/appinfo.js";
import { isOnline, withSpinner } from "../../util/index.js";
import { readHypermodeSettings } from "../../util/hypermode.js";
import BuildCommand from "../build/index.js";
import SDKInstallCommand from "../sdk/install/index.js";
import { BaseCommand } from "../../baseCommand.js";

const MANIFEST_FILE = "modus.json";
const ENV_FILES = [".env", ".env.local", ".env.development", ".env.dev", ".env.development.local", ".env.dev.local"];

export default class DevCommand extends BaseCommand {
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
    "no-build": Flags.boolean({
      aliases: ["nobuild"],
      description: "Don't build the app before running (or when watching for changes)",
    }),
    "no-watch": Flags.boolean({
      aliases: ["nowatch"],
      description: "Don't watch app code for changes",
    }),
    delay: Flags.integer({
      description: "Delay (in milliseconds) between file change detection and rebuild",
      default: 500,
    }),
  };

  static description = "Run a Modus app locally for development";

  static examples = ["modus dev", "modus dev ./my-app", "modus dev ./my-app --no-watch"];

  async run(): Promise<void> {
    const { args, flags } = await this.parse(DevCommand);

    const appPath = path.resolve(args.path);
    if (!(await fs.exists(path.join(appPath, MANIFEST_FILE)))) {
      this.log(chalk.red(`A ${MANIFEST_FILE} file was not found at ${appPath}`));
      this.log(chalk.red("Please either execute the modus command from the app directory, or specify the path to the app you want to run."));
      this.exit(1);
    }

    const app = await getAppInfo(appPath);
    const { sdk, sdkVersion } = app;
    const prerelease = vi.isPrerelease(sdkVersion) || flags.prerelease;

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    if (!(await vi.sdkVersionIsInstalled(sdk, sdkVersion))) {
      await SDKInstallCommand.run([sdk, sdkVersion, "--no-logo"]);
    }

    let runtimeVersion = flags.runtime;
    if (runtimeVersion) {
      const runtimeText = `Modus Runtime ${runtimeVersion}`;
      if (!(await vi.runtimeVersionIsInstalled(runtimeVersion))) {
        if (await isOnline()) {
          await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
            try {
              await installer.installRuntime(runtimeVersion!);
            } catch (e) {
              spinner.fail(chalk.red("Failed to download " + runtimeText));
              throw e;
            }
            spinner.succeed(chalk.dim("Installed " + runtimeText));
          });
        } else {
          this.logError(`${runtimeText} is not installed, and you are offline. Please try again when you have an internet connection.`);
          this.exit(1);
        }
      }
    } else if (await isOnline()) {
      const version = await vi.findLatestCompatibleRuntimeVersion(sdk, sdkVersion, prerelease);
      if (version && !(await vi.runtimeVersionIsInstalled(version))) {
        const runtimeText = `Modus Runtime ${version}`;
        await withSpinner(chalk.dim("Downloading and installing " + runtimeText), async (spinner) => {
          try {
            await installer.installRuntime(version!);
          } catch (e) {
            spinner.fail(chalk.red("Failed to download " + runtimeText));
            throw e;
          }
          spinner.succeed(chalk.dim("Installed " + runtimeText));
        });
      }
      if (!version) {
        this.logError("Could not find a compatible Modus runtime version. Please try again.");
        return;
      }
      runtimeVersion = version;
    } else {
      const version = await vi.findCompatibleInstalledRuntimeVersion(sdk, sdkVersion, prerelease);
      if (!version) {
        this.logError("Could not find a compatible Modus runtime version. Please try again when you have an internet connection.");
        return;
      }
      runtimeVersion = version;
    }

    const ext = os.platform() === "win32" ? ".exe" : "";
    const runtimePath = path.join(vi.getRuntimePath(runtimeVersion), "modus_runtime" + ext);

    // Build the app on first run (this includes copying the manifest file)
    if (!flags["no-build"]) {
      await BuildCommand.run([appPath, "--no-logo"]);
    }

    const appBuildPath = path.join(appPath, "build");

    // Copy env files to the build directory
    await copyEnvFiles(appPath, appBuildPath);

    // Read Hypermode settings if they exist, so they can be forwarded to the runtime
    const hypSettings = await readHypermodeSettings();

    const env = {
      ...process.env,
      MODUS_ENV: "dev",
      HYP_EMAIL: hypSettings.email,
      HYP_API_KEY: hypSettings.apiKey,
      HYP_WORKSPACE_ID: hypSettings.workspaceId,
    };

    // Spawn the runtime child process
    const child = spawn(runtimePath, ["-appPath", appBuildPath, "-refresh=1s"], {
      stdio: ["inherit", "inherit", "pipe"],
      env: env,
    });
    child.stderr.pipe(process.stderr);

    // Handle the runtime process exit
    child.on("close", (code) => {
      // note: can't use "this.exit" here because it would throw an unhandled exception
      // but "process.exit" works fine.
      if (code) {
        this.log(chalk.magentaBright(`Runtime terminated with code ${code}`) + "\n");
        process.exit(code);
      } else {
        this.log(chalk.magentaBright("Runtime terminated successfully.") + "\n");
        process.exit();
      }
    });

    // Forward SIGINT and SIGTERM to the child process for graceful shutdown from user ctrl+c or kill.
    process.on("SIGINT", () => {
      if (child && !child.killed) {
        child.kill("SIGINT");
      }
    });
    process.on("SIGTERM", () => {
      if (child && !child.killed) {
        child.kill("SIGTERM");
      }
    });

    // Watch for changes in the app directory and rebuild the app when changes are detected
    if (!flags["no-watch"]) {
      this.watchForEnvFileChanges(appPath, child.stderr);
      this.watchForManifestChanges(appPath, child.stderr);

      if (!flags["no-build"]) {
        this.watchForAppCodeChanges(appPath, sdk, child.stderr, flags.delay);
      }
    }
  }

  private watchForEnvFileChanges(appPath: string, runtimeOutput: Readable) {
    // Whenever any of the env files change, copy to the build directory.
    // The runtime will automatically reload them when it detects a change to the copies in the build folder.

    const onAddOrChange = async (sourcePath: string) => {
      const filename = path.basename(sourcePath);
      const outputPath = path.join(appPath, "build", filename);
      try {
        runtimeOutput.pause();
        this.log();
        this.log(chalk.magentaBright(`Detected change in ${filename} file. Applying...`));
        await fs.copyFile(sourcePath, outputPath);
      } catch (e) {
        this.log(chalk.red(`Failed to copy ${filename} to build directory.`), e);
      } finally {
        this.log();
        runtimeOutput.resume();
      }
    };

    chokidar
      .watch(appPath, {
        ignored: (filePath) => {
          if (filePath === appPath) {
            return false;
          }
          return !ENV_FILES.includes(path.basename(filePath));
        },
        ignoreInitial: true,
        persistent: true,
      })
      .on("add", onAddOrChange)
      .on("change", onAddOrChange)
      .on("unlink", async (sourcePath) => {
        const filename = path.basename(sourcePath);
        const outputPath = path.join(appPath, "build", filename);
        try {
          runtimeOutput.pause();
          this.log();
          this.log(chalk.magentaBright(`Detected ${filename} deleted. Applying...`));
          await fs.unlink(outputPath);
        } catch (e) {
          this.log(chalk.red(`Failed to delete ${filename} from build directory.`), e);
        } finally {
          this.log();
          runtimeOutput.resume();
        }
      });
  }

  private watchForManifestChanges(appPath: string, runtimeOutput: Readable) {
    // Whenever the manifest file changes, copy it to the build directory.
    // The runtime will automatically reload the manifest when it detects a change to the copy in the build folder.

    const sourcePath = path.join(appPath, MANIFEST_FILE);
    const outputPath = path.join(appPath, "build", MANIFEST_FILE);

    const onAddOrChange = async () => {
      try {
        runtimeOutput.pause();
        this.log();
        this.log(chalk.magentaBright("Detected manifest change. Applying..."));
        await fs.copyFile(sourcePath, outputPath);
      } catch (e) {
        this.log(chalk.red(`Failed to copy ${MANIFEST_FILE} to build directory.`), e);
      } finally {
        this.log();
        runtimeOutput.resume();
      }
    };

    chokidar
      .watch(sourcePath, {
        ignoreInitial: true,
        persistent: true,
      })
      .on("add", onAddOrChange)
      .on("change", onAddOrChange)
      .on("unlink", async () => {
        try {
          runtimeOutput.pause();
          this.log();
          this.log(chalk.magentaBright("Detected manifest deleted. Applying..."));
          await fs.unlink(outputPath);
        } catch (e) {
          this.log(chalk.red(`Failed to delete ${MANIFEST_FILE} from build directory.`), e);
        } finally {
          this.log();
          runtimeOutput.resume();
        }
      });
  }

  private watchForAppCodeChanges(appPath: string, sdk: SDK, runtimeOutput: Readable, delay: number) {
    let lastModified = 0;
    let lastBuild = 0;
    let paused = true;
    setInterval(async () => {
      if (paused) {
        return;
      }
      paused = true;

      if (lastBuild > lastModified) {
        return;
      }
      lastBuild = Date.now();

      try {
        runtimeOutput.pause();
        this.log();
        this.log(chalk.magentaBright("Detected source code change. Rebuilding..."));
        this.log();
        await BuildCommand.run([appPath, "--no-logo"]);
      } catch {
        this.log(chalk.magenta("Waiting for more changes..."));
        this.log(chalk.dim("Press Ctrl+C at any time to stop the server."));
      } finally {
        runtimeOutput.resume();
      }
    }, delay);

    const globs = getGlobsToWatch(sdk);

    // NOTE: The built-in fs.watch or fsPromises.watch is insufficient for our needs.
    // Instead, we use chokidar for consistent behavior in cross-platform file watching.
    const pmOpts: pm.PicomatchOptions = { posixSlashes: true };
    chokidar
      .watch(appPath, {
        ignored: (filePath, stats) => {
          const relativePath = path.relative(appPath, filePath);
          if (!stats || !relativePath) return false;

          let ignore = false;
          if (pm(globs.excluded, pmOpts)(relativePath)) {
            ignore = true;
          } else if (stats.isFile()) {
            ignore = !pm(globs.included, pmOpts)(relativePath);
          }

          if (process.env.MODUS_DEBUG) {
            this.log(chalk.dim(`${ignore ? "ignored: " : "watching:"}  ${relativePath}`));
          }
          return ignore;
        },
        ignoreInitial: true,
        persistent: true,
      })
      .on("all", () => {
        lastModified = Date.now();
        paused = false;
      });
  }

  private logError(message: string) {
    this.log(chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

function getGlobsToWatch(sdk: SDK) {
  const included: string[] = [];
  const excluded: string[] = [".git/**", "build/**"];

  switch (sdk) {
    case SDK.AssemblyScript:
      included.push("**/*.ts", "**/asconfig.json", "**/tsconfig.json", "**/package.json");
      excluded.push("node_modules/**");
      break;

    case SDK.Go:
      included.push("**/*.go", "**/go.mod");
      excluded.push("**/*_generated.go", "**/*.generated.go", "**/*_test.go");
      break;

    default:
      throw new Error(`Unsupported SDK: ${sdk}`);
  }
  return { included, excluded };
}

async function copyEnvFiles(appPath: string, buildPath: string): Promise<void> {
  for (const file of ENV_FILES) {
    const src = path.join(appPath, file);
    const dest = path.join(buildPath, file);
    if (await fs.exists(src)) {
      await fs.copyFile(src, dest);
    } else if (await fs.exists(dest)) {
      await fs.unlink(dest);
    }
  }
}

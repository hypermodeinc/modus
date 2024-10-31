import { Command, Flags } from "@oclif/core";
import chalk from "chalk";
import os from "node:os";
import * as vi from "../../util/versioninfo.js";
import { SDK } from "../../custom/globals.js";
import { getHeader } from "../../custom/header.js";
import { getGoVersion, getNPMVersion, getTinyGoVersion } from "../../util/systemVersions.js";
import { BaseCommand } from "../baseCommand.js";

export default class InfoCommand extends BaseCommand {
  static args = {};

  static flags = {};

  static description = "Show Modus info";

  static examples = ["modus info"];

  async run(): Promise<void> {
    const { flags } = await this.parse(InfoCommand);

    if (!flags["no-logo"]) {
      this.log(getHeader(this.config.version));
    }

    const info = await this.logInfo();
    this.displayInfo(info);
  }

  async logInfo(): Promise<Record<string, string | string[]>> {
    await this.logRuntimeVersions();
    await this.logSDKVersions();

    return {
      "Modus Installation Path": this.config.root,
      "Operating System": `${os.type()} ${os.release()}`,
      "Platform Architecture": process.arch,
      "Node.js Version": process.version,
      "NPM Version": (await getNPMVersion()) || "Not installed",
      "Go Version": (await getGoVersion()) || "Not installed",
      "TinyGo Version": (await getTinyGoVersion()) || "Not installed",
    };
  }

  displayInfo(info: Record<string, string | string[]>): void {
    for (const [key, value] of Object.entries(info)) {
      if (Array.isArray(value)) {
        this.log(`${chalk.bold(key)}:`);
        value.forEach((v) => this.log(`  ${v}`));
      } else if (value) {
        this.log(`${chalk.bold(key)}: ${value}`);
      }
    }
  }

  async logRuntimeVersions() {
    const versions = await vi.getInstalledRuntimeVersions();
    if (versions.length > 0) {
      this.log(chalk.bold.cyan("Installed Runtimes:"));
      for (const version of versions) {
        this.log(`• Modus Runtime ${version}`);
      }
    } else {
      this.log(chalk.yellow("No Modus runtimes installed"));
    }

    this.log();
  }

  async logSDKVersions() {
    let found = false;
    for (const sdk of Object.values(SDK)) {
      const versions = await vi.getInstalledSdkVersions(sdk);
      if (versions.length === 0) {
        continue;
      }
      if (!found) {
        this.log(chalk.bold.cyan("Installed SDKs:"));
        found = true;
      }

      for (const version of versions) {
        this.log(`• Modus ${sdk} SDK ${version}`);
      }
    }

    if (!found) {
      this.log(chalk.yellow("No Modus SDKs installed"));
    }

    this.log();
  }
}

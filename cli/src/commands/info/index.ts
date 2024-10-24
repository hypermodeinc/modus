import { Command, Flags } from "@oclif/core";
import chalk from "chalk";
import os from "node:os";
import * as vi from "../../util/versioninfo.js";
import { SDK } from "../../custom/globals.js";
import { getHeader } from "../../custom/header.js";
import { execFileWithExitCode } from "../../util/cp.js";

export default class InfoCommand extends Command {
  static args = {};

  static flags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    nologo: Flags.boolean({
      aliases: ["no-logo"],
      hidden: true,
    }),
  };

  static description = "Show Modus info";

  static examples = ["modus info"];

  async run(): Promise<void> {
    const { flags } = await this.parse(InfoCommand);

    if (!flags.nologo) {
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
      "NPM Version": await this.getNPMVersion(),
      "Go Version": await this.getGoVersion(),
      "TinyGo Version": await this.getTinyGoVersion(),
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

  async getNPMVersion(): Promise<string> {
    const { stdout } = await execFileWithExitCode("npm", ["--version"], { shell: true });
    return stdout.trim();
  }

  async getGoVersion(): Promise<string> {
    try {
      const { stdout } = await execFileWithExitCode("go", ["version"], { shell: true });
      return stdout.trim().split(" ")[2];
    } catch (error) {
      return "Not installed";
    }
  }

  async getTinyGoVersion(): Promise<string> {
    try {
      const { stdout } = await execFileWithExitCode("tinygo", ["version"], { shell: true });
      return stdout.trim().split(" ")[2];
    } catch (error) {
      return "Not installed";
    }
  }
}

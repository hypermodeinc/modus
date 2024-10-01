import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import { cpSync } from "node:fs";
import os from "node:os";
import path from "node:path";
import { expandHomeDir } from "../../../util/index.js";
import { Metadata } from "../../../util/metadata.js";

const versions = ["0.12.0", "0.12.1","0.12.2","0.12.3","0.12.4","0.12.5","0.12.6"];
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

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKInstallCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk install <version>");
    let version = args.version?.trim().toLowerCase().replace("v", "");
    const platform = os.platform();
    const arch = os.arch();
    const file = "modus-runtime-v" + version + "-" + platform + "-" + arch + (platform === "win32" ? ".exe" : "");
    
    if (version === "all") {
      for (const version of versions) {
        cpSync(path.join(path.dirname(import.meta.url.replace("file:", "")), "../../../../runtime-bin/" + "modus-runtime-v" + version + "-" + platform + "-" + arch + (platform === "win32" ? ".exe" : "")), expandHomeDir("~/.hypermode/sdk/" + version + "/runtime" + (platform === "win32" ? ".exe" : "")));
      }
      this.log("Installed versions 0.12.0-0.12.6");
      return;
    } else if (version === "latest") {
      version = (await Metadata.getLatestRuntime())!;
    }

    cpSync(path.join(path.dirname(import.meta.url.replace("file:", "")), "../../../../runtime-bin/" + file), expandHomeDir("~/.hypermode/sdk/" + version + "/runtime" + (platform === "win32" ? ".exe" : "")));
    this.log("Installed Modus v" + version);
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}
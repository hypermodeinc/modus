import { Args, Command } from "@oclif/core";
import chalk from "chalk";

export default class SDKInstallCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to install",
      hidden: false,
      required: false,
    }),
  };

  static description = "Install a specific SDK version";

  static examples = ["modus install v0.0.0", "modus install latest"];

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKInstallCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk install <version>")
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

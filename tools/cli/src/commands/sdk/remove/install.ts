import { Args, Command } from "@oclif/core";
import chalk from "chalk";
import { readdirSync } from "node:fs";
import { expandHomeDir } from "../../../util/index.js";

export default class SDKRemoveCommand extends Command {
  static args = {
    version: Args.string({
      description: "v0.0.0-|-SDK version to remove",
      hidden: false,
      required: false,
    }),
  };

  static description = "Uninstall a specific SDK version";

  static examples = ["modus remove v0.0.0", "modus remove all"];

  static flags = {};

  async run(): Promise<void> {
    const { args } = await this.parse(SDKRemoveCommand);
    if (!args.version) this.logError("No version specified! Run modus sdk remove <version>"), process.exit(0);
    const version = args.version.trim().toLowerCase();
    
    if (version === "all") {

    }

    const versions = readdirSync(expandHomeDir("~/.hypermode/sdk"));
    for (const v of versions) {
      
    }
  }

  private logError(message: string) {
    this.log("\n" + chalk.red(" ERROR ") + chalk.dim(": " + message));
  }
}

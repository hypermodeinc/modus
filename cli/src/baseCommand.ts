import { Command, Flags } from "@oclif/core";
import { checkForUpdates } from "./util/updateNotifier.js";

export abstract class BaseCommand extends Command {
  static baseFlags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
      helpGroup: "Global",
    }),
    "no-logo": Flags.boolean({
      aliases: ["nologo"],
      hidden: true,
    }),
  };

  async init(): Promise<void> {
    await super.init();

    const cmd = this.id?.split(" ")[0];
    if (cmd !== "version" && !this.argv.includes("--version") && !this.argv.includes("-v")) {
      await checkForUpdates(this.config.version);
    }
  }
}

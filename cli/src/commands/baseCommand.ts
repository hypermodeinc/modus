import { Command, Flags } from "@oclif/core";

export abstract class BaseCommand extends Command {
  static baseFlags = {
    help: Flags.help({
      char: "h",
      helpLabel: "-h, --help",
      description: "Show help message",
    }),
    "no-logo": Flags.boolean({
      aliases: ["nologo"],
      hidden: true,
    }),
  };
}

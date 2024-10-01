import { Command, Help, Interfaces } from "@oclif/core";
import chalk from "chalk";

import { CLI_VERSION } from "./globals.js";

export default class CustomHelp extends Help {
  formatRoot(): string {
    let out = "";
    out += chalk.bold.blueBright("Modus") + " Framework CLI " + chalk.dim("(v" + CLI_VERSION + ")") + "\n\n";

    // Usage: modus <command> [...flags] [...args]
    out += chalk.bold("Usage: modus") + " " + chalk.dim("<command>") + " " + chalk.bold.blueBright("[...flags]") + " " + chalk.bold("[...args]");
    return out;
  }

  formatCommand(command: Command.Loadable): string {
    let out = "";
    out += chalk.bold("Usage:") + " " + chalk.bold.blueBright("modus") + " " + command.id;
    return out;
  }

  formatCommands(commands: Command.Loadable[]): string {
    let out = "";
    out += chalk.bold("Commands:") + "\n";
    for (const command of commands) {
      if (command.id === "autocomplete") continue;

      const rawName = command.id.includes(":") ? command.id.split(":")[1] : command.id;
      const name = chalk.bold.blueBright(rawName);
      const prePadding = " ".repeat(Math.max(1, 10 - rawName.length));
      const args =
        Object.keys(command.args).length > 0
          ? Object.entries(command.args)
            .map((v) => {
              if (!v[1].hidden && v[1].required) {
                if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                  return v[1].description.split("-|-")[0];
                }

                return v[0];
              }

              return "";
            })
            .join(" ")
          : "";
      const postPadding = " ".repeat(Math.max(20 - args.length, 1));
      const description = command.description!;
      const aliases = command.aliases.length > 0 ? chalk.dim(" (" + command.aliases.join("/") + ")") : "";

      out += "  " + name + prePadding + chalk.dim(args) + postPadding + description + aliases + "\n";
    }
    return out.trim();
  }

  formatTopic(topic: Interfaces.Topic): string {
    let out = "";
    out += chalk.bold.blueBright("Modus") + " Help " + chalk.dim("(v0.0.0)") + "\n\n";
    if (topic.description) out += chalk.dim(topic.description) + "\n";

    out += chalk.bold("Usage: modus " + topic.name) + " " + chalk.bold.blue("[command]") + "\n";
    return out;
  }

  formatTopics(topics: Interfaces.Topic[]): string {
    const padding = 30;
    let out = "";
    if (topics.find((v) => !v.hidden)) out += chalk.bold("Tools:") + "\n";
    else return out;

    for (const topic of topics) {
      if (topic.hidden) continue;
      out += "  " + chalk.bold.blue(topic.name) + " ".repeat(Math.max(1, padding - topic.name.length)) + topic.description;
    }
    return out.trim();
  }

  formatRootFooter(): string {
    let out = "";
    out += "View the docs:                  " + chalk.blueBright("https://docs.hypermode.com/introduction") + "\n";
    out += "View the repo:                  " + chalk.blueBright("https://github.com/HypermodeAI/modus") + "\n";

    out += "\n";
    out += "Made with 💖 by " + chalk.magentaBright("https://hypermode.ai");
    return out;
  }

  async showRootHelp(): Promise<void> {
    let rootTopics = this.sortedTopics;
    let rootCommands = this.sortedCommands;
    const state = this.config.pjson?.oclif?.state;
    if (state) {
      this.log(state === "deprecated" ? `${this.config.bin} is deprecated` : `${this.config.bin} is in ${state}.\n`);
    }
    this.log(this.formatRoot());
    this.log("");
    if (!this.opts.all) {
      rootTopics = rootTopics.filter((t) => !t.name.includes(":"));
      rootCommands = rootCommands.filter((c) => !c.id.includes(":"));
    }
    if (rootTopics.length > 0) {
      this.log(this.formatTopics(rootTopics));
      this.log("");
    }
    if (rootCommands.length > 0) {
      rootCommands = rootCommands.filter((c) => c.id);
      this.log(this.formatCommands(rootCommands));
      this.log("");
    }
    this.log(this.formatRootFooter());
  }

  async showTopicHelp(topic: Interfaces.Topic) {
    const { name } = topic;
    const commands = this.sortedCommands.filter((c) => c.id.startsWith(name + ":"));
    const state = this.config.pjson?.oclif?.state;
    if (state) this.log(`This topic is in ${state}.\n`);
    this.log(this.formatTopic(topic));
    if (commands.length > 0) {
      this.log(this.formatCommands(commands));
      this.log("");
    }
  }

  async showCommandHelp(command: Command.Loadable): Promise<void> {
    const name = command.id.replaceAll(":", " ");

    const args =
      Object.keys(command.args).length > 0
        ? Object.entries(command.args)
          .map((v) => {
            if (!v[1].hidden && v[1].required) {
              if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                return v[1].description.split("-|-")[0];
              }
              return v[0];
            }
            return "";
          })
          .join(" ")
        : "";
    const flags = Object.keys(command.flags);

    this.log(chalk.bold.blueBright("Modus") + " Help " + chalk.dim("(v0.0.0)") + "\n");

    if (command.description) this.log(chalk.dim(command.description));

    this.log(chalk.bold("Usage:") + " " + chalk.bold("modus " + name) + (args.length > 0 ? " [...args]" : "") + (flags.length > 0 ? chalk.blueBright(" [...flags]") : "") + "\n");
    // if (examples) {
    //     this.log();
    //     this.log(chalk.bold("Examples:") + "\n");
    //     for (const example of examples) this.log("  " + chalk.dim(example));
    // }

    if (flags.length) {
      const flag_margin = 16;
      this.log(chalk.bold("Flags:"));
      for (const flag of Object.values(command.flags)) this.log("  " + chalk.bold.blueBright("--" + flag.name) + " ".repeat(flag_margin - flag.name.length) + flag.description);
    }

    if (args) {
      const flag_margin = 16;
      this.log(chalk.bold("Args:"));
      for (const flag of Object.values(command.args)) this.log("  " + chalk.bold.blueBright(flag.name) + " ".repeat(flag_margin - flag.name.length) + flag.description);
    }
  }
}

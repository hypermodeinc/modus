/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command, Help, Interfaces } from "@oclif/core";
import { getHeader } from "./header.js";
import chalk from "chalk";

export default class CustomHelp extends Help {
  private target_pad = 15;
  private pre_pad = 0;
  private post_pad = 0;

  formatHeader(): string {
    return getHeader(this.config.version);
  }

  formatRoot(): string {
    return `${chalk.bold("Usage:")} modus ${"<command or tool> [flags] [args]"}`;
  }

  formatCommands(commands: Command.Loadable[]): string {
    let out = "";
    out += chalk.bold("Commands:") + "\n";

    for (const command of commands) {
      if (command.id === "autocomplete") continue;

      const rawName = command.id.includes(":") ? command.id.split(":")[1] : command.id;
      const name = chalk.bold.blueBright(rawName);
      const prePadding = " ".repeat(Math.max(1, this.pre_pad - rawName.length));
      const args =
        Object.keys(command.args).length > 0
          ? Object.entries(command.args)
              .map((v) => {
                if (!v[1].hidden && v[1].required) {
                  if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                    return v[1].description.split("-|-")[0];
                  }
                }
                return "";
              })
              .join(" ")
          : "";
      const postPadding = " ".repeat(Math.max(1, this.post_pad - args.length));
      const description = command.description!;
      const aliases = command.aliases.length > 0 ? chalk.dim(" (" + command.aliases.join("/") + ")") : "";

      out += "  " + name + prePadding + chalk.dim(args) + postPadding + description + aliases + "\n";
    }
    return out.trim();
  }

  formatTopic(topic: Interfaces.Topic): string {
    let out = "";
    if (topic.description) out += chalk.hex("#A585FF")(topic.description) + "\n";
    out += `${chalk.bold("Usage:")} modus ${topic.name} ${chalk.blueBright("<command> [...flags] [...args]")}\n`;
    return out;
  }

  formatTopics(topics: Interfaces.Topic[]): string {
    let out = "";
    if (topics.find((v) => !v.hidden)) out += chalk.bold("Tools:") + "\n";
    else return out;

    for (const topic of topics) {
      if (topic.hidden) continue;
      out += "  " + chalk.bold.blue(topic.name) + " ".repeat(Math.max(1, this.pre_pad + this.post_pad - topic.name.length)) + topic.description + "\n";
    }
    return out.trim();
  }

  formatFlags(topics: Interfaces.Topic[]): string {
    let out = "";
    if (topics.find((v) => !v.hidden)) out += chalk.bold("Flags:") + "\n";
    else return out;

    for (const topic of topics) {
      if (topic.hidden) continue;
      out += "  " + chalk.bold.blue(topic.name) + " ".repeat(Math.max(1, this.pre_pad + this.post_pad - topic.name.length)) + topic.description + "\n";
    }
    return out.trim();
  }

  formatFooter(): string {
    let out = "";
    const links = [
      { name: "Docs", url: "https://docs.hypermode.com/modus" },
      { name: "GitHub", url: "https://github.com/hypermodeinc/modus" },
      { name: "Discord", url: "https://discord.hypermode.com" },
    ];

    for (const link of links) {
      out += `${link.name}: ${" ".repeat(Math.max(1, this.pre_pad + this.post_pad - link.name.length))}${link.url}\n`;
    }

    out += "\n";
    out += "Made with ♥︎ by Hypermode";
    out += "\n";

    return out;
  }

  async showRootHelp(): Promise<void> {
    this.log(this.formatHeader());

    let rootTopics = this.sortedTopics;
    let rootCommands = this.sortedCommands;
    const state = this.config.pjson?.oclif?.state;
    if (state) {
      this.log(state === "deprecated" ? `${this.config.bin} is deprecated` : `${this.config.bin} is in ${state}.\n`);
    }

    this.log(this.formatRoot());
    this.log();
    if (!this.opts.all) {
      rootTopics = rootTopics.filter((t) => !t.name.includes(":"));
      rootCommands = rootCommands.filter((c) => !c.id.includes(":"));
    }

    for (const command of rootCommands) {
      if (command.id.length > this.pre_pad) this.pre_pad = command.id.length;
      const args =
        Object.keys(command.args).length > 0
          ? Object.entries(command.args)
              .map((v) => {
                if (!v[1].hidden && v[1].required) {
                  if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                    return v[1].description.split("-|-")[0];
                  }
                }
                return "";
              })
              .join(" ")
          : "";
      if (args.length > this.post_pad) this.post_pad = args.length;
    }
    this.post_pad = 6 + this.pre_pad + this.post_pad > this.target_pad ? 6 + this.pre_pad + this.post_pad - this.target_pad : this.target_pad - this.pre_pad;
    this.pre_pad += 2;

    if (rootCommands.length > 0) {
      rootCommands = rootCommands.filter((c) => c.id);
      this.log(this.formatCommands(rootCommands));
      this.log();
    }

    if (rootTopics.length > 0) {
      this.log(this.formatTopics(rootTopics));
      this.log();
    }

    const globalFlagTopics: Interfaces.Topic[] = [
      {
        name: "--help",
        description: "Show help message",
      },
      {
        name: "--version",
        description: "Show Modus version",
      },
    ];
    this.log(this.formatFlags(globalFlagTopics));
    this.log();

    this.log(this.formatFooter());
  }

  async showTopicHelp(topic: Interfaces.Topic) {
    this.log(this.formatHeader());

    const { name } = topic;
    const commands = this.sortedCommands.filter((c) => c.id.startsWith(name + ":"));
    for (const command of commands) {
      if (command.id.split(":")[1].length > this.pre_pad) this.pre_pad = command.id.split(":")[1].length;
      const args =
        Object.keys(command.args).length > 0
          ? Object.entries(command.args)
              .map((v) => {
                if (!v[1].hidden && v[1].required) {
                  if (v[1].description && v[1].description.indexOf("-|-") > 0) {
                    return v[1].description.split("-|-")[0];
                  }
                }
                return "";
              })
              .join(" ")
          : "";
      if (args.length > this.post_pad) this.post_pad = args.length;
    }
    this.post_pad = 6 + this.pre_pad + this.post_pad > this.target_pad ? 6 + this.pre_pad + this.post_pad - this.target_pad : this.target_pad - this.pre_pad;
    this.pre_pad += 2;
    const state = this.config.pjson?.oclif?.state;
    if (state) this.log(`This topic is in ${state}.\n`);
    this.log(this.formatTopic(topic));
    if (commands.length > 0) {
      this.log(this.formatCommands(commands));
      this.log();
    }

    this.log(this.formatFooter());
  }

  async showCommandHelp(command: Command.Loadable): Promise<void> {
    this.log(this.formatHeader());

    const margin = 20;
    const name = command.id.replaceAll(":", " ");
    const args = Object.keys(command.args);
    const flags = Object.keys(command.flags);

    if (command.description) this.log(chalk.hex("#A585FF")(command.description));

    const argsStr = args.length > 0 ? " " + args.map((a) => `<${a}>`).join(" ") : "";
    const flagsStr = flags.length > 0 ? " [flags]" : "";
    this.log(`${chalk.bold("Usage:")} modus ${name}${chalk.blueBright(argsStr + flagsStr)}\n`);

    if (args.length) {
      this.log(chalk.bold("Args:"));
      for (let arg of Object.values(command.args)) {
        let usage = "";
        let desc = arg.description;
        if (arg.description?.includes("-|-")) {
          usage = arg.description.split("-|-")[0];
          desc = arg.description.split("-|-")[1];
        }
        let defaultValue = "";
        if (arg.default) {
          defaultValue = chalk.dim(` (default: '${arg.default}')`);
        }
        this.log("  " + chalk.bold.blueBright(arg.name) + " ".repeat(Math.max(1, margin + 2 - arg.name.length)) + desc + defaultValue);
      }
      this.log();
    }

    if (flags.length) {
      this.log(chalk.bold("Flags:"));
      for (const flag of Object.values(command.flags)) {
        if (flag.hidden) continue;

        const flagOptions = flag.char ? `-${flag.char}, --${flag.name}` : `--${flag.name}`;
        this.log("  " + chalk.bold.blueBright(flagOptions) + " ".repeat(margin + 2 - flagOptions.length) + flag.description);
      }
      this.log();
    }

    this.log(this.formatFooter());
  }
}

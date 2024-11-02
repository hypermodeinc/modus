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

const FIRST_PAD = 2;
const SECOND_PAD = 20;

export default class CustomHelp extends Help {
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
      const description = command.description!;
      const aliases = command.aliases.length > 0 ? chalk.dim(" (" + command.aliases.join("/") + ")") : "";

      out += getMessageWithPad(name, description + aliases, FIRST_PAD, SECOND_PAD, rawName.length) + "\n";
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

      out += getMessageWithPad(chalk.bold.blue(topic.name), topic.description || "", FIRST_PAD, SECOND_PAD, topic.name.length) + "\n";
    }
    return out.trim();
  }

  formatFooter(): string {
    let out = "";

    const links = [
      { name: "Docs", url: "https://docs.hypermode.com/modus" },
      { name: "GitHub", url: "https://github.com/hypermodeinc/modus" },
      { name: "Discord", url: "https://discord.hypermode.com (#modus)" },
    ];

    for (const link of links) {
      out += getMessageWithPad(link.name, link.url, 0, SECOND_PAD) + "\n";
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

    if (rootCommands.length > 0) {
      rootCommands = rootCommands.filter((c) => c.id);
      this.log(this.formatCommands(rootCommands));
      this.log();
    }

    if (rootTopics.length > 0) {
      this.log(this.formatTopics(rootTopics));
      this.log();
    }

    this.log();

    this.log(this.formatFooter());
  }

  printRootFlags() {
    const globalFlags = [
      {
        name: "--help",
        shortName: "-h",
        description: "Show help message",
      },
      {
        name: "--version",
        shortName: "-v",
        description: "Show Modus version",
      },
    ];

    let out = "";
    out += chalk.bold("Flags:") + "\n";

    for (const flag of globalFlags) {
      const fullName = flag.shortName ? `${flag.shortName}, ${flag.name}` : flag.name;
      out += getMessageWithPad(chalk.bold.blue(fullName), flag.description || "", FIRST_PAD, SECOND_PAD, fullName.length) + "\n";
    }

    this.log(out.trim());
  }

  async showTopicHelp(topic: Interfaces.Topic) {
    this.log(this.formatHeader());

    const { name } = topic;
    const commands = this.sortedCommands.filter((c) => c.id.startsWith(name + ":"));

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
        let desc = arg.description || "";
        if (arg.description?.includes("-|-")) {
          usage = arg.description.split("-|-")[0];
          desc = arg.description.split("-|-")[1];
        }
        if (arg.default) {
          desc += chalk.dim(` (default: ${arg.default})`);
        }

        if (arg.options) {
          desc += chalk.dim(` (options: ${arg.options.join(", ")})`);
        }

        this.log(getMessageWithPad(chalk.bold.blueBright(arg.name), desc, FIRST_PAD, SECOND_PAD, arg.name.length));
      }
      this.log();
    }

    if (flags.length) {
      this.log(chalk.bold("Flags:"));
      for (const flag of Object.values(command.flags)) {
        if (flag.hidden) continue;

        const flagOptions = flag.char ? `-${flag.char}, --${flag.name}` : `--${flag.name}`;

        let desc = flag.description || "";
        // The interface doesn't have `options`, but `options` is a valid property of a flag
        if ("options" in flag && flag.options) {
          desc += chalk.dim(` (${flag.options.join(", ")})`);
        }

        if (flag) this.log(getMessageWithPad(chalk.bold.blueBright(flagOptions), desc, FIRST_PAD, SECOND_PAD, flagOptions.length));
      }
      this.log();
    }

    this.log(this.formatFooter());
  }
}

/**
 * Used for printing help messages like:
 *   build               Build a Modus app
 *   dev                 Run a Modus app locally for development
 *
 * @param firstMessageLength The message can be ANSI escaped and have a length different from the actual string length
 */
function getMessageWithPad(formattedFirstMsg: string, formattedSecondMsg: string, firstPad: number, secondPad: number, firstMessageLength?: number): string {
  const padLeft = " ".repeat(firstPad);
  const actualFirstMsgLength = firstMessageLength ?? formattedFirstMsg.length;
  const padMiddle = " ".repeat(secondPad - firstPad - actualFirstMsgLength);

  return padLeft + formattedFirstMsg + padMiddle + formattedSecondMsg;
}

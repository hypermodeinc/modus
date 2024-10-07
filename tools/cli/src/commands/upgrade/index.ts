/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command, Flags } from "@oclif/core";

export default class Upgrade extends Command {
  static args = {
    // person: Args.string({ description: "Person to say hello to", required: true }),
  };

  static description = "Upgrade a Modus component";

  static examples = [];

  static flags = {};

  async run(): Promise<void> {
    const { args, flags } = await this.parse(Upgrade);

  }
}

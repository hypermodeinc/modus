/*
 * Copyright 2024 Hypermode, Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode, Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import { Command } from "@oclif/core";
import { rm } from "fs/promises";
import path from "path";
import { fileURLToPath } from "url";

const BASE_DIR = path.join(path.dirname(fileURLToPath(import.meta.url)), "../../../");

export default class UninstallCommand extends Command {
    static description = "Uninstall the Modus CLI";

    static examples = [];

    static flags = {};

    async run(): Promise<void> {
        const { args, flags } = await this.parse(UninstallCommand);
        await rm(BASE_DIR, { recursive: true, force: true });
    }
}

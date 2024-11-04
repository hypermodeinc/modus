/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import gradient from "gradient-string"

const logo = `
▗▖  ▗▖ ▗▄▖ ▗▄▄▄ ▗▖ ▗▖ ▗▄▄▖
▐▛▚▞▜▌▐▌ ▐▌▐▌  █▐▌ ▐▌▐▌   
▐▌  ▐▌▐▌ ▐▌▐▌  █▐▌ ▐▌ ▝▀▚▖
▐▌  ▐▌▝▚▄▞▘▐▙▄▄▀▝▚▄▞▘▗▄▄▞▘
`

export function getLogo(): string {
  return gradient(["#E3BFFF", "#602AF8"]).multiline(logo)
}

/*
 * SPDX-FileCopyrightText: © Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import gradient from "gradient-string";

const logo = `
▗▖  ▗▖ ▗▄▖ ▗▄▄▄ ▗▖ ▗▖ ▗▄▄▖
▐▛▚▞▜▌▐▌ ▐▌▐▌  █▐▌ ▐▌▐▌   
▐▌  ▐▌▐▌ ▐▌▐▌  █▐▌ ▐▌ ▝▀▚▖
▐▌  ▐▌▝▚▄▞▘▐▙▄▄▀▝▚▄▞▘▗▄▄▞▘
`;

export function getLogo(): string {
  return gradient(["#E3BFFF", "#602AF8"]).multiline(logo);
}

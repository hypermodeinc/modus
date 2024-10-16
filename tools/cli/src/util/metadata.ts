/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as globals from "../custom/globals.js";

export class Metadata {
  // CLI version (may want to have it read package.json instead)
  static cli_version: string = globals.CLI_VERSION;

  // Installed runtimes - sorted by date
  static runtimes: string[] = [];

  static initialize(): void {
    // search hypermode.json / package.json for wanted runtime version
    // search current runtimes. populate
  }

  static async getLatestRuntimeVersion(): Promise<string | null> {
    try {
      const tag = await findLatestReleaseTag(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, globals.usePrereleaseVersions);
      if (tag) {
        return tag.slice(globals.GitHubRuntimeTagPrefix.length);
      }
    } catch (e) {
      console.error(e);
    }

    return null;
  }
}

async function findLatestReleaseTag(owner: string, repo: string, prefix: string, prerelease: boolean): Promise<string | null> {
  let page = 1;
  while (true) {
    const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases?page=${page}`, {
      headers: {
        Accept: "application/vnd.github.v3+json",
        "X-GitHub-Api-Version": "2022-11-28",
      },
    });

    if (!response.ok) {
      throw new Error(`Error fetching releases: ${response.statusText}`);
    }

    const releases = await response.json();
    if (releases.length === 0) {
      break;
    }

    for (const release of releases) {
      if (!prerelease && release.prerelease) {
        continue;
      }

      if (prefix && !release.tag_name.startsWith(prefix)) {
        continue;
      }

      return release.tag_name;
    }

    page++;
  }

  return null;
}

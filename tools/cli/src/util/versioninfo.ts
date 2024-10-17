/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import * as globals from "../custom/globals.js";

export async function getLatestRuntimeVersion(prerelease: boolean): Promise<string | undefined> {
  try {
    const tag = await findLatestReleaseTag(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, prerelease);
    if (tag) {
      return tag.slice(globals.GitHubRuntimeTagPrefix.length);
    }
  } catch (e) {
    console.error(e);
  }
}

export async function runtimeReleaseExists(version: string): Promise<boolean> {
  return releaseExists(globals.GitHubOwner, globals.GitHubRepo, `${globals.GitHubRuntimeTagPrefix}${version}`);
}

async function releaseExists(owner: string, repo: string, tag: string): Promise<boolean> {
  const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases/tags/${encodeURIComponent(tag)}`, {
    headers: {
      Accept: "application/vnd.github.v3+json",
      "X-GitHub-Api-Version": "2022-11-28",
    },
  });

  return response.ok;
}

async function findLatestReleaseTag(owner: string, repo: string, prefix: string, prerelease: boolean): Promise<string | undefined> {
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
      return;
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
}

/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import semver from "semver";
import path from "node:path";
import * as fs from "./fs.js";
import * as globals from "../custom/globals.js";

export async function getLatestRuntimeVersion(prerelease: boolean): Promise<string | undefined> {
  try {
    let tag = await findLatestReleaseTag(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, prerelease);
    if (!tag && !prerelease) {
      // If no stable release was found, look for a prerelease
      tag = await findLatestReleaseTag(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, true);
    }
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

export async function findCompatibleSdkVersion(sdk: globals.SDK, runtimeVersion: string): Promise<string | undefined> {
  const versionParts = runtimeVersion.split(".");
  const versionPrefix = versionParts.slice(0, 2).join(".") + ".";
  const prerelease = versionParts.length > 2 && versionParts[2].includes("-");

  const sdkPrefix = globals.GetSdkTagPrefix(sdk);
  const searchPrefix = `${sdkPrefix}${versionPrefix}`; // ex: "sdk/assemblyscript/v0.13."
  const tag = await findLatestReleaseTag(globals.GitHubOwner, globals.GitHubRepo, searchPrefix, prerelease);
  if (tag) {
    return tag.slice(sdkPrefix.length);
  }
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

export async function getLatestInstalledVersion(): Promise<string | undefined> {
  const dir = path.join(globals.ModusHomeDir, "sdk");
  if (!(await fs.exists(dir))) {
    return;
  }

  const subdirs = (await fs.readdir(dir, { withFileTypes: true })).filter((e) => e.isDirectory()).map((e) => e.name);
  return semver.sort(subdirs).pop();
}

export async function getLatestTemplatesArchive(mainVersion: string, sdk: string): Promise<{ version: string; path: string } | undefined> {
  const dir = path.join(globals.ModusHomeDir, "sdk", mainVersion);
  if (!(await fs.exists(dir))) {
    return;
  }

  const prefix = `templates_${sdk}_v`;
  const suffix = ".tar.gz";

  const versions = (await fs.readdir(dir)).filter((f) => f.startsWith(prefix) && f.endsWith(suffix)).map((f) => f.slice(prefix.length, -suffix.length));

  semver.sort(versions);
  const latestVersion = versions.pop();
  if (latestVersion) {
    return {
      version: "v" + latestVersion,
      path: path.join(dir, prefix + latestVersion + suffix),
    };
  }
}

export async function versionIsInstalled(version: string): Promise<boolean> {
  const dir = path.join(globals.ModusHomeDir, "sdk", version);
  return await fs.exists(dir);
}

export async function getInstalledVersions(): Promise<string[]> {
  const dir = path.join(globals.ModusHomeDir, "sdk");
  if (!(await fs.exists(dir))) {
    return [];
  }

  const versions = (await fs.readdir(dir, { withFileTypes: true })).filter((e) => e.isDirectory()).map((e) => e.name);
  return semver.rsort(versions);
}

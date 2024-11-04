/*
 * Copyright 2024 Hypermode Inc.
 * Licensed under the terms of the Apache License, Version 2.0
 * See the LICENSE file that accompanied this code for further details.
 *
 * SPDX-FileCopyrightText: 2024 Hypermode Inc. <hello@hypermode.com>
 * SPDX-License-Identifier: Apache-2.0
 */

import semver from "semver";
import * as globals from "./constants.js";

function getGitHubApiHeaders() {
  const headers = new Headers({
    Accept: "application/vnd.github.v3+json",
    "X-GitHub-Api-Version": "2022-11-28",
    "User-Agent": "Modus CLI",
  });
  if (process.env.GITHUB_TOKEN) {
    headers.append("Authorization", `Bearer ${process.env.GITHUB_TOKEN}`);
  }
  return headers;
}

export function isPrerelease(version: string): boolean {
  if (version.startsWith("v")) {
    version = version.slice(1);
  }
  return !!semver.prerelease(version);
}

export async function getLatestSdkVersion(sdk: globals.SDK, includePrerelease: boolean): Promise<string | undefined> {
  return await getLatestVersion(globals.GitHubOwner, globals.GitHubRepo, globals.GetSdkTagPrefix(sdk), includePrerelease);
}

export async function getLatestRuntimeVersion(includePrerelease: boolean): Promise<string | undefined> {
  return await getLatestVersion(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, includePrerelease);
}

export async function getLatestCliVersion(includePrerelease: boolean): Promise<string | undefined> {
  return await getLatestVersion(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubCliTagPrefix, includePrerelease);
}

async function getLatestVersion(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string | undefined> {
  try {
    let tag = await findLatestReleaseTag(owner, repo, prefix, includePrerelease);
    if (!tag && !includePrerelease) {
      // If no stable release was found, look for a prerelease
      tag = await findLatestReleaseTag(owner, repo, prefix, true);
    }
    if (tag) {
      return tag.slice(prefix.length);
    }
  } catch (e) {
    console.error(e);
  }
}

export async function getAllSdkVersions(sdk: globals.SDK, includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubOwner, globals.GitHubRepo, globals.GetSdkTagPrefix(sdk), includePrerelease);
}

export async function getAllCliVersions(includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubCliTagPrefix, includePrerelease);
}

export async function getAllRuntimeVersions(includePrerelease: boolean): Promise<string[]> {
  return await getAllVersions(globals.GitHubOwner, globals.GitHubRepo, globals.GitHubRuntimeTagPrefix, includePrerelease);
}

async function getAllVersions(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string[]> {
  try {
    let tags = await getAllReleaseTags(owner, repo, prefix, includePrerelease);
    if (tags.length === 0 && !includePrerelease) {
      // If no stable release was found, look for prereleases
      tags = await getAllReleaseTags(owner, repo, prefix, true);
    }
    let versions = tags.map((tag) => {
      let version = tag.slice(prefix.length);
      if (version.startsWith("v")) {
        version = version.slice(1);
      }
      return version;
    });
    versions = semver.rsort(versions);
    return versions.map((v) => "v" + v);
  } catch (e) {
    console.error(e);
    return [];
  }
}

const headers = getGitHubApiHeaders();

async function findLatestReleaseTag(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string | undefined> {
  let page = 1;
  while (true) {
    const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases?page=${page}`, { headers });

    if (!response.ok) {
      throw new Error(`Error fetching releases: ${response.statusText}`);
    }

    const releases = await response.json();
    if (releases.length === 0) {
      return;
    }

    for (const release of releases) {
      if (!includePrerelease && release.prerelease) {
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

async function getAllReleaseTags(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string[]> {
  const results: string[] = [];

  let page = 1;
  while (true) {
    const response = await fetch(`https://api.github.com/repos/${owner}/${repo}/releases?per_page=100&page=${page}`, { headers });

    if (!response.ok) {
      throw new Error(`Error fetching releases: ${response.statusText}`);
    }

    const releases = await response.json();
    if (releases.length === 0) {
      return results;
    }

    for (const release of releases) {
      if (!includePrerelease && release.prerelease) {
        continue;
      }

      if (prefix && !release.tag_name.startsWith(prefix)) {
        continue;
      }

      results.push(release.tag_name);
    }

    page++;
  }
}

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
import { execSync } from "child_process";

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
    throw new Error(`Error fetching latest release tag of ${owner}/${repo}: ${e}`);
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
    throw new Error(`Error fetching release tags of ${owner}/${repo}: ${e}`);
  }
}

function execGitCommand(command: string): string {
  try {
    return execSync(command, {
      encoding: "utf8",
      timeout: 10000,
    }).trim();
  } catch (error) {
    console.error(`Error executing git command: ${command}`, error);
    throw error;
  }
}

async function findLatestReleaseTag(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string | undefined> {
  const repoUrl = `https://github.com/${owner}/${repo}`;
  const normalizedPrefix = prefix.endsWith("/") ? prefix : prefix;
  const pattern = `refs/tags/${normalizedPrefix}*`;

  const command = `git ls-remote --refs ${repoUrl} ${pattern}`;
  const output = execGitCommand(command);

  if (!output) {
    return undefined;
  }

  // Format of each line: <commit-hash>\trefs/tags/<tag-name>
  const tags = output
    .split("\n")
    .map((line) => {
      const parts = line.split("\t");
      if (parts.length !== 2) return null;

      const tagName = parts[1].replace("refs/tags/", "");

      if (!tagName.startsWith(normalizedPrefix)) {
        return null;
      }

      return tagName;
    })
    .filter((tag) => tag !== null) as string[];

  if (tags.length === 0) {
    return undefined;
  }

  let filteredTags = tags;
  if (!includePrerelease) {
    const stableVersions = tags.filter((tag) => {
      const version = tag.slice(normalizedPrefix.length);
      return !isPrerelease(version);
    });

    if (stableVersions.length > 0) {
      filteredTags = stableVersions;
    }
  }

  const versions = filteredTags.map((tag) => {
    const version = tag.slice(normalizedPrefix.length);
    return version.startsWith("v") ? version.slice(1) : version;
  });

  if (versions.length === 0) {
    return undefined;
  }

  const latestVersion = semver.rsort(versions)[0];
  return normalizedPrefix + (latestVersion.startsWith("v") ? latestVersion : "v" + latestVersion);
}

async function getAllReleaseTags(owner: string, repo: string, prefix: string, includePrerelease: boolean): Promise<string[]> {
  const repoUrl = `https://github.com/${owner}/${repo}`;
  const normalizedPrefix = prefix.endsWith("/") ? prefix : prefix;
  const pattern = `refs/tags/${normalizedPrefix}*`;

  const command = `git ls-remote --refs ${repoUrl} ${pattern}`;
  const output = execGitCommand(command);

  if (!output) {
    return [];
  }

  // Format of each line: <commit-hash>\trefs/tags/<tag-name>
  const tags = output
    .split("\n")
    .map((line) => {
      const parts = line.split("\t");
      if (parts.length !== 2) return null;

      // Extract tag name from refs/tags/prefix/vX.Y.Z
      const tagName = parts[1].replace("refs/tags/", "");

      // Skip if it doesn't start with the prefix
      if (!tagName.startsWith(normalizedPrefix)) {
        return null;
      }

      return tagName;
    })
    .filter((tag) => tag !== null) as string[];

  if (!includePrerelease) {
    const stableVersions = tags.filter((tag) => {
      const version = tag.slice(normalizedPrefix.length);
      return !isPrerelease(version);
    });

    if (stableVersions.length > 0) {
      return stableVersions;
    }
  }

  return tags;
}

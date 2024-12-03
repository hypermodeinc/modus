/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

// These classes are used by the example functions in the index.ts file.

@json
export class Person {
  constructor(
    uid: string = "",
    firstName: string = "",
    lastName: string = "",
    dType: string[] = [],
  ) {
    this.uid = uid;
    this.firstName = firstName;
    this.lastName = lastName;
    this.dType = dType;
  }
  uid: string = "";
  firstName: string = "";
  lastName: string = "";


  @alias("dgraph.type")
  dType: string[] = [];
}


@json
export class PeopleData {
  people!: Person[];
}


@json
export class Plugin {
  constructor(
    uid: string = "",
    id: string = "",
    name: string = "",
    version: string = "",
    language: string = "",
    sdkVersion: string = "",
    buildId: string = "",
    buildTime: string = "",
    gitRepo: string = "",
    gitCommit: string = "",
    dType: string[] = [],
  ) {
    this.uid = uid;
    this.id = id;
    this.name = name;
    this.version = version;
    this.language = language;
    this.sdkVersion = sdkVersion;
    this.buildId = buildId;
    this.buildTime = buildTime;
    this.gitRepo = gitRepo;
    this.gitCommit = gitCommit;
    this.dType = dType;
  }
  uid: string = "";
  id: string = "";
  name: string = "";
  version: string = "";
  language: string = "";


  @alias("sdk_version")
  sdkVersion: string = "";


  @alias("build_id")
  buildId: string = "";


  @alias("build_time")
  buildTime: string = "";


  @alias("git_repo")
  gitRepo: string = "";


  @alias("git_commit")
  gitCommit: string = "";


  @alias("dgraph.type")
  dType: string[] = [];
}


@json
export class PluginData {
  plugins!: Plugin[];
}

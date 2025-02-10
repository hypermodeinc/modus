/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

// These classes are used by the example functions in the index.ts file.

@json
export class Quote {

  @alias("q")
  quote!: string;


  @alias("a")
  author!: string;
}

export class Image {
  contentType!: string;
  data!: Uint8Array;
}


@json
export class Issue {
  title!: string;
  body!: string;

  // The URL of the issue on GitHub, after the issue is created.
  @alias("html_url")
  url: string | null = null;
}

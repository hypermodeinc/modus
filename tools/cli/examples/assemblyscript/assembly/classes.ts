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
  data!: ArrayBuffer;
}


@json
export class Issue {
  title!: string;
  body!: string;

  // The URL of the issue on GitHub, after the issue is created.
  @alias("html_url")
  url: string | null = null;
}

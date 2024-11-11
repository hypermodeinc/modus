/*
 * This example is part of the Modus project, licensed under the Apache License 2.0.
 * You may modify and use this example in accordance with the license.
 * See the LICENSE file that accompanied this code for further details.
 */

import { http } from "@hypermode/modus-sdk-as";
import { Quote, Image, Issue } from "./classes";

// This function makes a simple HTTP GET request to example.com,
// and returns the HTML text of the response.
export function getExampleHtml(): string {
  const response = http.fetch("https://example.com/");
  return response.text();
}

// This function makes a request to an API that returns data in JSON format,
// and returns an object representing the data.
// It also demonstrates how to check the HTTP response status.
export function getRandomQuote(): Quote {
  const request = new http.Request("https://zenquotes.io/api/random");

  const response = http.fetch(request);
  if (!response.ok) {
    throw new Error(
      `Failed to fetch quote. Received: ${response.status} ${response.statusText}`,
    );
  }

  // The API returns an array of quotes, but we only want the first one.
  return response.json<Quote[]>()[0];
}

// This function makes a request to an API that returns an image, and returns the image data.
export function getRandomImage(width: i32, height: i32): Image {
  const url = `https://picsum.photos/${width}/${height}`;
  const response = http.fetch(url);
  const contentType = response.headers.get("Content-Type")!;

  return {
    contentType,
    data: response.body,
  };
}

/*
This function demonstrates a more complex HTTP call.
It makes a POST request to the GitHub API to create an issue.

See https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28#create-an-issue

To use it, you must add a GitHub personal access token to your secrets.
Create a fine-grained token at https://github.com/settings/tokens?type=beta with access
to write issues to the repository you want to use, then add it to the appropriate secret
store for your environment.  (See the modus documentation for details.)

NOTE: Do not pass the Authorization header in code when creating the request.
That would be a security risk, as the token could be exposed in the source code repository.

Instead, configure the headers in the modus.json manifest as follows:

	"hosts": {
	  "github": {
	    "baseUrl": "https://api.github.com/",
	    "headers": {
	      "Authorization": "Bearer {{AUTH_TOKEN}}"
	    }
	  }
	}

The Modus runtime will retrieve the token from your secrets and add it to the request.
*/
export function createGithubIssue(
  owner: string,
  repo: string,
  title: string,
  body: string,
): Issue {
  // The URL for creating an issue in a GitHub repository.
  const url = `https://api.github.com/repos/${owner}/${repo}/issues`;

  // Create a new request with the URL, method, and headers.
  const request = new http.Request(url, {
    method: "POST",
    headers: http.Headers.from([
      // Do not pass an Authorization header here. See note above.
      ["Accept", "application/vnd.github+json"],
      ["X-GitHub-Api-Version", "2022-11-28"],
      ["Content-Type", "application/json"],
    ]),

    // The request body will be sent as JSON from the Issue object passed here.
    body: http.Content.from(<Issue>{ title, body }),
  } as http.RequestOptions);

  // Send the request and check the response status.
  // NOTE: If you are using a private repository, and you get a 404 error, that could
  // be an authentication issue. Make sure you have created a token as described above.
  const response = http.fetch(request);
  if (!response.ok) {
    throw new Error(
      `Failed to create issue. Received: ${response.status} ${response.statusText}`,
    );
  }

  // The response will contain the issue data, including the URL of the issue on GitHub.
  return response.json<Issue>();
}
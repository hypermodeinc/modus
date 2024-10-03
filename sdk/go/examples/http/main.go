package main

import (
	"errors"
	"fmt"

	"github.com/hypermodeAI/functions-go/pkg/http"
)

// This function makes a simple HTTP GET request to example.com, and returns the HTML text of the response.
func GetExampleHtml() (string, error) {
	response, err := http.Fetch("https://example.com/")
	if err != nil {
		return "", err
	}
	return response.Text(), nil
}

// This function makes a request to an API that returns data in JSON format, and returns an object representing the data.
// It also demonstrates how to check the HTTP response status.
func GetRandomQuote() (*Quote, error) {
	request := http.NewRequest("https://zenquotes.io/api/random")

	response, err := http.Fetch(request)
	if err != nil {
		return nil, err
	}
	if !response.Ok() {
		return nil, fmt.Errorf("Failed to fetch quote. Received: %d %s", response.Status, response.StatusText)
	}

	// The API returns an array of quotes, but we only want the first one.
	var quotes []Quote
	response.JSON(&quotes)
	if len(quotes) == 0 {
		return nil, errors.New("expected at least one quote in the response, but none were found")
	}
	return &quotes[0], nil
}

// This function makes a request to an API that returns an image, and returns the image data.
func GetRandomImage(width, height int) (*Image, error) {
	url := fmt.Sprintf("https://picsum.photos/%d/%d", width, height)
	response, err := http.Fetch(url)
	if err != nil {
		return nil, err
	}

	result := &Image{
		ContentType: *response.Headers.Get("Content-Type"),
		Data:        response.Body,
	}

	return result, nil
}

/*
This function demonstrates a more complex HTTP call.
It makes a POST request to the GitHub API to create an issue.

See https://docs.github.com/en/rest/issues/issues?apiVersion=2022-11-28#create-an-issue

To use it, you must add a GitHub personal access token to your Hypermode secrets.
Create a fine-grained token at https://github.com/settings/tokens?type=beta with access
to write issues to the repository you want to use.

NOTE: Do not pass the Authorization header when creating the request in code.
That would be a security risk, as the token could be exposed in the source code repository.

Instead, configure the headers in the hypermode.json manifest as follows:

	"hosts": {
	  "github": {
	    "baseUrl": "https://api.github.com/",
	    "headers": {
	      "Authorization": "Bearer {{AUTH_TOKEN}}"
	    }
	  }
	}

The Hypermode Runtime will retrieve the token from your secrets and add it to the request.
To set the secret, after committing the changes to your code and manifest file, go to the
Hypermode Console UI, find the "github" host and add your token to the AUTH_TOKEN secret.
*/
func CreateGithubIssue(owner, repo, title, body string) (*Issue, error) {

	// The URL for creating an issue in a GitHub repository.
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)

	// Create a new request with the URL, method, and headers.
	request := http.NewRequest(url, &http.RequestOptions{
		Method: "POST",
		Headers: map[string]string{
			// Do not pass an Authorization header here. See note above.
			"Accept":               "application/vnd.github+json",
			"X-GitHub-Api-Version": "2022-11-28",
			"Content-Type":         "application/json",
		},

		// The request body will be sent as JSON from the Issue object passed here.
		Body: &Issue{
			Title: title,
			Body:  body,
		},
	})

	// Send the request and check the response status.
	// NOTE: If you are using a private repository, and you get a 404 error, that could
	// be an authentication issue. Make sure you have created a token as described above.
	response, err := http.Fetch(request)
	if err != nil {
		return nil, err
	}
	if !response.Ok() {
		return nil, fmt.Errorf("Failed to create issue. Received: %d %s", response.Status, response.StatusText)
	}

	// The response will contain the issue data, including the URL of the issue on GitHub.
	var issue *Issue
	response.JSON(issue)
	return issue, nil
}

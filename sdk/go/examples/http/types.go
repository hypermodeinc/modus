package main

// These types are used by the example functions in the main.go file.

type Quote struct {
	Quote  string `json:"q"`
	Author string `json:"a"`
}

type Image struct {
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
}

type Issue struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	URL   *string `json:"html_url"`
}

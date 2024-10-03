/*
 * Copyright 2024 Hypermode, Inc.
 */

package http

import (
	"encoding/json"
)

type Response struct {
	Status     uint16
	StatusText string
	Headers    *Headers
	Body       []byte
}

func (r *Response) Ok() bool {
	return r.Status >= 200 && r.Status < 300
}

func (r *Response) Text() string {
	return string(r.Body)
}

func (r *Response) JSON(result any) {
	if err := json.Unmarshal(r.Body, result); err != nil {
		panic(err)
	}
}

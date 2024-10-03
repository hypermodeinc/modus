/*
 * Copyright 2024 Hypermode, Inc.
 */

package http

import (
	"strings"
)

type Header struct {
	Name   string
	Values []string
}

type Headers struct {
	data map[string]*Header
}

func NewHeaders[T [][]string | map[string]string | map[string][]string](value T) *Headers {
	h := &Headers{}
	switch t := any(value).(type) {
	case [][]string:
		for _, entry := range t {
			for i := 1; i < len(entry); i++ {
				h.Append(entry[0], entry[i])
			}
		}
	case map[string]string:
		for name, value := range t {
			h.Append(name, value)
		}
	case map[string][]string:
		for name, values := range t {
			for _, value := range values {
				h.Append(name, value)
			}
		}
	}
	return h
}

func (h *Headers) Append(name, value string) {
	key := strings.ToLower(name)
	if h.data == nil {
		h.data = make(map[string]*Header)
		h.data[key] = &Header{name, []string{value}}
	} else if header, ok := h.data[key]; !ok {
		h.data[key] = &Header{name, []string{value}}
	} else {
		header.Values = append(header.Values, value)
		h.data[key] = header
	}
}

func (h *Headers) Entries() [][]string {
	entries := make([][]string, 0, len(h.data))
	for _, header := range h.data {
		for _, value := range header.Values {
			entries = append(entries, []string{header.Name, value})
		}
	}
	return entries
}

func (h *Headers) Get(name string) *string {
	key := strings.ToLower(name)
	if header, ok := h.data[key]; ok {
		result := strings.Join(header.Values, ",")
		return &result
	}

	return nil
}

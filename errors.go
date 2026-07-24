package main

import "fmt"

type HTTPStatusError struct {
	URL    string
	Code   int
	Status string
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("GET %q: unexpected status %s", e.URL, e.Status)
}

type ContentTypeError struct {
	URL         string
	ContentType string
}

func (e *ContentTypeError) Error() string {
	return fmt.Sprintf(
		"GET %q: unsupported content type %q",
		e.URL,
		e.ContentType,
	)
}

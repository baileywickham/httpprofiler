package main

import (
	"fmt"
)

type HTTPRequest struct {
	// An http request header found here
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages
	// Defining this struct allows more flexibility over just a get request,
	// and can be extended to test other request types.
	// Start-line
	HTTPMethod    string
	RequestTarget string
	HTTPVersion   string
	// Headers
	GeneralHeaders string
	RequestHeaders string
	EntityHeaders  string
	// Body
	Body string
}

func (r *HTTPRequest) String() string {
	return fmt.Sprintf("%s %s %s\r\n%s%s%s\r\n\r\n%s",
		r.HTTPMethod, r.RequestTarget, r.HTTPVersion,
		r.GeneralHeaders, r.RequestHeaders, r.EntityHeaders,
		r.Body)
}

func createGetRequest(url string) string {
	r := HTTPRequest{
		HTTPMethod:     "GET",
		RequestTarget:  "/", //url,
		HTTPVersion:    "HTTP/1.0",
		GeneralHeaders: "",
		RequestHeaders: "",
		EntityHeaders:  "",
		Body:           "",
	}
	return r.String()
}

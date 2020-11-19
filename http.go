package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
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

type HTTPResponse struct {
	// An http repsonse found here
	// https://developer.mozilla.org/en-US/docs/Web/HTTP/Messages

	// Start-line
	HTTPVersion string
	StatusCode  string
	StatusText  string
	// Headers
	GeneralHeaders  string
	ResponseHeaders string
	EntityHeaders   string
	// Body
	Body string
	Size int
}

func (r *HTTPRequest) String() string {
	// Add startline to header
	s := fmt.Sprintf("%s %s %s\r\n", r.HTTPMethod, r.RequestTarget, r.HTTPVersion)
	//%s%s%s\r\n\r\n%s,
	// Append headers
	if r.GeneralHeaders != "" {
		s += r.GeneralHeaders
	}
	if r.RequestHeaders != "" {
		s += r.RequestHeaders
	}
	if r.EntityHeaders != "" {
		s += r.EntityHeaders
	}
	// Append break signifying start of body
	s += "\r\n\r\n"
	// Append body
	if r.Body != "" {
		s += r.Body
	}
	return s
}

func createGetRequest(u url.URL, keepAlive bool) string {
	// create a new get request. Eventually can be generalized to
	// have more header fields, but for profiling these are all that are needed.
	r := HTTPRequest{
		HTTPMethod:    "GET",
		RequestTarget: u.RequestURI(),
		// Force 1.0 to avoid chunking
		// HTTP1.1 may be supported in future versions of this code
		HTTPVersion:    "HTTP/1.0",
		GeneralHeaders: fmt.Sprintf("Host: %s\n", u.Hostname()),
		RequestHeaders: "",
		EntityHeaders:  "",
		Body:           "",
	}
	// If we want to use a single tcp connection we add the keep-alive header
	if keepAlive {
		r.GeneralHeaders += "Connection: Keep-Alive"
	}
	return r.String()
}

func readResponse(conn net.Conn) HTTPResponse {
	var contentLen int
	r := bufio.NewReader(conn)
	_startline, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	// Read in the startline and split on spaces.
	startline := strings.Fields(_startline)

	if len(startline) < 3 {
		panic("Malformated startline")
	}

	// Create the http Response struct, fill the startline
	rsp := HTTPResponse{
		HTTPVersion: startline[0],
		StatusCode:  startline[1],
		StatusText:  startline[2],
	}
	//
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if line == "\r\n" || line == "" || line == "\r" {
			break
		}
		if strings.Contains(line, "Content-Length: ") {
			// hack to get content length
			// This is the case where we have a body
			contentLen, _ = strconv.Atoi(strings.Fields(line)[1])
		}
		rsp.GeneralHeaders += line
	}
	body := make([]byte, contentLen)
	_, err = io.ReadFull(r, body)
	// Error reading body, should only happen on TCP problems
	if err != nil {
		panic(err)
	}

	rsp.Size = contentLen
	rsp.Body = string(body)
	return rsp
}

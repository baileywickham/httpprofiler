package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"net/url"
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
	s := fmt.Sprintf("%s %s %s\r\n", r.HTTPMethod, r.RequestTarget, r.HTTPVersion)
	//%s%s%s\r\n\r\n%s,
	if r.GeneralHeaders != "" {
		s += r.GeneralHeaders
	}
	if r.RequestHeaders != "" {
		s += r.RequestHeaders
	}
	if r.EntityHeaders != "" {
		s += r.EntityHeaders
	}
	s += "\r\n\r\n"
	if r.Body != "" {
		s += r.Body
	}
	return s
}

func createGetRequest(u url.URL) string {
	r := HTTPRequest{
		HTTPMethod:     "GET",
		RequestTarget:  u.RequestURI(), //url,
		HTTPVersion:    "HTTP/1.0",
		GeneralHeaders: fmt.Sprintf("Host: %s", u.Hostname()),
		RequestHeaders: "",
		EntityHeaders:  "",
		Body:           "",
	}
	return r.String()
}

func readResponse(conn net.Conn) HTTPResponse {
	r := bufio.NewReader(conn)
	_startline, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	startline := strings.Fields(_startline)

	rsp := HTTPResponse{
		HTTPVersion: startline[0],
		StatusCode:  startline[1],
		StatusText:  startline[2],
	}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			panic(err)
		}
		if line == "\r\n" {
			break
		}
		//if strings.Contains(line, "Content-Length:") {
		//	length := strings.Fields(line)[1]
		//	println(length)
		//}
		rsp.GeneralHeaders += line
	}
	body, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	rsp.Size = len(body)
	rsp.Body = string(body)
	return rsp
}

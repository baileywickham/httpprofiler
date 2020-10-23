package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Profile struct {
	n       int
	rspTime time.Duration
	size    int
	rsp     HTTPResponse
}

var _url = flag.String("url", "http://cloudflare.com", "URL to profile")
var n = flag.Int("profile", 1, "Number of requests to send")
var keepalive = flag.Bool("keepalive", false, "Attempt to use a keepalive request to use the same TCP connection")
var v = flag.Bool("v", false, "Print responses as they are recieved")

func parseURL(_url string) url.URL {
	// Parse the URL into a URL struct
	u, err := url.Parse(_url)
	if err != nil {
		panic(err)
	}
	return *u
}

func startProfileKeepAlive(u url.URL) []Profile {
	// profile with a single TCP request
	profs := make([]Profile, *n)

	// Open a tcp connection with the url requsted. Use http or https
	conn, err := net.Dial("tcp", u.Hostname()+":"+u.Scheme)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	// Create a HTTP get request, true indicates the keepalive header should
	// be included
	request := createGetRequest(u, true)
	for i := 0; i < *n; i++ {
		start := time.Now()
		// Send that request
		fmt.Fprintf(conn, request)
		rsp := readResponse(conn)
		duration := time.Since(start)
		p := Profile{
			n:       i,
			rspTime: duration,
			rsp:     rsp}
		if *v {
			printRsp(rsp, i, duration)
		}
		profs[i] = p

		// Stop on closed connection
		if strings.Contains(rsp.GeneralHeaders+rsp.ResponseHeaders+rsp.EntityHeaders,
			"Connection: close") {
			panic("Connection closed, keepalive response not sent")
		}
	}
	return profs
}

func startProfile(u url.URL) []Profile {
	// startProfileTCP includes the time of the TCP connection,
	// not only the http request
	profs := make([]Profile, *n)

	for i := 0; i < *n; i++ {
		conn, err := net.Dial("tcp", u.Hostname()+":"+u.Scheme)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		// Create a HTTP get request, no keepalive
		request := createGetRequest(u, false)
		start := time.Now()
		// Send that request
		fmt.Fprintf(conn, request)
		rsp := readResponse(conn)
		duration := time.Since(start)
		p := Profile{
			n:       i,
			rspTime: duration,
			rsp:     rsp}
		if *v {
			printRsp(rsp, i, duration)
		}
		profs[i] = p
	}
	return profs
}

func printRsp(rsp HTTPResponse, n int, t time.Duration) {
	fmt.Printf("Request number: %d\nStatus Code: %s\nResponse time: %dms\n",
		n, rsp.StatusCode, t.Milliseconds())
}

func evaluate(profs []Profile) {
	// A success response is a 200 response
	sucesses := 0
	times := make([]time.Duration, n)
	errors := make([]string, 0)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("Number of requests: %d\n", *n)

	for i, p := range profs {
		code, _ := strconv.Atoi(p.rsp.StatusCode)
		if code >= 200 && code < 300 {
			sucesses++
		} else {
			errors = append(errors, p.rsp.StatusCode)
		}
		times[i] = p.rspTime
	}
}

func main() {
	flag.Parse()

	u := parseURL(*_url)

	if *keepalive {
		evaluate(startProfileKeepAlive(u))
	} else {
		evaluate(startProfile(u))
	}
}

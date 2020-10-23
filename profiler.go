package main

import (
	"flag"
	"fmt"
	"net"
	"net/url"
	"time"
)

type Profile struct {
	n       int
	rspTime time.Duration
	size    int
	rsp     HTTPResponse
}

func parseURL(_url string) url.URL {
	// Parse the URL into a URL struct
	u, err := url.Parse(_url)
	if err != nil {
		panic(err)
	}
	return *u
}

func startProfile(u url.URL, n int, v bool) []Profile {
	// profile with a single TCP request
	profs := make([]Profile, n)

	// Open a tcp connection with the url requsted. Use http or https
	conn, err := net.Dial("tcp", u.Hostname()+":"+u.Scheme)
	if err != nil {
		panic(err)
	}

	// Create a HTTP get request
	request := createGetRequest(u)
	for i := 0; i <= n; i++ {
		start := time.Now()
		// Send that request
		fmt.Fprintf(conn, request)
		rsp := readResponse(conn)
		p := Profile{
			n:       i,
			rspTime: time.Since(start),
			rsp:     rsp}
		if v {
			printRsp(rsp, n)
		}
		profs[i] = p
	}
	return profs
}

func startProfileTCP(u url.URL, n int, v bool) []Profile {
	// startProfileTCP includes the time of the TCP connection,
	// not only the http request
	profs := make([]Profile, n)

	for i := 0; i <= n; i++ {
		conn, err := net.Dial("tcp", u.Hostname()+":"+u.Scheme)
		defer conn.Close()
		if err != nil {
			panic(err)
		}
		// Create a HTTP get request
		request := createGetRequest(u)
		start := time.Now()
		// Send that request
		fmt.Fprintf(conn, request)
		rsp := readResponse(conn)
		p := Profile{
			n:       i,
			rspTime: time.Since(start),
			rsp:     rsp}
		if v {
			printRsp(rsp, n)
		}
		profs[i] = p
	}
	return profs
}

func printRsp(rsp HTTPResponse, n int) {

}

func evaluate(profs []Profile, n int) {
	//times := make([]time.Duration, n)
	//for i, p := range profs {
	//	times[i] = p.rspTime
	//}
}

func main() {
	_url := flag.String("url", "http://cloudflare.com", "URL to profile")
	n := flag.Int("profile", 1, "Number of requests to send")
	tcp := flag.Bool("tcp", false, "Start a new tcp connection between requests?")
	v := flag.Bool("v", false, "Print responses as they are recieved")
	flag.Parse()

	u := parseURL(*_url)

	if *tcp {
		evaluate(startProfileTCP(u, *n, *v), *n)
	} else {
		evaluate(startProfile(u, *n, *v), *n)
	}
}

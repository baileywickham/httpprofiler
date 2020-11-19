package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/montanaflynn/stats"
)

type Profile struct {
	n       int
	rspTime time.Duration
	size    int
	rsp     HTTPResponse
}

var _url = flag.String("url", "http://cloudflare.com", "URL to profile")
var n = flag.Int("profile", 0, "Number of profile requests to send, if -1 then print the body and exit")
var keepalive = flag.Bool("keepalive", false, "Attempt to use a keepalive request to use the same TCP connection, fails on Connection: closed response")
var v = flag.Bool("verbose", false, "Print responses as they are recieved")

func parseURL(_url string) url.URL {
	// Parse the URL into a URL struct
	u, err := url.Parse(_url)
	if err != nil {
		panic(err)
	}
	return *u
}

func printBody(u url.URL) {
	fmt.Printf("Fetching %s\n", u.String())
	conn, err := net.Dial("tcp", u.Hostname()+":"+u.Scheme)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	// Create a HTTP get request, true indicates the keepalive header should
	// be included
	request := createGetRequest(u, true)
	fmt.Fprintf(conn, request)
	rsp := readResponse(conn)
	println(rsp.Body)

}

func startProfileKeepAlive(u url.URL) []Profile {
	// profile with a single TCP request
	profs := make([]Profile, *n)

	// Open a tcp connection with the url requsted. Use http
	// https is not currently supported
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
			log.Fatal("Connection closed, keep-alive response not sent")
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

func evaluate(profs []Profile, u url.URL) {
	// A success response is a 200 response
	successes := 0
	times := make([]time.Duration, *n)
	timesFloat := make([]float64, *n)
	sizes := make([]int, *n)
	errCodes := make([]string, 0)

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println("URL: ", u.String())

	for i, p := range profs {
		code, _ := strconv.Atoi(p.rsp.StatusCode)
		if code >= 200 && code < 300 {
			successes++
		} else {
			errCodes = append(errCodes, p.rsp.StatusCode)
		}
		times[i] = p.rspTime
		timesFloat[i] = float64(p.rspTime)
		sizes[i] = p.rsp.Size
	}
	minTime, maxTime := MinMaxDuration(times)
	fmt.Printf("Number of requests: %d\n", *n)
	fmt.Printf("Fastest request: %dms\n", minTime.Milliseconds())
	fmt.Printf("Slowest request: %dms\n", maxTime.Milliseconds())
	mean, _ := stats.Mean(timesFloat)
	fmt.Printf("Mean time: %dms\n", time.Duration(mean).Milliseconds())
	median, _ := stats.Median(timesFloat)
	fmt.Printf("Median time: %dms\n", time.Duration(median).Milliseconds())
	fmt.Printf("Percent successful: %.2f\n", 100*float64(successes)/float64(*n))
	fmt.Printf("Non 2xx error codes: %v\n", errCodes)
	minSize, maxSize := MinMaxInt(sizes)
	fmt.Printf("Smallest response body: %d bytes\n", minSize)
	fmt.Printf("Largest response body: %d bytes\n", maxSize)
}

func main() {
	flag.Parse()

	u := parseURL(*_url)
	if *n == 0 {
		// print body mode
		printBody(u)
		return
	}
	// Use a single TCP request
	if *keepalive {
		evaluate(startProfileKeepAlive(u), u)
	} else {
		evaluate(startProfile(u), u)
	}
}

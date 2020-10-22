package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
)

func startProfile(url string) {
	conn, err := net.Dial("tcp", url+":80")
	if err != nil {
		panic(err)
	}
	request := createGetRequest(url)
	fmt.Fprintf(conn, request)
	//	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		panic(err)
	}
	fmt.Println(status)
}

func echoHelp() {
	println("Profile a url")
}

func main() {
	url := flag.String("url", "localhost", "url to profile")
	help := flag.Bool("help", false, "print this message")
	flag.Parse()
	if *help {
		echoHelp()
		return
	}
	startProfile(*url)
}

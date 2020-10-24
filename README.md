# httpprofiler
Profile websites using custom http requests

Note, this project does not yet support https.
```
Usage of ./httpprofiler:
  -keepalive
        Attempt to use a keepalive connection to use the same TCP connection, fails reguarly
  -profile int
        Number of requests to send (default 2)
  -url string
        URL to profile (default "http://cloudflare.com")
  -verbose
        Print responses as they are recieved
```
![bw.baileywickham.workers.dev](workers.png)

## Use
A Dockerfile is provided which runs a default request

Install the stats package:
```bash
go get github.com/montanaflynn/stats
```

Build and run the script
```golang
go bulid .
./httpprofiler -url http://google.com
```

## Project structure
- `http.go` http specific funcitons and structs
- `profile.go` main program logic, tcp connection, printing
- `helper.go` min, max helper functions, maybe eventually median, mean funcs

## FAQ
### Why did you not use goroutines?
While `conn` is safe across goroutines, using multiple threads/goroutines would mess up the timing of responses. This project also would not benifit from parallelization.

### Why use the stats package?
Dealing with time.Duration can be painful because many of the built in methods like sort don't deal with it. Therefore it's easier to convert Durations and use a library.

### What is the keepalive option?
Passing in -keepalive attempts to use a single tcp connection for all http requests. This can be speedy, but can also cause problems when the server closes the connection on you early. Because my http is hacked together, this option reguarly fails. **Use at your own risk**.

### Why are there no tests?
Testing this program would be a little difficult. I could test that I am getting a response, and that it is well formatted, but this is testing the server as much as my program. I also don't really want to test the output, so I would have to rewrite to make the correct functions exposed which is more work than it's worth for this. 

### Why does the size in the profile only include the body, not the headers?
When looking at a response, I think the size should only really refer to the size of the body because the headers change with the request, or the method of which you are making the request. I don't think it makes sense to have `curl` and `httpprofiler` return different sizes on the same static webpage. 


## Findings
![cloudflare/example.com](cloudflare.png)

Here we test against `cloudflare.com` and `example.com`. Looking at the cloudflare response you will see something that you will see often with this program: failure on a 301 response. This is because most of the modern internet forces redirects from http->https. This is a good thing! However, my program doesn't yet support https, and a 301 response is not a 2xx response, so it is counted as a failure. `example.com` on the otherhand does not redirect to https and returns as you would expect. This matches what you would expect from a profile. 

One thing that would be interesting to add would be a stddev mesaurement, to calculate for variation. 



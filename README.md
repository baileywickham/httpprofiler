# httpprofiler
Profile websites using custom http requests

Note, this project does not yet support https.
```
Usage of ./httpprofiler:
  -keepalive
        Attempt to use a keepalive connection to use the same TCP connection, fails on Connection: closed response
  -profile int
        Number of requests to send (default 2)
  -url string
        URL to profile (default "http://cloudflare.com")
  -verbose
        Print responses as they are recieved
```
![bw.baileywickham.workers.dev](workers.png)

## Use
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



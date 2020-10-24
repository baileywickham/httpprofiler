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




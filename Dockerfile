FROM golang:1.15.0-alpine3.12

WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...

EXPOSE 8080
CMD go run . -url http://bw.baileywickham.workers.dev -verbose -profile 10


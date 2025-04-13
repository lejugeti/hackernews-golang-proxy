# hackernews-golang-proxy

This repository contains a Golang gRPC server and client that proxy [HackerNews API](https://github.com/HackerNews/API).

The project leverages protobuf gRPC generator along with [github.com/peterhellberg/hn](https://github.com/peterhellberg/hn) HackerNews client.

## Golang version

Version used for development : `go version go1.24.2 darwin/arm64`

## gRPC server/client

The server and client are based on the services generated from [this protobuf schema](./grpc_news.proto).

Command to generate:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"

protoc --go_out=./generated \
    --go_opt=paths=source_relative \
    --go-grpc_out=./generated \
    --go-grpc_opt=paths=source_relative \
    grpc_news.proto
```

## Server

The HackerNews proxy server can fetch from HackerNews API :

- First nth top stories
- User details

Results are stored in cache to speed up future requests. Note that each cached data have a time to live, meaning that after a defined period following its addition, it will be evicted automatically from the cache.

### Usage

```bash
go run server/main.go up
```

## Client

Command line tool used to send commands to the HackerNews proxy server.

### Flags

- -list: Fetches top stories
- -max: Indicate the number of stories to fetch (default: 10)
- -timeout: Timeout in seconds before the client cutting connection with the server (default: 20)
- -whois: Fetches user details based on its nickname

Note that either the `-list` or the `-whois` flags **must be used**. These flags however **cannot be used together**.

### Usage

```bash
# fetch user details based on his/her nickname
go run client/main.go -whois fra

# fetch 10 first HackerNews top stories
go run client/main.go -list -max 10

# increase timeout if the request takes too much time
go run client/main.go -list -max 50 -timeout 40
```
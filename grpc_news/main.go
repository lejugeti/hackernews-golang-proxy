package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	grpcHn "hackernews/generated"

	proxyServer "hackernews/grpc_news/server"
	sts "hackernews/grpc_news/stories"
	us "hackernews/grpc_news/users"
)

const port int = 50051

func main() {
    if len(os.Args) == 1 || os.Args[1] != "up" {
        fmt.Println("Usage : go run grpc_news/main.go up")
        return
    }

    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()

	hnServer := proxyServer.NewHnProxyServer(
		sts.NewHackernewsStoriesProxy(),
		us.NewHackernewsUserProxy(),
	)

    grpcHn.RegisterHnServiceServer(s, &hnServer)
    log.Printf("server listening at %v", listener.Addr())
    if err := s.Serve(listener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"google.golang.org/grpc"

	grpcHn "hackernews/generated"

	"hackernews/server/cache"
	proxyServer "hackernews/server/server"
	sts "hackernews/server/stories"
	us "hackernews/server/users"

	hn "github.com/peterhellberg/hn"
)

const port int = 50051
const defaultCacheTtlSeconds uint32 = 40
const defaultClientTimeoutSeconds uint32 = 20

func main() {
    if len(os.Args) == 1 || os.Args[1] != "up" {
        fmt.Println("Usage : go run server/main.go up")
        return
    }

    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()

	timeToLiveDuration := time.Second * time.Duration(defaultCacheTtlSeconds)
	userCache := cache.NewTimeToLiveCache[string, *us.User](timeToLiveDuration)
	storiesCache := cache.NewTimeToLiveCache[int, *sts.Story](timeToLiveDuration)

	clientTimeout := time.Duration(defaultClientTimeoutSeconds) * time.Second
	hnClient := hn.NewClient(&http.Client{Timeout: clientTimeout})

	hnServer := proxyServer.NewHnProxyServer(
		sts.NewHackernewsStoriesProxy(*hnClient, storiesCache),
		us.NewHackernewsUserProxy(*hnClient, userCache),
	)

    grpcHn.RegisterHnServiceServer(s, &hnServer)
    log.Printf("server listening at %v", listener.Addr())
    if err := s.Serve(listener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

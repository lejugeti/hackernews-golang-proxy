package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	hn "hackernews/generated"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
    hn.UnimplementedHnServiceServer
}

func (s *server) GetTopStories(_ context.Context, _ *hn.TopStoriesRequest) (*hn.TopStories, error) {
    log.Print("list news")

    return &hn.TopStories{}, nil
}

func (s *server) Whois(_ context.Context, userRequest *hn.UserInfoRequest) (*hn.User, error) {
    log.Print("whos")

    return &hn.User{Nick: userRequest.GetName()}, nil
}

const port int = 50051

func main() {
    if len(os.Args) == 0 || os.Args[1] != "up" {
        fmt.Println("Usage : go run grpc_news/main.go up")
        return
    }

    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    hn.RegisterHnServiceServer(s, &server{})
    if err := s.Serve(listener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
    log.Printf("server listening at %v", listener.Addr())

}

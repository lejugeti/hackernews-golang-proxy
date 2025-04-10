package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	hn "hackernews/generated"
)
const serverAddress = "localhost:50051"

const listFlag string = "list"
const whoisFlag string = "whois"

var (
    listNewsMode = flag.Bool(listFlag, false, "List the top news from HackerNews front page")
    userName = flag.String(whoisFlag, "", "Retrieve information on user passed as input")
)

func main() {
    flag.Parse()

    if !*listNewsMode && len(*userName) == 0 {
        fmt.Printf("Use at least and only one argument among -%s and -%s\n", listFlag, whoisFlag)
        flag.PrintDefaults()
        return
    } else if *listNewsMode && len(*userName) != 0 {
        fmt.Printf("Arguments -%s and -%s cannot be used together\n", listFlag, whoisFlag)
        flag.PrintDefaults()
        return
    }

    // Set up a connection to the server.
    conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Cannot connect to server: %v", err)
    }
    defer conn.Close()
    c := hn.NewHnServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    if *listNewsMode {
        r, err := c.GetTopStories(ctx, &hn.TopStoriesRequest{})
        if err != nil {
            log.Fatalf("Error: %v", err)
        }
        log.Printf("Top series: %v", r.GetStories())
    } else if len(*userName) != 0 {
        r, err := c.Whois(ctx, &hn.UserInfoRequest{Name: *userName})
        if err != nil {
            log.Fatalf("Error: %v", err)
        }
        log.Printf("User name: %s", r.GetNick())
    }
    
}

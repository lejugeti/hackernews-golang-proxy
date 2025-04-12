package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcHn "hackernews/generated"
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
    c := grpcHn.NewHnServiceClient(conn)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second)
    defer cancel()

    if *listNewsMode {
        r, err := c.GetTopStories(ctx, &grpcHn.TopStoriesRequest{})
        if err != nil {
            log.Fatalf("Error: %v", err)
        }
        log.Printf("Top series: %v", r.GetStories())
    } else if len(*userName) != 0 {
        user, err := c.Whois(ctx, &grpcHn.UserInfoRequest{Name: *userName})
        if err != nil {
            log.Fatalf("Error: %v", err)
        }

		fmt.Printf("User:   %s\n", user.GetNickname())
		fmt.Printf("Karma:  %d\n", user.GetKarma())
		fmt.Printf("About:  %s\n", user.GetAbout())
		fmt.Printf("Joined: %s\n", time.Unix(user.GetJoinedAt(), 0).Format(time.DateOnly))
    }
    
}

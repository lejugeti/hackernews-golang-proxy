package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"

	grpcHn "hackernews/generated"
)

func GetTopStories(client *grpcHn.HnServiceClient, context *context.Context, maxStoriesCount *int) {
	
	if *maxStoriesCount <= 0 {
		fmt.Println("Stories number to fetch must be a positive number")
		return
	}

	request := grpcHn.TopStoriesRequest{StoryNumber: uint32(*maxStoriesCount)}
	topStories, err := (*client).GetTopStories(*context, &request)
    
	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			fmt.Printf("Server took too long to answer the request. You can consider adding more timeout with the -%s flag\n", timeoutFlag)
			return
		} else {
			fmt.Printf("Error: %v\n", err.Error())
			return
		}
	}
	
	for _, story := range topStories.Stories {
		fmt.Printf("- %s\n", story.Title)
		fmt.Printf("  %s\n", story.Url)
		fmt.Println()
	}
}

func GetUserInfo(client *grpcHn.HnServiceClient, context *context.Context, userName *string) {
	if *userName == "" {
		fmt.Println("Please provide a username to fetch user details")
		return
	}

	request := grpcHn.UserInfoRequest{Name: *userName}
	user, err := (*client).Whois(*context, &request)
			
	if status.Code(err) == codes.NotFound {
		fmt.Printf("User '%s' does not exist in HackerNews\n", *userName)
		return
	} else if status.Code(err) == codes.DeadlineExceeded {
		fmt.Printf("Server took too long to answer the request. You can consider adding more timeout with the -%s flag\n", timeoutFlag)
		return
	}else if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return
	}

	fmt.Printf("User:   %s\n", user.GetNickname())
	fmt.Printf("Karma:  %d\n", user.GetKarma())
	fmt.Printf("About:  %s\n", user.GetAbout())
	fmt.Printf("Joined: %s\n", time.Unix(user.GetJoinedAt(), 0).Format(time.DateOnly))
}

const serverAddress = "localhost:50051"

const listFlag string = "list"
const newsNumberFlag string = "max"
const timeoutFlag string = "timeout"
const whoisFlag string = "whois"

var (
    userName = flag.String(whoisFlag, "", "Retrieve information on user passed as input")
    isListMode = flag.Bool(listFlag, false, "Number of top news from HackerNews front page to fetch")
    newsNumber = flag.Int(newsNumberFlag, 10, fmt.Sprintf("Max number of news to fetch. Must be used along with the -%s flag", listFlag))
    timeoutSeconds = flag.Int(timeoutFlag, 20, "Timeout in seconds before client cutting connection to server")
)

func main() {
    flag.Parse()

	var isUserMode bool = *userName != ""

    if !*isListMode && !isUserMode {
        fmt.Printf("Use at least and only one argument among -%s and -%s\n", listFlag, whoisFlag)
        flag.PrintDefaults()
        return
    } else if *isListMode && isUserMode {
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
    
	client := grpcHn.NewHnServiceClient(conn)

	maxTimeToWait := time.Duration(*timeoutSeconds) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), maxTimeToWait)
    defer cancel()

    if *isListMode {
        GetTopStories(&client, &ctx, newsNumber)
    } else if isUserMode {
        GetUserInfo(&client, &ctx, userName)
    } else {
		fmt.Print("Client could not choose any mode to fetch information from HackerNews")
		flag.PrintDefaults()
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcHn "hackernews/generated"

	sts "hackernews/grpc_news/stories"
	us "hackernews/grpc_news/users"
)

type hackernewsServer struct {
    grpcHn.UnimplementedHnServiceServer
	UserService us.UserService
	StoriesService sts.StoriesService
}

func (s *hackernewsServer) GetTopStories(_ context.Context, storiesRequest *grpcHn.TopStoriesRequest) (*grpcHn.TopStories, error) {
	stories, err := s.StoriesService.GetTopStories(storiesRequest.GetStoryNumber())

	if err != nil {
		log.Printf("Error while retrieving top stories. Cause: %s\n", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error while retrieving top stories. Caused by: %s", err.Error())
	}

	var mappedStories = make([]*grpcHn.Story, len(*stories))

	for i, story := range *stories {
		mappedStories[i] = &grpcHn.Story {
			Title: story.Title,
			Url: story.Url,
		}
	}

    return &grpcHn.TopStories{Stories: mappedStories}, nil
}

// Fetches information about a user based on his/her nickname
func (s *hackernewsServer) Whois(_ context.Context, userRequest *grpcHn.UserInfoRequest) (*grpcHn.User, error){
	if userRequest.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "could not fetch user details because no user nickname was provided")
	}

	user, err := s.UserService.GetUserInfo(userRequest.GetName())
	
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user '%s' not found", userRequest.Name)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "could not get user information. Caused by: %s", err.Error())
	}

	return &grpcHn.User{
		Nickname: user.Nickname,
		About: user.About,
		Karma: user.Karma,
		JoinedAt: int64(user.Joined.Unix()),
	}, nil	
}

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
	hnServer := hackernewsServer{
		UserService: us.NewHackernewsUserProxy(),
		StoriesService: sts.NewHackernewsStoriesProxy()}

    grpcHn.RegisterHnServiceServer(s, &hnServer)
    log.Printf("server listening at %v", listener.Addr())
    if err := s.Serve(listener); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}

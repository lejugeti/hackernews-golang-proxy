package server

import (
	"context"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpcHn "hackernews/generated"

	sts "hackernews/grpc_news/stories"
	us "hackernews/grpc_news/users"
)

type hackernewsProxyServer struct {
    grpcHn.UnimplementedHnServiceServer // necessary for grpc to work
	UserService us.UserService
	StoriesService sts.StoriesService
}

func NewHnProxyServer(storiesService sts.StoriesService, userService us.UserService) hackernewsProxyServer {
	return hackernewsProxyServer{
		StoriesService: storiesService,
		UserService: userService,
	}
}

// Fetches first nth top stories and their basic information
func (s *hackernewsProxyServer) GetTopStories(_ context.Context, storiesRequest *grpcHn.TopStoriesRequest) (*grpcHn.TopStories, error) {
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
func (s *hackernewsProxyServer) Whois(_ context.Context, userRequest *grpcHn.UserInfoRequest) (*grpcHn.User, error){
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
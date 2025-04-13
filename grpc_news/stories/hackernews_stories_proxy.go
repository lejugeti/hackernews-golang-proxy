package stories

import (
	"net/http"
	"time"

	hn "github.com/peterhellberg/hn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type hackernewsStoriesProxy struct {}


func (hsp *hackernewsStoriesProxy) GetTopStories(maxStoryCount uint32) (*[]Story, error) {
	hn := hn.NewClient(&http.Client{Timeout: time.Duration(10 * time.Second)})

	idsStories, err := hn.TopStories()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred during top stories fetch. Cause: %v", err)
	}

	var stories = make([]Story, maxStoryCount)

	for i, storyId := range idsStories[:maxStoryCount] {
		rawStory, err := hn.Item(storyId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "could not fetch story '%d'. Cause: %v", storyId, err)
		}	

		stories[i] = Story{
			Id: rawStory.ID,
			Title: rawStory.Title,
			Url: rawStory.URL,
		}
	}

	return &stories, nil
}

func NewHackernewsStoriesProxy() (StoriesService) {
	return &hackernewsStoriesProxy{}
}
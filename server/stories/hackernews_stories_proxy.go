package stories

import (
	"hackernews/server/cache"
	"net/http"
	"time"

	hn "github.com/peterhellberg/hn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type hackernewsStoriesProxy struct {
	cache cache.Cache[int, *Story]
}

func NewHackernewsStoriesProxy(cache cache.Cache[int, *Story]) (StoriesService) {
	return &hackernewsStoriesProxy{
		cache: cache,
	}
}

func (hsp *hackernewsStoriesProxy) GetTopStories(maxStoryCount uint32) (*[]Story, error) {
	hnClient := hn.NewClient(&http.Client{Timeout: time.Duration(10 * time.Second)})

	idsStories, err := hnClient.TopStories()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred during top stories fetch. Cause: %v", err)
	}

	var stories = make([]Story, maxStoryCount)

	for i, storyId := range idsStories[:maxStoryCount] {
		storyFromCache, storyIsCached  := hsp.cache.Get(storyId)

		if storyIsCached {
			stories[i] = *storyFromCache
		} else {
			story, err := hsp.fetchStory(hnClient, storyId)

			if err != nil {
				return nil, status.Errorf(codes.Internal, "error encountered while fetching top stories. Cause: %v", err)
			}

			stories[i] = *story
			hsp.cache.Add(storyId, story)
		}
	}

	return &stories, nil
}

func (hsp *hackernewsStoriesProxy) fetchStory(client *hn.Client, id int) (*Story, error) {
	rawStory, err := client.Item(id)
	
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch story '%d'. Cause: %v", id, err)
	}	

	return &Story{
		Id: rawStory.ID,
		Title: rawStory.Title,
		Url: rawStory.URL,
	}, nil
}

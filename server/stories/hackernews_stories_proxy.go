package stories

import (
	"hackernews/server/cache"

	hn "github.com/peterhellberg/hn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type hackernewsStoriesProxy struct {
	hnClient hn.Client
	cache cache.Cache[int, *Story]
}

func NewHackernewsStoriesProxy(client hn.Client, cache cache.Cache[int, *Story]) (StoriesService) {
	return &hackernewsStoriesProxy{
		hnClient: client,
		cache: cache,
	}
}

func (hsp *hackernewsStoriesProxy) GetTopStories(maxStoryCount uint32) (*[]Story, error) {
	idsStories, err := hsp.hnClient.TopStories()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error occurred during top stories fetch. Cause: %v", err)
	}

	var stories = make([]Story, maxStoryCount)

	for i, storyId := range idsStories[:maxStoryCount] {
		storyFromCache, storyIsCached  := hsp.cache.Get(storyId)

		if storyIsCached {
			stories[i] = *storyFromCache
		} else {
			story, err := hsp.fetchStory(storyId)

			if err != nil {
				return nil, status.Errorf(codes.Internal, "error encountered while fetching top stories. Cause: %v", err)
			}

			stories[i] = *story
			hsp.cache.Add(storyId, story)
		}
	}

	return &stories, nil
}

func (hsp *hackernewsStoriesProxy) fetchStory(id int) (*Story, error) {
	rawStory, err := hsp.hnClient.Item(id)
	
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch story '%d'. Cause: %v", id, err)
	}	

	return &Story{
		Id: rawStory.ID,
		Title: rawStory.Title,
		Url: rawStory.URL,
	}, nil
}

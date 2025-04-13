package stories

import (
	"errors"
	"hackernews/server/cache"
	"testing"

	hn "github.com/peterhellberg/hn"
)

type MockHnLiveService struct {
	MockedTopStories func() ([]int, error)
	MockedMaxItem func() (int, error)
	MockedUpdates func() (*hn.Updates, error)
}

func (m MockHnLiveService) TopStories() ([]int, error) {
	return m.MockedTopStories()
}

func (m MockHnLiveService) MaxItem() (int, error) {
	return m.MockedMaxItem()
}

func (m MockHnLiveService) Updates() (*hn.Updates, error) {
	return m.MockedUpdates()
}

type MockHnItemService struct {
	MockedItem func(id int) (*hn.Item, error)
}

func (m MockHnItemService) Get(id int) (*hn.Item, error) {
	return m.MockedItem(id)
}

func TestGetTopStoriesReturnErrorIfTopStoriesFetchFails(t *testing.T) {
	// GIVEN
	mockLiveService := MockHnLiveService{}
	mockLiveService.MockedTopStories = func() ([]int, error) {
		return nil, errors.New("top stories fetch fail")
	}

	var client hn.Client = hn.Client{Live: mockLiveService}
	var storiesCache cache.Cache[int, *Story] = cache.NewTimeToLiveCache[int, *Story](10)
	var service StoriesService = NewHackernewsStoriesProxy(client, storiesCache)

	// WHEN
	_, err := service.GetTopStories(1)

	// THEN
	if err == nil {
		t.Error("error should be raised if nickname empty")
	}
}

func TestGetTopStoriesShouldGetStoriesFromCache(t *testing.T) {
	// GIVEN
	topStories := []int{0, 1}

	mockLiveService := MockHnLiveService{}
	mockLiveService.MockedTopStories = func() ([]int, error) {
		return topStories, nil
	}

	var client hn.Client = hn.Client{Live: mockLiveService}
	var storiesCache cache.Cache[int, *Story] = cache.NewTimeToLiveCache[int, *Story](10)
	var service StoriesService = NewHackernewsStoriesProxy(client, storiesCache)

	for _, storyId := range topStories {
		storiesCache.Add(storyId, &Story{Id: storyId})
	}

	// WHEN
	stories, err := service.GetTopStories(uint32(len(topStories)))

	// THEN
	if err != nil {
		t.Error("no error should be met")
	} else if len(*stories) != len(topStories) {
		t.Error("should have found 2 stories")
	}

	for i, expectedStoryId := range topStories {
		actualStory := (*stories)[i]
		
		if actualStory.Id != expectedStoryId {
			t.Errorf("Found actual id '%d' but expected '%d'", actualStory.Id, expectedStoryId)
		}
	}
}

func TestGetTopStoriesShouldFetchStories(t *testing.T) {
	// GIVEN
	storyId := 0
	topStoriesIds := []int{storyId}
	topStory := hn.Item{ID: storyId}

	mockLiveService := MockHnLiveService{}
	mockLiveService.MockedTopStories = func() ([]int, error) {
		return topStoriesIds, nil
	}

	mockItemService := MockHnItemService{}
	mockItemService.MockedItem = func(id int) (*hn.Item, error) {
		return &topStory, nil
	}

	var client hn.Client = hn.Client{Live: mockLiveService, Items: mockItemService}
	var storiesCache cache.Cache[int, *Story] = cache.NewTimeToLiveCache[int, *Story](10)
	var service StoriesService = NewHackernewsStoriesProxy(client, storiesCache)

	// WHEN
	stories, err := service.GetTopStories(uint32(len(topStoriesIds)))

	// THEN
	if err != nil {
		t.Error("no error should be met")
	} else if len(*stories) != len(topStoriesIds) {
		t.Error("different number of stories found than expected")
	}

	actualStory := (*stories)[0]
	if actualStory.Id != storyId {
		t.Errorf("Found actual id '%d' but expected '%d'", actualStory.Id, storyId)
	}

	storyFromCache, storyIsCached := storiesCache.Get(storyId)
	if !storyIsCached {
		t.Error("story should have been added to cache")
	} else if storyFromCache.Id != storyId {
		t.Errorf("found story with id '%d' in cache but expected '%d' instead", storyFromCache.Id, storyId)
	}
}

func TestGetTopStoriesShouldReturnErrorIfFetchStoryFails(t *testing.T) {
	// GIVEN
	storyId := 0
	topStoriesIds := []int{storyId}

	mockLiveService := MockHnLiveService{}
	mockLiveService.MockedTopStories = func() ([]int, error) {
		return topStoriesIds, nil
	}

	mockItemService := MockHnItemService{}
	mockItemService.MockedItem = func(id int) (*hn.Item, error) {
		return nil, errors.New("item fetch fail")
	}

	var client hn.Client = hn.Client{Live: mockLiveService, Items: mockItemService}
	var storiesCache cache.Cache[int, *Story] = cache.NewTimeToLiveCache[int, *Story](10)
	var service StoriesService = NewHackernewsStoriesProxy(client, storiesCache)

	// WHEN
	stories, err := service.GetTopStories(uint32(len(topStoriesIds)))

	// THEN
	if err == nil {
		t.Error("should encounter an error on item fetch failure")
	} else if stories != nil {
		t.Error("should not return any story")
	}

	_, storyIsCached := storiesCache.Get(storyId)
	if storyIsCached {
		t.Error("story should not have been cached because it couldn't be fetched")
	}
}

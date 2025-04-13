package stories

type StoriesService interface {
	GetTopStories(maxStoryCount uint32) (*[]Story, error)
}
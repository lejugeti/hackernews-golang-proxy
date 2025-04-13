package users

import (
	"errors"
	"hackernews/server/cache"
	"testing"

	hn "github.com/peterhellberg/hn"
)

type MockHnUserService struct {
	MockedGet func(id string) (*hn.User, error)
}

func (m MockHnUserService) Get(id string) (*hn.User, error) {
	return m.MockedGet(id)
}

func TestGetUserInfoShouldErrorIfNicknameEmpty(t *testing.T) {
	// GIVEN
	var client hn.Client = hn.Client{}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

    nickname := ""

	// WHEN
	_, err := service.GetUserInfo(nickname)

	// THEN
	if err == nil {
		t.Error("error should be raised if nickname empty")
	}
}

func TestGetUserInfoShouldGetUserFromCache(t *testing.T) {
	// GIVEN
    var client hn.Client = hn.Client{}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

	nickname := "antwan"
	userCache.Add(nickname, &User{})

	// WHEN
	user, err := service.GetUserInfo(nickname)

	// THEN
	if err != nil {
		t.Error("no error should occured")
	} else if user == nil {
		t.Error("user should exist")
	}
}

func TestGetUserInfoShouldGetUserFromCacheEvenIfNil(t *testing.T) {
    // GIVEN
	var client hn.Client = hn.Client{}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

	nickname := "antwan"
	userCache.Add(nickname, nil)

	// WHEN
	user, err := service.GetUserInfo(nickname)

	// THEN
	if err != nil {
		t.Error("no error should occured")
	} else if user != nil {
		t.Error("user should be nil")
	}
}

func TestGetUserInfoShouldFetchUserFromHn(t *testing.T) {
	// GIVEN
	mockUserService := MockHnUserService{}
	mockUserService.MockedGet = func(id string) (*hn.User, error) {
		hnUser := hn.User{
			ID: "id",
			About: "about",
			Karma: 123,
			Created: 1234,
		}
		return &hnUser, nil
	}

	var client hn.Client = hn.Client{Users: mockUserService}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

	nickname := "antwan"

	// WHEN
	user, err := service.GetUserInfo(nickname)

	// THEN
	if err != nil {
		t.Error("no error should occured")
	} else if user == nil {
		t.Error("user should exist")
	}

	_, userExists := userCache.Get(nickname)
	if !userExists {
		t.Error("user should exist in cache")
	}
}

func TestGetUserInfoShouldFetchUserFromHnEvenIfUserNotFound(t *testing.T) {
	// GIVEN
	mockUserService := MockHnUserService{}
	mockUserService.MockedGet = func(id string) (*hn.User, error) {
		return &hn.User{}, nil
	}
	
	var client hn.Client = hn.Client{Users: mockUserService}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

	nickname := "antwan"

	// WHEN
	user, err := service.GetUserInfo(nickname)

	// THEN
	if err != nil {
		t.Error("no error should occured")
	} else if user != nil {
		t.Error("user should not have been found")
	}

	user, userExists := userCache.Get(nickname)
	if !userExists {
		t.Error("user should exist in cache")
	}
	if user != nil {
		t.Error("user should be nil")
	}
}

func TestGetUserInfoShouldReturnErrorIfFetchFails(t *testing.T) {
	// GIVEN
	mockUserService := MockHnUserService{}
	mockUserService.MockedGet = func(id string) (*hn.User, error) {
		return nil, errors.New("fetch user info fail")
	}
	
	var client hn.Client = hn.Client{Users: mockUserService}
	var userCache cache.Cache[string, *User] = cache.NewTimeToLiveCache[string, *User](10)
	var service UserService = NewHackernewsUserProxy(client, userCache)

	nickname := "antwan"

	// WHEN
	user, err := service.GetUserInfo(nickname)

	// THEN
	if err == nil {
		t.Error("error should have been returned")
	} else if user != nil {
		t.Error("user should not have been found")
	}

	_, userExists := userCache.Get(nickname)
	if userExists {
		t.Error("user should not exist in cache")
	}
}
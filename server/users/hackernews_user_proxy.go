package users

import (
	"time"

	"hackernews/server/cache"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	hn "github.com/peterhellberg/hn"
)

type hackernewsUserProxy struct {
	hnClient hn.Client
	cache cache.Cache[string, *User]
}

func NewHackernewsUserProxy(client hn.Client, cache cache.Cache[string, *User]) (UserService) {
	return &hackernewsUserProxy{
		hnClient: client,
		cache: cache,
	}
}

func (us *hackernewsUserProxy) GetUserInfo(nickname string) (*User, error) {
	if nickname == "" {
		return nil, status.Error(codes.InvalidArgument, "user nickname is required to get user info")
	}

	userFromCache, userIsCached := us.cache.Get(nickname)

	if userIsCached {
		return userFromCache, nil
	} else {
		user, err := us.fetchUserDetails(nickname)

		if err != nil {
			return nil, status.Errorf(codes.Internal, "error occurred while fetching user '%s' details. Cause: %v", nickname, err)
		}

		us.cache.Add(nickname, user)

		return user, nil
	}
}

func (us *hackernewsUserProxy) fetchUserDetails(nickname string) (*User, error) {
	if nickname == "" {
		return nil, status.Error(codes.InvalidArgument, "user nickname must be provided in order to fetch user details")
	}

	userInfo, err := us.hnClient.User(nickname)
	
	if err != nil {
		return nil, err
  	} else if us.userNotFound(userInfo) {
		return nil, nil
	}

	var user = User{
		Nickname: userInfo.ID,
		About: userInfo.About,
		Karma: uint64(userInfo.Karma),
		Joined: time.Unix(int64(userInfo.Created), 0)}

	return &user, nil
}

func (us *hackernewsUserProxy) userNotFound(user *hn.User) bool {
	return user.ID == "" && user.About == "" && user.Karma == 0 && user.Created == 0
}

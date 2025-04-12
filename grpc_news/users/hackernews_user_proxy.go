package users

import (
	"time"

	hn "github.com/peterhellberg/hn"
)

type hackernewsUserProxy struct {}

func (us *hackernewsUserProxy) GetUserInfo(nickname string) (*User, error) {
	userInfo, err := hn.DefaultClient.User(nickname)
	
	if err != nil {
		return nil, err
  	}

	var user = User{
		Nickname: userInfo.ID,
		About: userInfo.About,
		Karma: uint64(userInfo.Karma),
		Joined: time.Unix(int64(userInfo.Created), 0)}

	return &user, nil
}

func NewHackernewsUserProxy() (UserService) {
	return &hackernewsUserProxy{}
}
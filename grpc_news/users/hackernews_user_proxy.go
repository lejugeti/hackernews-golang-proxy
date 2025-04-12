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

func NewHackernewsUserProxy() (UserService) {
	return &hackernewsUserProxy{}
}
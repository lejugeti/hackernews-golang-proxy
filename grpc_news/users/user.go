package users

import "time"

type User struct {
	Nickname string;
	Karma uint64;
	About string;
	Joined time.Time;
}
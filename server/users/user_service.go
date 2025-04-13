package users

type UserService interface {
	GetUserInfo(nickname string) (*User, error)
}
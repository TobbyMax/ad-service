package user

type User struct {
	ID       int64
	Nickname string `validate:"min:1"`
	Email    string `validate:"min:1"`
}

package ads

import "time"

type Ad struct {
	ID          int64
	Title       string `validate:"min:1; max:99"`
	Text        string `validate:"min:1; max:499"`
	AuthorID    int64
	Published   bool
	DateCreated time.Time
	DateChanged time.Time
}

type AdList struct {
	Data []Ad
}

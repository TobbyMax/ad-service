package app

import "time"

const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "Mon, 2 Jan 2006 15:04:05 MST"
)

func ParseDate(s *string) (*time.Time, error) {
	if s == nil {
		return nil, nil
	}
	date, err := time.Parse(DateLayout, *s)
	if err != nil {
		return nil, err
	}
	return &date, nil
}

func FormatDate(date time.Time) string {
	return date.Format(DateTimeLayout)
}

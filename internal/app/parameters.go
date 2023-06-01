package app

import "time"

type ListAdsParams struct {
	Published *bool
	Uid       *int64
	Date      *time.Time
	Title     *string
}

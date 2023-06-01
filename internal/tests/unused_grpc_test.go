package tests

import (
	"github.com/stretchr/testify/assert"
	grpcSvc "homework10/internal/ports/grpc"
	"testing"
)

// For percentage!

func TestCreateAdRequest(t *testing.T) {
	var uid int64 = 2009
	testcases := []*grpcSvc.CreateAdRequest{
		{Title: "Dang!", Text: "Favourite Part", UserId: &uid},
		{UserId: &uid},
		{Title: "Dang!", Text: "Favourite Part"},
		{},
		nil,
	}

	for _, tc := range testcases {
		if tc == nil || tc.UserId == nil {
			assert.Equal(t, int64(0), tc.GetUserId())
		} else {
			assert.Equal(t, uid, tc.GetUserId())
		}
		if tc == nil {
			assert.Equal(t, "", tc.GetTitle())
			assert.Equal(t, "", tc.GetText())
		} else {
			assert.Equal(t, tc.Title, tc.GetTitle())
			assert.Equal(t, tc.Text, tc.GetText())
		}
	}
}

func TestChangeAdStatusRequest(t *testing.T) {
	var (
		uid int64 = 2009
		id  int64 = 2016
	)
	testcases := []*grpcSvc.ChangeAdStatusRequest{
		{AdId: &id, UserId: &uid, Published: true},
		{AdId: &id, UserId: &uid},
		{UserId: &uid, Published: true},
		{AdId: &id, Published: true},
		{},
		nil,
	}

	for _, tc := range testcases {
		if tc == nil || tc.UserId == nil {
			assert.Equal(t, int64(0), tc.GetUserId())
		} else {
			assert.Equal(t, uid, tc.GetUserId())
		}
		if tc == nil || tc.AdId == nil {
			assert.Equal(t, int64(0), tc.GetAdId())
		} else {
			assert.Equal(t, id, tc.GetAdId())
		}
		if tc == nil {
			assert.Equal(t, false, tc.GetPublished())
		} else {
			assert.Equal(t, tc.Published, tc.GetPublished())
		}
	}
}

func TestUpdateAdRequest(t *testing.T) {
	var (
		uid int64 = 2009
		id  int64 = 2016
	)
	testcases := []*grpcSvc.UpdateAdRequest{
		{AdId: &id, Title: "Dang!", Text: "Favourite Part", UserId: &uid},
		{AdId: &id, UserId: &uid},
		{Title: "Dang!", Text: "Favourite Part", UserId: &uid},
		{AdId: &id, Title: "Dang!", Text: "Favourite Part"},
		{Title: "Dang!", Text: "Favourite Part"},
		{},
		nil,
	}

	for _, tc := range testcases {
		if tc == nil || tc.UserId == nil {
			assert.Equal(t, int64(0), tc.GetUserId())
		} else {
			assert.Equal(t, uid, tc.GetUserId())
		}
		if tc == nil || tc.AdId == nil {
			assert.Equal(t, int64(0), tc.GetAdId())
		} else {
			assert.Equal(t, id, tc.GetAdId())
		}
		if tc == nil {
			assert.Equal(t, "", tc.GetTitle())
			assert.Equal(t, "", tc.GetText())
		} else {
			assert.Equal(t, tc.Title, tc.GetTitle())
			assert.Equal(t, tc.Text, tc.GetText())
		}
	}
}

func TestListAdRequest(t *testing.T) {
	var (
		uid   int64 = 2009
		pub         = true
		date        = "2022-01-01"
		title       = "Dang!"
	)
	testcases := []*grpcSvc.ListAdRequest{
		{UserId: &uid, Published: &pub, Date: &date, Title: &title},
		{Published: &pub, Date: &date, Title: &title},
		{UserId: &uid, Date: &date, Title: &title},
		{UserId: &uid, Published: &pub, Title: &title},
		{UserId: &uid, Published: &pub, Date: &date},
		{},
		nil,
	}

	for _, tc := range testcases {
		if tc == nil || tc.UserId == nil {
			assert.Equal(t, int64(0), tc.GetUserId())
		} else {
			assert.Equal(t, *tc.UserId, tc.GetUserId())

		}
		if tc == nil || tc.Published == nil {
			assert.Equal(t, false, tc.GetPublished())
		} else {
			assert.Equal(t, *tc.Published, tc.GetPublished())
		}
		if tc == nil || tc.Date == nil {
			assert.Equal(t, "", tc.GetDate())
		} else {
			assert.Equal(t, *tc.Date, tc.GetDate())

		}
		if tc == nil || tc.Title == nil {
			assert.Equal(t, "", tc.GetTitle())
		} else {
			assert.Equal(t, *tc.Title, tc.GetTitle())
		}
	}
}

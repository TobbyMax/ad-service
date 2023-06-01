package grpc

import (
	"errors"
	"github.com/TobbyMax/validator"
	"google.golang.org/grpc/codes"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/user"
)

var ErrMissingArgument = errors.New("required argument is missing")

func AdSuccessResponse(ad *ads.Ad) *AdResponse {
	return &AdResponse{
		Id:          ad.ID,
		Title:       ad.Title,
		Text:        ad.Text,
		AuthorId:    ad.AuthorID,
		Published:   ad.Published,
		DateCreated: app.FormatDate(ad.DateCreated),
		DateChanged: app.FormatDate(ad.DateChanged),
	}
}

func AdListSuccessResponse(al *ads.AdList) *ListAdResponse {
	response := ListAdResponse{List: make([]*AdResponse, 0)}

	for _, ad := range al.Data {
		response.List = append(response.List, AdSuccessResponse(&ad))
	}
	return &response
}

func UserSuccessResponse(u *user.User) *UserResponse {
	return &UserResponse{
		Id:    u.ID,
		Name:  u.Nickname,
		Email: u.Email,
	}
}

func GetErrorCode(err error) codes.Code {
	switch {
	case errors.As(err, &validator.ValidationErrors{}):
		return codes.InvalidArgument
	case errors.Is(err, app.ErrForbidden):
		return codes.PermissionDenied
	case errors.Is(err, app.ErrAdNotFound):
		fallthrough
	case errors.Is(err, app.ErrUserNotFound):
		return codes.NotFound
	}
	return codes.Internal
}

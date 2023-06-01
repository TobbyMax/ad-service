package app

import (
	"context"
	"fmt"
	"github.com/TobbyMax/validator"
	"homework10/internal/ads"
	"homework10/internal/user"
	"time"
)

var (
	ErrForbidden    = fmt.Errorf("forbidden")
	ErrAdNotFound   = fmt.Errorf("ad with such id does not exist")
	ErrUserNotFound = fmt.Errorf("user with such id does not exist")
)

type AdApp interface {
	CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error)
	ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error)
	UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error)
	GetAd(ctx context.Context, id int64) (*ads.Ad, error)
	DeleteAd(ctx context.Context, id int64, uid int64) error

	ListAds(ctx context.Context, params ListAdsParams) (*ads.AdList, error)
}

type UserApp interface {
	CreateUser(ctx context.Context, nickname string, email string) (*user.User, error)
	GetUser(ctx context.Context, id int64) (*user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname string, email string) (*user.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

type App interface {
	AdApp
	UserApp
}

type AdRepository interface {
	AddAd(ctx context.Context, ad ads.Ad) (int64, error)
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)
	UpdateAdStatus(ctx context.Context, id int64, published bool, date time.Time) error
	UpdateAdContent(ctx context.Context, id int64, title string, text string, date time.Time) error
	DeleteAdByID(ctx context.Context, id int64) error

	GetAdList(ctx context.Context, params ListAdsParams) (*ads.AdList, error)
}

type UserRepository interface {
	AddUser(ctx context.Context, u user.User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*user.User, error)
	UpdateUser(ctx context.Context, id int64, nickname string, email string) error
	DeleteUserByID(ctx context.Context, id int64) error
}

type Repository interface {
	AdRepository
	UserRepository
}

type Application struct {
	repository Repository
}

func NewApp(repo Repository) App {
	return NewAdApp(repo)
}

func NewAdApp(repo Repository) *Application {
	return &Application{repository: repo}
}

func (a Application) CreateAd(ctx context.Context, title string, text string, uid int64) (*ads.Ad, error) {
	ad := ads.Ad{Title: title, Text: text, AuthorID: uid, Published: false, DateCreated: time.Now().UTC()}
	ad.DateChanged = ad.DateCreated
	if err := validator.Validate(ad); err != nil {
		return nil, err
	}

	id, err := a.repository.AddAd(ctx, ad)
	if err != nil {
		return nil, err
	}
	ad.ID = id

	return &ad, nil
}

func (a Application) GetAd(ctx context.Context, id int64) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a Application) ChangeAdStatus(ctx context.Context, id int64, uid int64, published bool) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != uid {
		return nil, ErrForbidden
	}

	ad.Published = published
	ad.DateChanged = time.Now().UTC()

	err = a.repository.UpdateAdStatus(ctx, id, published, ad.DateChanged)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a Application) UpdateAd(ctx context.Context, id int64, uid int64, title string, text string) (*ads.Ad, error) {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ad.AuthorID != uid {
		return nil, ErrForbidden
	}

	ad.Title = title
	ad.Text = text
	ad.DateChanged = time.Now().UTC()

	if err := validator.Validate(*ad); err != nil {
		return nil, err
	}

	err = a.repository.UpdateAdContent(ctx, id, title, text, ad.DateChanged)
	if err != nil {
		return nil, err
	}

	return ad, nil
}

func (a Application) ListAds(ctx context.Context, params ListAdsParams) (*ads.AdList, error) {
	p := true
	if params.Published == nil && params.Uid == nil && params.Date == nil && params.Title == nil {
		params.Published = &p
	}
	al, err := a.repository.GetAdList(ctx, params)

	if err != nil {
		return nil, err
	}

	return al, nil
}

func (a Application) CreateUser(ctx context.Context, nickname string, email string) (*user.User, error) {
	u := user.User{Nickname: nickname, Email: email}

	if err := validator.Validate(u); err != nil {
		return nil, err
	}

	id, err := a.repository.AddUser(ctx, u)
	if err != nil {
		return nil, err
	}
	u.ID = id

	return &u, nil
}

func (a Application) GetUser(ctx context.Context, id int64) (*user.User, error) {
	u, err := a.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a Application) UpdateUser(ctx context.Context, id int64, nickname string, email string) (*user.User, error) {
	u, err := a.repository.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.Nickname = nickname
	u.Email = email

	if err := validator.Validate(*u); err != nil {
		return nil, err
	}

	err = a.repository.UpdateUser(ctx, id, nickname, email)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (a Application) DeleteAd(ctx context.Context, id int64, uid int64) error {
	ad, err := a.repository.GetAdByID(ctx, id)
	if err != nil {
		return err
	}
	if ad.AuthorID != uid {
		return ErrForbidden
	}
	err = a.repository.DeleteAdByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (a Application) DeleteUser(ctx context.Context, id int64) error {
	err := a.repository.DeleteUserByID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

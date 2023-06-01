package adrepo

import (
	"context"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/user"
	"sync"
	"time"
)

func New() app.Repository {
	return NewRepositoryMap()
}

type RepositoryMap struct {
	sync.Mutex
	adTable   map[int64]ads.Ad
	userTable map[int64]user.User
	user2ads  map[int64]map[int64]struct{}
}

func NewRepositoryMap() *RepositoryMap {
	return &RepositoryMap{
		adTable:   make(map[int64]ads.Ad),
		userTable: make(map[int64]user.User),
		user2ads:  make(map[int64]map[int64]struct{}),
	}
}

func (r *RepositoryMap) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.userTable[ad.AuthorID]; !ok {
		return 0, app.ErrUserNotFound
	}
	ad.ID = int64(len(r.adTable))
	r.adTable[ad.ID] = ad
	r.user2ads[ad.AuthorID][ad.ID] = struct{}{}
	return ad.ID, nil
}

func (r *RepositoryMap) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	r.Lock()
	defer r.Unlock()
	if ad, ok := r.adTable[id]; !ok {
		return nil, app.ErrAdNotFound
	} else {
		return &ad, nil
	}
}

func (r *RepositoryMap) UpdateAdStatus(ctx context.Context, id int64, published bool, date time.Time) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.adTable[id]; !ok {
		return app.ErrAdNotFound
	}
	ad := r.adTable[id]
	ad.Published = published
	ad.DateChanged = date
	r.adTable[id] = ad
	return nil
}

func (r *RepositoryMap) UpdateAdContent(ctx context.Context, id int64, title string, text string, date time.Time) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.adTable[id]; !ok {
		return app.ErrAdNotFound
	}
	ad := r.adTable[id]
	ad.Title = title
	ad.Text = text
	ad.DateChanged = date
	r.adTable[id] = ad
	return nil
}

func (r *RepositoryMap) GetAdList(ctx context.Context, params app.ListAdsParams) (*ads.AdList, error) {
	r.Lock()
	defer r.Unlock()
	al := ads.AdList{Data: make([]ads.Ad, 0)}
	for _, ad := range r.adTable {
		if params.Published == nil || *params.Published == ad.Published {
			if (params.Uid == nil || *params.Uid == ad.AuthorID) && (params.Title == nil || *params.Title == ad.Title) {
				if year, month, day := ad.DateCreated.Date(); params.Date == nil ||
					(params.Date.Year() == year && params.Date.Month() == month && params.Date.Day() == day) {
					al.Data = append(al.Data, ad)
				}
			}
		}
	}
	return &al, nil
}

func (r *RepositoryMap) AddUser(ctx context.Context, u user.User) (int64, error) {
	r.Lock()
	defer r.Unlock()
	u.ID = int64(len(r.userTable))
	r.userTable[u.ID] = u
	r.user2ads[u.ID] = make(map[int64]struct{})
	return u.ID, nil
}

func (r *RepositoryMap) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	r.Lock()
	defer r.Unlock()
	if u, ok := r.userTable[id]; !ok {
		return nil, app.ErrUserNotFound
	} else {
		return &u, nil
	}
}

func (r *RepositoryMap) UpdateUser(ctx context.Context, id int64, nickname string, email string) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.userTable[id]; !ok {
		return app.ErrUserNotFound
	}
	u := r.userTable[id]
	u.Nickname = nickname
	u.Email = email
	r.userTable[id] = u
	return nil
}

func (r *RepositoryMap) DeleteAdByID(ctx context.Context, id int64) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.adTable[id]; !ok {
		return app.ErrAdNotFound
	}
	delete(r.adTable, id)
	return nil
}

func (r *RepositoryMap) DeleteUserByID(ctx context.Context, id int64) error {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.userTable[id]; !ok {
		return app.ErrUserNotFound
	}
	for adID := range r.user2ads[id] {
		delete(r.adTable, adID)
	}
	delete(r.user2ads, id)
	delete(r.userTable, id)
	return nil
}

package adrepo

import (
	"context"
	"github.com/TobbyMax/ad-service.git/internal/ads"
	"github.com/TobbyMax/ad-service.git/internal/app"
	"github.com/TobbyMax/ad-service.git/internal/user"
	"github.com/jackc/pgx/v5"
	"time"
)

type RepositoryPG struct {
	conn *pgx.Conn
}

func NewRepositoryPG(conn *pgx.Conn) *RepositoryPG {
	return &RepositoryPG{conn: conn}
}

func (r *RepositoryPG) GetUserByID(ctx context.Context, user_id int64) (*user.User, error) {
	q := `select u.id, u.name from users u where u.id = $1`

	u := &user.User{}

	if err := r.conn.QueryRow(ctx, q, user_id).Scan(&u.ID, &u.Nickname, &u.Email); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *RepositoryPG) AddUser(ctx context.Context, user user.User) (int64, error) {
	q := `insert into users(name) values($1) returning id`

	var user_id int64

	if err := r.conn.QueryRow(ctx, q, user.Nickname).Scan(&user_id); err != nil {
		return user_id, err
	}

	return user_id, nil
}

func (r *RepositoryPG) UpdateUser(ctx context.Context, id int64, nickname string, email string) error {
	q := `update users set nickname = $1 email = $2 where id = $3`

	commandTag, err := r.conn.Exec(ctx, q, nickname, email, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return app.ErrUserNotFound
	}

	return nil
}

func (r *RepositoryPG) DeleteAdByID(ctx context.Context, id int64) error {

}

func (r *RepositoryPG) DeleteUserByID(ctx context.Context, id int64) error {
	q := `delete from users where id=$1`

	commandTag, err := r.conn.Exec(ctx, q, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return app.ErrUserNotFound
	}

	return nil
}

func (r *RepositoryPG) AddAd(ctx context.Context, ad ads.Ad) (int64, error) {
	//r.Lock()
	//defer r.Unlock()
	//if _, ok := r.userTable[ad.AuthorID]; !ok {
	//	return 0, app.ErrUserNotFound
	//}
	//ad.ID = int64(len(r.adTable))
	//r.adTable[ad.ID] = ad
	//r.user2ads[ad.AuthorID][ad.ID] = struct{}{}
	//return ad.ID, nil
}

func (r *RepositoryPG) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	//r.Lock()
	//defer r.Unlock()
	//if ad, ok := r.adTable[id]; !ok {
	//	return nil, app.ErrAdNotFound
	//} else {
	//	return &ad, nil
	//}
}

func (r *RepositoryPG) UpdateAdStatus(ctx context.Context, id int64, published bool, date time.Time) error {
	//r.Lock()
	//defer r.Unlock()
	//if _, ok := r.adTable[id]; !ok {
	//	return app.ErrAdNotFound
	//}
	//ad := r.adTable[id]
	//ad.Published = published
	//ad.DateChanged = date
	//r.adTable[id] = ad
	//return nil
}

func (r *RepositoryPG) UpdateAdContent(ctx context.Context, id int64, title string, text string, date time.Time) error {
	//r.Lock()
	//defer r.Unlock()
	//if _, ok := r.adTable[id]; !ok {
	//	return app.ErrAdNotFound
	//}
	//ad := r.adTable[id]
	//ad.Title = title
	//ad.Text = text
	//ad.DateChanged = date
	//r.adTable[id] = ad
	//return nil
}

func (r *RepositoryPG) GetAdList(ctx context.Context, params app.ListAdsParams) (*ads.AdList, error) {
	//r.Lock()
	//defer r.Unlock()
	//al := ads.AdList{Data: make([]ads.Ad, 0)}
	//for _, ad := range r.adTable {
	//	if params.Published == nil || *params.Published == ad.Published {
	//		if (params.Uid == nil || *params.Uid == ad.AuthorID) && (params.Title == nil || *params.Title == ad.Title) {
	//			if year, month, day := ad.DateCreated.Date(); params.Date == nil ||
	//				(params.Date.Year() == year && params.Date.Month() == month && params.Date.Day() == day) {
	//				al.Data = append(al.Data, ad)
	//			}
	//		}
	//	}
	//}
	//return &al, nil
}

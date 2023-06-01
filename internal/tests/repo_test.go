package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/user"
	"log"
	"testing"
	"time"
)

type RepoSuite struct {
	suite.Suite
	Repo app.Repository
	Ctx  context.Context
}

func (suite *RepoSuite) SetupTest() {
	log.Println("Setting Up Test")
	suite.Ctx = context.Background()
	suite.Repo = adrepo.New()
}

func (suite *RepoSuite) TearDownTest() {
	log.Println("Tearing Down Test")
}

func (suite *RepoSuite) TestRepo_AddUser() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	id, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	suite.Equal(int64(0), id)
}

func (suite *RepoSuite) TestRepo_AddMultipleUsers() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	for i := 0; i < 100; i++ {
		id, err := suite.Repo.AddUser(suite.Ctx, u)
		suite.NoError(err)
		suite.Equal(int64(i), id)
	}
}

func (suite *RepoSuite) TestRepo_GetUser() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	id, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	res, err := suite.Repo.GetUserByID(suite.Ctx, id)
	suite.NoError(err)
	suite.Equal(u, *res)
}

func (suite *RepoSuite) TestRepo_GetUserError() {
	res, err := suite.Repo.GetUserByID(suite.Ctx, 1)
	suite.Error(err)
	suite.Nil(res)
	suite.ErrorIs(app.ErrUserNotFound, err)
}

func (suite *RepoSuite) TestRepo_UpdateUser() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	id, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	err = suite.Repo.UpdateUser(suite.Ctx, id, "KDot", "money@trees.com")
	suite.NoError(err)
	res, err := suite.Repo.GetUserByID(suite.Ctx, id)
	suite.NoError(err)
	suite.Equal(user.User{ID: 0, Nickname: "KDot", Email: "money@trees.com"}, *res)
}

func (suite *RepoSuite) TestRepo_UpdateUserError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	_, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	err = suite.Repo.UpdateUser(suite.Ctx, 1, "KDot", "money@trees.com")
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *RepoSuite) TestRepo_DeleteUser() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	id, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	err = suite.Repo.DeleteUserByID(suite.Ctx, id)
	suite.NoError(err)
	_, err = suite.Repo.GetUserByID(suite.Ctx, id)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *RepoSuite) TestRepo_DeleteUserError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	_, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	err = suite.Repo.DeleteUserByID(suite.Ctx, 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *RepoSuite) TestRepo_AddAd() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	id, err := suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	suite.Equal(int64(0), id)
}

func (suite *RepoSuite) TestRepo_AddAdError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	_, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: 2009}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *RepoSuite) TestRepo_GetAd() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	id, err := suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	res, err := suite.Repo.GetAdByID(suite.Ctx, id)
	suite.NoError(err)
	suite.Equal(ad, *res)
}

func (suite *RepoSuite) TestRepo_GetAdError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	_, err = suite.Repo.GetAdByID(suite.Ctx, 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *RepoSuite) TestRepo_UpdateAdStatus() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	id, err := suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	t := time.Now().UTC()
	err = suite.Repo.UpdateAdStatus(suite.Ctx, id, true, t)
	suite.NoError(err)
	res, err := suite.Repo.GetAdByID(suite.Ctx, id)
	suite.NoError(err)
	suite.Equal(id, res.ID)
	suite.Equal(ad.Title, res.Title)
	suite.Equal(ad.Text, res.Text)
	suite.Equal(true, res.Published)
	suite.Equal(t, res.DateChanged)
}

func (suite *RepoSuite) TestRepo_UpdateAdStatusError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	t := time.Now().UTC()
	err = suite.Repo.UpdateAdStatus(suite.Ctx, 1, true, t)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *RepoSuite) TestRepo_UpdateAdContent() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	id, err := suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	t := time.Now().UTC()
	err = suite.Repo.UpdateAdContent(suite.Ctx, id, "Apparently", "by J.Cole", t)
	suite.NoError(err)
	res, err := suite.Repo.GetAdByID(suite.Ctx, id)
	suite.NoError(err)
	suite.Equal(id, res.ID)
	suite.Equal("Apparently", res.Title)
	suite.Equal("by J.Cole", res.Text)
	suite.Equal(false, res.Published)
	suite.Equal(t, res.DateChanged)
}

func (suite *RepoSuite) TestRepo_UpdateAdContentError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	t := time.Now().UTC()
	err = suite.Repo.UpdateAdContent(suite.Ctx, 1, "Apparently", "by J.Cole", t)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *RepoSuite) TestRepo_DeleteAd() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	id, err := suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	err = suite.Repo.DeleteAdByID(suite.Ctx, id)
	suite.NoError(err)
	_, err = suite.Repo.GetAdByID(suite.Ctx, id)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *RepoSuite) TestRepo_DeleteAdError() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	err = suite.Repo.DeleteAdByID(suite.Ctx, 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *RepoSuite) TestRepo_GetAdList() {
	u := user.User{Nickname: "Mac Miller", Email: "swimmig@circles.com"}
	uid, err := suite.Repo.AddUser(suite.Ctx, u)
	suite.NoError(err)
	ad := ads.Ad{Title: "Dang!", Text: "The Divine Feminine", AuthorID: uid}
	_, err = suite.Repo.AddAd(suite.Ctx, ad)
	suite.NoError(err)
	res, err := suite.Repo.GetAdList(suite.Ctx, app.ListAdsParams{})
	suite.NoError(err)
	suite.Len(res.Data, 1)
	suite.Equal(ad, res.Data[0])
}

func TestRepo(t *testing.T) {
	suite.Run(t, new(RepoSuite))
}

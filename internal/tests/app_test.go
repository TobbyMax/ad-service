package tests

import (
	"context"
	"github.com/TobbyMax/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/tests/mocks"
	"homework10/internal/user"
	"testing"
)

type AppTestSuite struct {
	suite.Suite
	Repo *mocks.Repository
	Ctx  context.Context
}

func (suite *AppTestSuite) SetupTest() {
	suite.Repo = mocks.NewRepository(suite.T())
	suite.Ctx = context.Background()
}

func (suite *AppTestSuite) TestApp_CreateAd() {
	id := int64(13)
	suite.Repo.On("AddAd", suite.Ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, nil).
		Once()
	service := app.NewApp(suite.Repo)
	ad, err := service.CreateAd(suite.Ctx, "title", "text", 1)
	suite.Nil(err)
	suite.Equal(id, ad.ID)
	suite.Equal("title", ad.Title)
	suite.Equal("text", ad.Text)
	suite.Equal(int64(1), ad.AuthorID)
}

func (suite *AppTestSuite) TestApp_CreateAd_NonExistentUser() {
	id := int64(13)
	suite.Repo.On("AddAd", suite.Ctx, mock.AnythingOfType("ads.Ad")).
		Return(id, app.ErrUserNotFound).
		Once()
	service := app.NewApp(suite.Repo)
	_, err := service.CreateAd(suite.Ctx, "title", "text", 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *AppTestSuite) TestApp_GetAd() {
	suite.Repo.On("GetAdByID", suite.Ctx, int64(0)).
		Return(&ads.Ad{}, nil).
		Once()
	service := app.NewApp(suite.Repo)
	ad, err := service.GetAd(suite.Ctx, 0)
	suite.Nil(err)
	suite.Equal(ads.Ad{}, *ad)
}

func (suite *AppTestSuite) TestApp_CreateAd_InvalidTitle() {
	service := app.NewApp(suite.Repo)
	_, err := service.CreateAd(suite.Ctx, "", "text", 1)
	suite.Error(err)
	e := &validator.ValidationErrors{}
	suite.ErrorAs(err, e)
}

func (suite *AppTestSuite) TestApp_CreateAd_InvalidText() {
	service := app.NewApp(suite.Repo)
	_, err := service.CreateAd(suite.Ctx, "title", "", 1)
	suite.Error(err)
	e := &validator.ValidationErrors{}
	suite.ErrorAs(err, e)
}

func (suite *AppTestSuite) TestApp_UpdateAd() {
	id := int64(0)
	title := "title"
	text := "text"
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("UpdateAdContent", suite.Ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateAd(suite.Ctx, id, int64(1), title, text)
	suite.Nil(err)
}

func (suite *AppTestSuite) TestApp_UpdateAd_NonExistentAd() {
	id := int64(0)
	title := "title"
	text := "text"
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateAd(suite.Ctx, id, int64(1), title, text)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *AppTestSuite) TestApp_UpdateAd_Forbidden() {
	id := int64(0)
	title := "title"
	text := "text"
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateAd(suite.Ctx, id, int64(1), title, text)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrForbidden)
}

func (suite *AppTestSuite) TestApp_UpdateAd_RepoError() {
	id := int64(0)
	title := "title"
	text := "text"
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("UpdateAdContent", suite.Ctx, id, title, text, mock.AnythingOfType("time.Time")).
		Return(ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateAd(suite.Ctx, id, int64(1), title, text)
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
}

func (suite *AppTestSuite) TestApp_ChangeAdStatus() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("UpdateAdStatus", suite.Ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.ChangeAdStatus(suite.Ctx, id, int64(1), true)
	suite.Nil(err)
}

func (suite *AppTestSuite) TestApp_ChangeAdStatus_NonExistentAd() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.ChangeAdStatus(suite.Ctx, id, int64(1), true)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *AppTestSuite) TestApp_ChangeAdStatus_Forbidden() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.ChangeAdStatus(suite.Ctx, id, int64(1), true)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrForbidden)
}

func (suite *AppTestSuite) TestApp_ChangeAdStatus_RepoError() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("UpdateAdStatus", suite.Ctx, id, true, mock.AnythingOfType("time.Time")).
		Return(ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.ChangeAdStatus(suite.Ctx, id, int64(1), true)
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
}

func (suite *AppTestSuite) TestApp_ListAds() {
	pub := true
	params := app.ListAdsParams{Published: &pub}
	suite.Repo.On("GetAdList", suite.Ctx, params).
		Return(nil, nil).
		Once()

	service := app.NewApp(suite.Repo)
	al, err := service.ListAds(suite.Ctx, params)
	suite.Nil(err)
	suite.Nil(al)
}

func (suite *AppTestSuite) TestApp_ListAds_AllNil() {
	params := app.ListAdsParams{}
	pub := true
	suite.Repo.On("GetAdList", suite.Ctx, app.ListAdsParams{Published: &pub}).
		Return(nil, nil).
		Once()

	service := app.NewApp(suite.Repo)
	al, err := service.ListAds(suite.Ctx, params)
	suite.Nil(err)
	suite.Nil(al)
}

func (suite *AppTestSuite) TestApp_ListAds_RepoError() {
	params := app.ListAdsParams{}
	pub := true
	suite.Repo.On("GetAdList", suite.Ctx, app.ListAdsParams{Published: &pub}).
		Return(nil, ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	al, err := service.ListAds(suite.Ctx, params)
	suite.Nil(al)
	suite.ErrorIs(err, ErrMock)
}

func (suite *AppTestSuite) TestApp_CreateUser() {
	id := int64(13)
	suite.Repo.On("AddUser", suite.Ctx, mock.AnythingOfType("user.User")).
		Return(id, nil).
		Once()
	service := app.NewApp(suite.Repo)
	u, err := service.CreateUser(suite.Ctx, "Mac Miller", "swimming@circles.com")
	suite.Nil(err)
	suite.Equal(id, u.ID)
	suite.Equal("Mac Miller", u.Nickname)
	suite.Equal("swimming@circles.com", u.Email)
}

func (suite *AppTestSuite) TestApp_GetUser() {
	suite.Repo.On("GetUserByID", suite.Ctx, int64(0)).
		Return(nil, nil).
		Once()
	service := app.NewApp(suite.Repo)
	u, err := service.GetUser(suite.Ctx, 0)
	suite.Nil(err)
	suite.Nil(u)
}

func (suite *AppTestSuite) TestApp_CreateUser_InvalidName() {
	service := app.NewApp(suite.Repo)
	u, err := service.CreateUser(suite.Ctx, "", "swimming@circles.com")
	suite.Error(err)
	e := &validator.ValidationErrors{}
	suite.ErrorAs(err, e)
	suite.Nil(u)
}

func (suite *AppTestSuite) TestApp_CreateUser_InvalidEmail() {
	service := app.NewApp(suite.Repo)
	u, err := service.CreateUser(suite.Ctx, "Mac", "")
	suite.Error(err)
	e := &validator.ValidationErrors{}
	suite.ErrorAs(err, e)
	suite.Nil(u)
}

func (suite *AppTestSuite) TestApp_CreateUser_RepoError() {
	id := int64(13)
	suite.Repo.On("AddUser", suite.Ctx, mock.AnythingOfType("user.User")).
		Return(id, ErrMock).
		Once()
	service := app.NewApp(suite.Repo)
	u, err := service.CreateUser(suite.Ctx, "Mac", "swimming@circles.com")
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
	suite.Nil(u)
}

func (suite *AppTestSuite) TestApp_UpdateUser() {
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	suite.Repo.On("GetUserByID", suite.Ctx, id).
		Return(&user.User{}, nil).
		Once()
	suite.Repo.On("UpdateUser", suite.Ctx, id, name, email).
		Return(nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateUser(suite.Ctx, id, name, email)
	suite.Nil(err)
}

func (suite *AppTestSuite) TestApp_UpdateUser_NonExistentID() {
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	suite.Repo.On("GetUserByID", suite.Ctx, id).
		Return(nil, app.ErrUserNotFound).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateUser(suite.Ctx, id, name, email)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrUserNotFound)
}

func (suite *AppTestSuite) TestApp_UpdateUser_InvalidName() {
	id := int64(0)
	name := ""
	email := "swimming@circles.com"
	suite.Repo.On("GetUserByID", suite.Ctx, id).
		Return(&user.User{}, nil).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateUser(suite.Ctx, id, name, email)
	suite.Error(err)
	e := &validator.ValidationErrors{}
	suite.ErrorAs(err, e)
}

func (suite *AppTestSuite) TestApp_UpdateUser_RepoError() {
	id := int64(0)
	name := "Mac Miller"
	email := "swimming@circles.com"
	suite.Repo.On("GetUserByID", suite.Ctx, id).
		Return(&user.User{}, nil).
		Once()
	suite.Repo.On("UpdateUser", suite.Ctx, id, name, email).
		Return(ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	_, err := service.UpdateUser(suite.Ctx, id, name, email)
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
}

func (suite *AppTestSuite) TestApp_DeleteUser() {
	id := int64(0)
	suite.Repo.On("DeleteUserByID", suite.Ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteUser(suite.Ctx, id)
	suite.Nil(err)
}

func (suite *AppTestSuite) TestApp_DeleteUser_RepoError() {
	id := int64(0)
	suite.Repo.On("DeleteUserByID", suite.Ctx, id).
		Return(ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteUser(suite.Ctx, id)
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
}

func (suite *AppTestSuite) TestApp_DeleteAd() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("DeleteAdByID", suite.Ctx, id).
		Return(nil).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteAd(suite.Ctx, id, 1)
	suite.Nil(err)
}

func (suite *AppTestSuite) TestApp_DeleteAd_Forbidden() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 0}, nil).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteAd(suite.Ctx, id, 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrForbidden)
}

func (suite *AppTestSuite) TestApp_DeleteAd_NonExistentAd() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(nil, app.ErrAdNotFound).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteAd(suite.Ctx, id, 1)
	suite.Error(err)
	suite.ErrorIs(err, app.ErrAdNotFound)
}

func (suite *AppTestSuite) TestApp_DeleteAd_RepoError() {
	id := int64(0)
	suite.Repo.On("GetAdByID", suite.Ctx, id).
		Return(&ads.Ad{AuthorID: 1}, nil).
		Once()
	suite.Repo.On("DeleteAdByID", suite.Ctx, id).
		Return(ErrMock).
		Once()

	service := app.NewApp(suite.Repo)
	err := service.DeleteAd(suite.Ctx, id, 1)
	suite.Error(err)
	suite.ErrorIs(err, ErrMock)
}

func TestAppSuite(t *testing.T) {
	suite.Run(t, new(AppTestSuite))
}

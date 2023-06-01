package tests

import (
	grpcPort "homework10/internal/ports/grpc"
)

func (suite *GRPCSuite) TestGRRPCCreateUser() {
	res, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err, "suite.Client.GetUser")

	suite.Equal("Oleg", res.Name)
	suite.Equal("ivanov@yandex.ru", res.Email)
}

func (suite *GRPCSuite) TestGRRPCGetUser() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	res, err := suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{Id: &user.Id})
	suite.NoError(err)

	suite.Equal(int64(0), res.Id)
	suite.Equal("Oleg", res.Name)
	suite.Equal("ivanov@yandex.ru", res.Email)
}

func (suite *GRPCSuite) TestGRRPCDeleteUser() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	_, err = suite.Client.DeleteUser(suite.Context, &grpcPort.DeleteUserRequest{Id: &user.Id})

	suite.NoError(err)

	_, err = suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{Id: &user.Id})
	suite.Error(err)
	suite.Equal(ErrUserNotFound.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRRPCUpdateUser() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	res, err := suite.Client.UpdateUser(suite.Context, &grpcPort.UpdateUserRequest{Id: &user.Id, Name: "Kanye", Email: "graduation@west.com"})
	suite.NoError(err)
	suite.Equal("Kanye", res.Name)
	suite.Equal("graduation@west.com", res.Email)
}

func (suite *GRPCSuite) TestGRRPCCreateAd() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	res, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	suite.NoError(err)
	suite.Equal(int64(0), res.Id)
	suite.Equal("Forest", res.Title)
	suite.Equal("Hill Drive", res.Text)
	suite.Equal(int64(0), res.AuthorId)
	suite.Equal(false, res.Published)
}

func (suite *GRPCSuite) TestGRRPCChangeAdStatus() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	suite.NoError(err)

	res, err := suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &ad.Id, UserId: &user.Id, Published: true})
	suite.NoError(err)

	suite.Equal(int64(0), res.Id)
	suite.Equal("Forest", res.Title)
	suite.Equal("Hill Drive", res.Text)
	suite.Equal(int64(0), res.AuthorId)
	suite.Equal(true, res.Published)
}

func (suite *GRPCSuite) TestGRRPCUpdateAd() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	suite.NoError(err)

	res, err := suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{AdId: &ad.Id, UserId: &user.Id, Title: "Corny", Text: "Low Key"})
	suite.NoError(err)

	suite.Equal(int64(0), res.Id)
	suite.Equal("Corny", res.Title)
	suite.Equal("Low Key", res.Text)
	suite.Equal(int64(0), res.AuthorId)
	suite.Equal(false, res.Published)
}

func (suite *GRPCSuite) TestGRRPCGetAd() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	suite.NoError(err)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{AdId: &ad.Id, UserId: &user.Id, Title: "Corny", Text: "Low Key"})
	suite.NoError(err)
	res, err := suite.Client.GetAd(suite.Context, &grpcPort.GetAdRequest{AdId: &ad.Id})
	suite.NoError(err)
	suite.Equal(int64(0), res.Id)
	suite.Equal("Corny", res.Title)
	suite.Equal("Low Key", res.Text)
	suite.Equal(int64(0), res.AuthorId)
	suite.Equal(false, res.Published)
}

func (suite *GRPCSuite) TestGRRPCDeleteAd() {
	user, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "ivanov@yandex.ru"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user.Id})
	suite.NoError(err)

	_, err = suite.Client.DeleteAd(suite.Context, &grpcPort.DeleteAdRequest{AdId: &ad.Id, AuthorId: &user.Id})
	suite.NoError(err)

	_, err = suite.Client.GetAd(suite.Context, &grpcPort.GetAdRequest{AdId: &ad.Id})
	suite.Error(err)

	suite.Equal(ErrAdNotFound.Error(), err.Error())
}

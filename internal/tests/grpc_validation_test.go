package tests

import (
	grpcPort "homework10/internal/ports/grpc"
	"strings"
)

func (suite *GRPCSuite) TestGRPCCreateUser_InvalidEmail() {
	_, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "abc"})
	suite.Error(err, "suite.Client.GetUser")

	suite.Equal(ErrInvalidEmail.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCUpdateUser_InvalidEmail() {
	res, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "MacMiller", Email: "swimming@circles.com"})
	suite.NoError(err)

	_, err = suite.Client.UpdateUser(suite.Context, &grpcPort.UpdateUserRequest{Id: &res.Id, Name: "MacMiller", Email: "good_am.ru"})
	suite.Error(err)
	suite.Equal(ErrInvalidEmail.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCCreateAd_EmptyTitle() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "", Text: "Hill Drive", UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCCreateAd_TooLongTitle() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	text := strings.Repeat("a", 101)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: text, Text: "Hill Drive", UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCCreateAd_TooLongText() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	text := strings.Repeat("a", 501)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: text, UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCUpdateAd_EmptyTitle() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user1.Id})
	suite.NoError(err)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{Title: "", Text: "Hill Drive", UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCUpdateAd_TooLongTitle() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user1.Id})
	suite.NoError(err)
	text := strings.Repeat("a", 101)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{Title: text, Text: "Hill Drive", UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCUpdateAd_EmptyText() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user1.Id})
	suite.NoError(err)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{Title: "Hill Drive", Text: "", UserId: &user1.Id})
	suite.Error(err)
}

func (suite *GRPCSuite) TestGRPCUpdateAd_TooLongText() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user1.Id})
	suite.NoError(err)
	text := strings.Repeat("a", 501)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{Title: "Hill Drive", Text: text, UserId: &user1.Id})
	suite.Error(err)
}

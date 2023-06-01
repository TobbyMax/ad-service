package tests

import (
	grpcPort "homework10/internal/ports/grpc"
)

func (suite *GRPCSuite) TestGRPCChangeStatusAdOfAnotherUser() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "GOMD", Text: "Role Modelz", UserId: &user1.Id})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &ad.Id, UserId: &user2.Id, Published: true})
	suite.Error(err)

	suite.Equal(ErrGRPCForbidden.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCUpdateAdOfAnotherUser() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Pimp", Text: "A Butterfly", UserId: &user2.Id})
	suite.NoError(err)

	_, err = suite.Client.UpdateAd(suite.Context, &grpcPort.UpdateAdRequest{AdId: &ad.Id, UserId: &user1.Id, Title: "Mr. Morale", Text: "The Big Steppers"})
	suite.Error(err)

	suite.Equal(ErrGRPCForbidden.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCCreateAd_ID() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	res, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Pimp", Text: "A Butterfly", UserId: &user2.Id})
	suite.NoError(err)
	suite.Equal(res.Id, int64(0))

	res, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Mr. Morale", Text: "The Big Steppers"})
	suite.NoError(err)
	suite.Equal(res.Id, int64(1))

	res, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "Cole World", Text: "Born Sinner"})
	suite.NoError(err)
	suite.Equal(res.Id, int64(2))
}

func (suite *GRPCSuite) TestGRPCDeleteAdOfAnotherUser() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{Title: "Forest", Text: "Hill Drive", UserId: &user1.Id})
	suite.NoError(err)

	_, err = suite.Client.DeleteAd(suite.Context, &grpcPort.DeleteAdRequest{AdId: &ad.Id, AuthorId: &user2.Id})
	suite.Error(err)

	suite.Equal(ErrGRPCForbidden.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCGetUser_NonExistentID() {
	_, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "MacMiller", Email: "blue_slide_park@hotmail.com"})
	suite.NoError(err)

	var id int64 = 1
	_, err = suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{Id: &id})
	suite.Error(err)

	suite.Equal(ErrUserNotFound.Error(), err.Error())
}

func (suite *GRPCSuite) TestGRPCGetUser_NoID() {
	_, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "MacMiller", Email: "blue_slide_park@hotmail.com"})
	suite.NoError(err)

	_, err = suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{})
	suite.Error(err)

	suite.Equal(ErrMissingArgument.Error(), err.Error())
}

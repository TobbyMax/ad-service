package tests

import (
	grpcPort "homework10/internal/ports/grpc"
	"time"
)

func (suite *GRPCSuite) TestGRPCListAds() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	publishedAd, err := suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{})

	suite.NoError(err)
	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, publishedAd.Id)
	suite.Equal(ads.List[0].Title, publishedAd.Title)
	suite.Equal(ads.List[0].Text, publishedAd.Text)
	suite.Equal(ads.List[0].AuthorId, publishedAd.AuthorId)
	suite.True(ads.List[0].Published)
}

func (suite *GRPCSuite) TestGRPCListAdsPublished() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	publishedAd, err := suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	published := true
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Published: &published})

	suite.NoError(err)
	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, publishedAd.Id)
	suite.Equal(ads.List[0].Title, publishedAd.Title)
	suite.Equal(ads.List[0].Text, publishedAd.Text)
	suite.Equal(ads.List[0].AuthorId, publishedAd.AuthorId)
	suite.True(ads.List[0].Published)
}

func (suite *GRPCSuite) TestGRPCListAdsNotPublished() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	_, err = suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	ad1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &ad1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	notPublishedAd, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	published := false
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Published: &published})

	suite.NoError(err)
	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, notPublishedAd.Id)
	suite.Equal(ads.List[0].Title, notPublishedAd.Title)
	suite.Equal(ads.List[0].Text, notPublishedAd.Text)
	suite.Equal(ads.List[0].AuthorId, notPublishedAd.AuthorId)
	suite.False(ads.List[0].Published)
}

func (suite *GRPCSuite) TestGRPCListAdsByUser() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{UserId: &user1.Id})

	suite.NoError(err)
	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, adByUser1.Id)
	suite.Equal(ads.List[0].Title, adByUser1.Title)
	suite.Equal(ads.List[0].Text, adByUser1.Text)
	suite.Equal(ads.List[0].AuthorId, adByUser1.AuthorId)
	suite.True(ads.List[0].Published)
}

func (suite *GRPCSuite) TestGRPCListAdsByDate() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	today := time.Now().UTC().Format(DateLayout)
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Date: &today})
	suite.NoError(err)

	suite.Len(ads.List, 3)
}

func (suite *GRPCSuite) TestGRPCListAdsByDate_Yesterday() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	yesterday := time.Now().UTC().Add(time.Duration(-24) * time.Hour).Format(DateLayout)
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Date: &yesterday})
	suite.NoError(err)

	suite.Len(ads.List, 0)
}

func (suite *GRPCSuite) TestGRPCListAdsByUserAndDate() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	today := time.Now().UTC().Format(DateLayout)
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Date: &today, UserId: &user2.Id})
	suite.NoError(err)

	suite.Len(ads.List, 2)
}

func (suite *GRPCSuite) TestGRPCListAdsByTitle() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	gomd, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &gomd.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Fire Squad", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	title := "GOMD"
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Title: &title})
	suite.NoError(err)

	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, gomd.Id)
	suite.Equal(ads.List[0].Title, gomd.Title)
	suite.Equal(ads.List[0].Text, gomd.Text)
	suite.Equal(ads.List[0].AuthorId, gomd.AuthorId)
	suite.True(ads.List[0].Published)
}

func (suite *GRPCSuite) TestGRPCListAdsByTitle_Multiple() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	title := "GOMD"
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{Title: &title})
	suite.NoError(err)

	suite.Len(ads.List, 2)
}

func (suite *GRPCSuite) TestGRPCListAdsByOptions() {
	user1, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "J.Cole", Email: "foresthill@drive.com"})
	suite.NoError(err)

	user2, err := suite.Client.CreateUser(suite.Context, &grpcPort.CreateUserRequest{Name: "Kendrick", Email: "section80@damn.com"})
	suite.NoError(err)

	adByUser1, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	target, err := suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser1.Id, UserId: &user1.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user1.Id, Title: "GOMD", Text: "Role Modelz"})
	suite.NoError(err)

	adByUser2, err := suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "GOMD", Text: "Cole World"})
	suite.NoError(err)

	_, err = suite.Client.ChangeAdStatus(suite.Context, &grpcPort.ChangeAdStatusRequest{AdId: &adByUser2.Id, UserId: &user2.Id, Published: true})
	suite.NoError(err)

	_, err = suite.Client.CreateAd(suite.Context, &grpcPort.CreateAdRequest{UserId: &user2.Id, Title: "Born Sinner", Text: "Cole World"})
	suite.NoError(err)

	today := time.Now().UTC().Format(DateLayout)
	published := true
	ads, err := suite.Client.ListAds(suite.Context, &grpcPort.ListAdRequest{UserId: &user1.Id, Date: &today, Title: &target.Title, Published: &published})
	suite.NoError(err)
	suite.Len(ads.List, 1)
	suite.Equal(ads.List[0].Id, target.Id)
	suite.Equal(ads.List[0].Title, target.Title)
	suite.Equal(ads.List[0].Text, target.Text)
	suite.Equal(ads.List[0].AuthorId, target.AuthorId)
	suite.True(ads.List[0].Published)
}

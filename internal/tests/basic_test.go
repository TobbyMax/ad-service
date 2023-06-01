package tests

func (suite *HTTPSuite) TestCreateAd() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)
	suite.Zero(response.Data.ID)
	suite.Equal(response.Data.Title, "hello")
	suite.Equal(response.Data.Text, "world")
	suite.Equal(response.Data.AuthorID, int64(uResponse.Data.ID))
	suite.False(response.Data.Published)
}

func (suite *HTTPSuite) TestChangeAdStatus() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)

	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	suite.NoError(err)
	suite.True(response.Data.Published)

	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	suite.NoError(err)
	suite.False(response.Data.Published)

	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	suite.NoError(err)
	suite.False(response.Data.Published)
}

func (suite *HTTPSuite) TestUpdateAd() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)

	response, err = suite.Client.updateAd(uResponse.Data.ID, response.Data.ID, "привет", "мир")
	suite.NoError(err)
	suite.Equal(response.Data.Title, "привет")
	suite.Equal(response.Data.Text, "мир")
}

func (suite *HTTPSuite) TestListAds() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)

	publishedAd, err := suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(uResponse.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAds()
	suite.NoError(err)
	suite.Len(ads.Data, 1)
	suite.Equal(ads.Data[0].ID, publishedAd.Data.ID)
	suite.Equal(ads.Data[0].Title, publishedAd.Data.Title)
	suite.Equal(ads.Data[0].Text, publishedAd.Data.Text)
	suite.Equal(ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	suite.True(ads.Data[0].Published)
}

func (suite *HTTPSuite) TestGetAd() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)

	response, err = suite.Client.getAd(response.Data.ID)
	suite.NoError(err)
	suite.Zero(response.Data.ID)
	suite.Equal(response.Data.Title, "hello")
	suite.Equal(response.Data.Text, "world")
	suite.Equal(response.Data.AuthorID, int64(uResponse.Data.ID))
	suite.False(response.Data.Published)
}

func (suite *HTTPSuite) TestDeleteAd() {
	user1, err := suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)

	ad1, err := suite.Client.createAd(user1.Data.ID, "Good News", "Dang!")
	suite.NoError(err)

	_, err = suite.Client.deleteAd(ad1.Data.ID, user1.Data.ID)
	suite.NoError(err)

	_, err = suite.Client.getAd(ad1.Data.ID)
	suite.Error(err)
	suite.ErrorIs(err, ErrNotFound)
}

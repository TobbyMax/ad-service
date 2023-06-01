package tests

func (suite *HTTPSuite) TestChangeStatusAdOfAnotherUser() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user1.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, resp.Data.ID, true)
	suite.ErrorIs(err, ErrForbidden)
}

func (suite *HTTPSuite) TestUpdateAdOfAnotherUser() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.updateAd(user1.Data.ID, resp.Data.ID, "title", "text")
	suite.ErrorIs(err, ErrForbidden)
}

func (suite *HTTPSuite) TestCreateAd_ID() {
	_, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(0))

	resp, err = suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(1))

	resp, err = suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(2))
}

func (suite *HTTPSuite) TestDeleteAdOfAnotherUser() {
	user1, err := suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Childish Gambino", "because@internet.com")
	suite.NoError(err)

	ad1, err := suite.Client.createAd(user1.Data.ID, "Good News", "Dang!")
	suite.NoError(err)

	_, err = suite.Client.deleteAd(ad1.Data.ID, user2.Data.ID)
	suite.Error(err)
	suite.ErrorIs(err, ErrForbidden)

	_, err = suite.Client.getAd(ad1.Data.ID)
	suite.NoError(err)
}

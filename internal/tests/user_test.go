package tests

func (suite *HTTPSuite) TestCreateUser() {
	response, err := suite.Client.createUser("TobbyMax", "agemax@gmail.com")
	suite.NoError(err)
	suite.Zero(response.Data.ID)
	suite.Equal("TobbyMax", response.Data.Nickname)
	suite.Equal("agemax@gmail.com", response.Data.Email)
}

func (suite *HTTPSuite) TestCreateUser_InvalidEmail() {
	_, err := suite.Client.createUser("TobbyMax", "abc")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestGetUser() {
	response, err := suite.Client.createUser("TobbyMax", "agemax@gmail.com")
	suite.NoError(err)

	response, err = suite.Client.getUser(response.Data.ID)
	suite.NoError(err)
	suite.Zero(response.Data.ID)
	suite.Equal("TobbyMax", response.Data.Nickname)
	suite.Equal("agemax@gmail.com", response.Data.Email)
}

func (suite *HTTPSuite) TestGetUser_NonExistentID() {
	_, err := suite.Client.createUser("MacMiller", "blue_slide_park@hotmail.com")
	suite.NoError(err)

	_, err = suite.Client.getUser(1)
	suite.ErrorIs(err, ErrNotFound)
}

func (suite *HTTPSuite) TestUpdateUser() {
	response, err := suite.Client.createUser("MacMiller", "swimming@circles.com")
	suite.NoError(err)

	response, err = suite.Client.updateUser(response.Data.ID, "MacMiller", "the_divine2016@feminine.ru")
	suite.NoError(err)
	suite.Equal("MacMiller", response.Data.Nickname)
	suite.Equal("the_divine2016@feminine.ru", response.Data.Email)
}

func (suite *HTTPSuite) TestUpdateUser_InvalidEmail() {
	response, err := suite.Client.createUser("MacMiller", "swimming@circles.com")
	suite.NoError(err)

	response, err = suite.Client.updateUser(response.Data.ID, "MacMiller", "good_am.ru")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestCreateUser_ID() {
	resp, err := suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(0))

	resp, err = suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(1))

	resp, err = suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)
	suite.Equal(resp.Data.ID, int64(2))
}

func (suite *HTTPSuite) TestDeleteUser() {
	user1, err := suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)

	ad1, err := suite.Client.createAd(user1.Data.ID, "Good News", "Dang!")
	suite.NoError(err)

	_, err = suite.Client.createUser("Mac Miller", "swimming@circles.com")
	suite.NoError(err)

	_, err = suite.Client.deleteUser(user1.Data.ID)
	suite.NoError(err)

	_, err = suite.Client.getUser(user1.Data.ID)
	suite.Error(err)
	suite.ErrorIs(ErrNotFound, err)

	_, err = suite.Client.getAd(ad1.Data.ID)
	suite.Error(err)
	suite.ErrorIs(err, ErrNotFound)
}

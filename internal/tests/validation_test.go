package tests

import "strings"

func (suite *HTTPSuite) TestCreateAd_EmptyTitle() {
	_, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "", "world")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestCreateAd_TooLongTitle() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	title := strings.Repeat("a", 101)

	_, err = suite.Client.createAd(user1.Data.ID, title, "world")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestCreateAd_EmptyText() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	_, err = suite.Client.createAd(user1.Data.ID, "title", "")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestCreateAd_TooLongText() {
	_, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	text := strings.Repeat("a", 501)

	_, err = suite.Client.createAd(user2.Data.ID, "title", text)
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestUpdateAd_EmptyTitle() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user1.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.updateAd(user1.Data.ID, resp.Data.ID, "", "new_world")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestUpdateAd_TooLongTitle() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user1.Data.ID, "hello", "world")
	suite.NoError(err)

	title := strings.Repeat("a", 101)

	_, err = suite.Client.updateAd(user1.Data.ID, resp.Data.ID, title, "world")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestUpdateAd_EmptyText() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	resp, err := suite.Client.createAd(user1.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.updateAd(user1.Data.ID, resp.Data.ID, "title", "")
	suite.ErrorIs(err, ErrBadRequest)
}

func (suite *HTTPSuite) TestUpdateAd_TooLongText() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	text := strings.Repeat("a", 501)

	resp, err := suite.Client.createAd(user1.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.updateAd(user1.Data.ID, resp.Data.ID, "title", text)
	suite.ErrorIs(err, ErrBadRequest)
}

package tests

import "time"

func (suite *HTTPSuite) TestCreateDate() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)
	suite.Zero(response.Data.ID)
	suite.Equal(response.Data.Title, "hello")
	suite.Equal(response.Data.Text, "world")
	suite.Equal(response.Data.AuthorID, uResponse.Data.ID)
	suite.False(response.Data.Published)

	suite.True(response.Data.DateCreated == response.Data.DateChanged)
	date, _ := time.Parse(DateTimeLayout, response.Data.DateCreated)
	suite.True(time.Since(date) < time.Hour)
}

func (suite *HTTPSuite) TestChangeDate() {
	uResponse, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(uResponse.Data.ID, "hello", "world")
	suite.NoError(err)

	time.Sleep(2 * time.Second)
	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, true)
	suite.NoError(err)
	suite.True(response.Data.Published)

	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	suite.NoError(err)
	suite.False(response.Data.Published)

	response, err = suite.Client.changeAdStatus(uResponse.Data.ID, response.Data.ID, false)
	suite.NoError(err)
	suite.False(response.Data.Published)

	suite.True(response.Data.DateCreated != response.Data.DateChanged)
	date, _ := time.Parse(DateTimeLayout, response.Data.DateChanged)
	suite.True(time.Since(date) < time.Hour)
}

package tests

import "time"

func (suite *HTTPSuite) TestListAdsPublished() {
	_, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	publishedAd, err := suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
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

func (suite *HTTPSuite) TestListAdsNotPublished() {
	_, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	notPublishedAd, err := suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByStatus(false)
	suite.NoError(err)
	suite.Len(ads.Data, 1)
	suite.Equal(ads.Data[0].ID, notPublishedAd.Data.ID)
	suite.Equal(ads.Data[0].Title, notPublishedAd.Data.Title)
	suite.Equal(ads.Data[0].Text, notPublishedAd.Data.Text)
	suite.Equal(ads.Data[0].AuthorID, notPublishedAd.Data.AuthorID)
	suite.False(ads.Data[0].Published)
}

func (suite *HTTPSuite) TestListAdsByUser() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	adByUser1, err := suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByUser(user1.Data.ID)
	suite.NoError(err)
	suite.Len(ads.Data, 1)
	suite.Equal(ads.Data[0].ID, adByUser1.Data.ID)
	suite.Equal(ads.Data[0].Title, adByUser1.Data.Title)
	suite.Equal(ads.Data[0].Text, adByUser1.Data.Text)
	suite.Equal(ads.Data[0].AuthorID, adByUser1.Data.AuthorID)
	suite.False(ads.Data[0].Published)
}

func (suite *HTTPSuite) TestListAdsByDate() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	_, err = suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByDate(time.Now().UTC().Format(DateLayout))
	suite.NoError(err)
	suite.Len(ads.Data, 3)
}

func (suite *HTTPSuite) TestListAdsByDate_Yesterday() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	_, err = suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByDate(time.Now().UTC().Add(time.Duration(-24) * time.Hour).Format(DateLayout))
	suite.NoError(err)
	suite.Len(ads.Data, 0)
}

func (suite *HTTPSuite) TestListAdsByUserAndDate() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	_, err = suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByUserAndDate(user2.Data.ID, time.Now().UTC().Format(DateLayout))
	suite.NoError(err)
	suite.Len(ads.Data, 2)
}

func (suite *HTTPSuite) TestListAdsByTitle() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	titledAd, err := suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "best cat", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByTitle(titledAd.Data.Title)
	suite.NoError(err)
	suite.Len(ads.Data, 1)
	suite.Equal(ads.Data[0].ID, titledAd.Data.ID)
	suite.Equal(ads.Data[0].Title, titledAd.Data.Title)
	suite.Equal(ads.Data[0].Text, titledAd.Data.Text)
	suite.Equal(ads.Data[0].AuthorID, titledAd.Data.AuthorID)
	suite.False(ads.Data[0].Published)
}

func (suite *HTTPSuite) TestListAdsByTitle_Multiple() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	titledAd, err := suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	response, err := suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user2.Data.ID, "GOMD", "not for sale")
	suite.NoError(err)

	ads, err := suite.Client.listAdsByTitle(titledAd.Data.Title)
	suite.NoError(err)
	suite.Len(ads.Data, 2)
}

func (suite *HTTPSuite) TestListAdsByOptions() {
	user1, err := suite.Client.createUser("J.Cole", "foresthill@drive.com")
	suite.NoError(err)

	user2, err := suite.Client.createUser("Kendrick", "section80@damn.com")
	suite.NoError(err)

	response, err := suite.Client.createAd(user1.Data.ID, "GOMD", "Role Modelz")
	suite.NoError(err)

	target, err := suite.Client.changeAdStatus(user1.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	_, err = suite.Client.createAd(user1.Data.ID, "GOMD", "Cole World")
	suite.NoError(err)

	response, err = suite.Client.createAd(user2.Data.ID, "hello", "world")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	response, err = suite.Client.createAd(user2.Data.ID, "GOMD", "not for sale")
	suite.NoError(err)

	_, err = suite.Client.changeAdStatus(user2.Data.ID, response.Data.ID, true)
	suite.NoError(err)

	ads, err := suite.Client.listAdsByOptions(user1.Data.ID, time.Now().UTC().Format(DateLayout), true, target.Data.Title)
	suite.NoError(err)
	suite.Len(ads.Data, 1)
	suite.Equal(ads.Data[0].ID, target.Data.ID)
	suite.Equal(ads.Data[0].Title, target.Data.Title)
	suite.Equal(ads.Data[0].Text, target.Data.Text)
	suite.Equal(ads.Data[0].AuthorID, target.Data.AuthorID)
	suite.True(ads.Data[0].Published)
}

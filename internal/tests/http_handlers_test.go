package tests

import (
	"github.com/TobbyMax/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"homework10/internal/ads"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin"
	"homework10/internal/tests/mocks"
	"homework10/internal/user"
	"log"
	"net/http/httptest"
	"testing"
)

type HTTPMockSuite struct {
	suite.Suite
	App    *mocks.App
	Client *testClient
}

func (suite *HTTPMockSuite) SetupTest() {
	log.Println("Setting Up Test")

	suite.App = mocks.NewApp(suite.T())
	server := httpgin.NewHTTPServer(":18080", suite.App)
	testServer := httptest.NewServer(server.Handler)

	suite.Client = &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (suite *HTTPMockSuite) TearDownTest() {
	log.Println("Tearing Down Test")
}

func (suite *HTTPMockSuite) TestHandler_CreateUser() {
	type args struct {
		badReq   bool
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful create",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "invalid email",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "validation error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      validator.ValidationErrors{},
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("CreateUser",
					mock.AnythingOfType("*gin.Context"),
					tc.args.nickname, tc.args.email,
				).
					Return(&user.User{Nickname: tc.args.nickname, Email: tc.args.email}, tc.args.err).
					Once()
			}
			var (
				response userResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.createUser(nil, nil)
			} else {
				response, err = suite.Client.createUser(tc.args.nickname, tc.args.email)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.nickname, response.Data.Nickname)
				suite.Equal(tc.args.email, response.Data.Email)
				suite.Equal(int64(0), response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_GetUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful get",
			args: args{
				id: 1,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("GetUser",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id,
				).
					Return(&user.User{ID: tc.args.id, Nickname: "Mac Miller", Email: "swimming@circles.com"}, tc.args.err).
					Once()
			}
			var (
				response userResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.getUser("hi")
			} else {
				response, err = suite.Client.getUser(tc.args.id)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal("Mac Miller", response.Data.Nickname)
				suite.Equal("swimming@circles.com", response.Data.Email)
				suite.Equal(tc.args.id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_UpdateUser() {
	type args struct {
		badId    bool
		badBody  bool
		id       int64
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful update",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "validation error",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      validator.ValidationErrors{},
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      app.ErrUserNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: id not int",
			args: args{
				badId:    true,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: unable to bind data",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				badBody:  true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "invalid email",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			id, name, email, e := tc.args.id, tc.args.nickname, tc.args.email, tc.args.err
			if tc.needMock {
				suite.App.On("UpdateUser",
					mock.AnythingOfType("*gin.Context"),
					id, name, email,
				).
					Return(&user.User{ID: id, Nickname: name, Email: email}, e).
					Once()
			}
			var (
				response userResponse
				err      error
			)
			if tc.args.badId {
				response, err = suite.Client.updateUser("hi", name, email)
			} else if tc.args.badBody {
				response, err = suite.Client.updateUser(id, 13, email)
			} else {
				response, err = suite.Client.updateUser(id, name, email)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(name, response.Data.Nickname)
				suite.Equal(email, response.Data.Email)
				suite.Equal(id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_DeleteUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful get",
			args: args{
				id: 1,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("DeleteUser",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id,
				).
					Return(tc.args.err).
					Once()
			}
			var err error
			if tc.args.badReq {
				_, err = suite.Client.deleteUser("hi")
			} else {
				_, err = suite.Client.deleteUser(tc.args.id)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_CreateAd() {
	type args struct {
		badReq bool
		title  string
		text   string
		uid    int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful create",
			args: args{
				title: "DAMN.",
				text:  "by Kendrick Lamar",
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "failed dependency",
			args: args{
				err: app.ErrUserNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrFailedDependency)
				return true
			},
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "validation error",
			args: args{
				title: "",
				text:  "by Kendrick Lamar",
				err:   validator.ValidationErrors{},
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				err:   ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("CreateAd",
					mock.AnythingOfType("*gin.Context"),
					tc.args.title, tc.args.text, tc.args.uid,
				).
					Return(&ads.Ad{Title: tc.args.title, Text: tc.args.text}, tc.args.err).
					Once()
			}
			var (
				response adResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.createAd(tc.args.uid, 1, tc.args.text)
			} else {
				response, err = suite.Client.createAd(tc.args.uid, tc.args.title, tc.args.text)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.title, response.Data.Title)
				suite.Equal(tc.args.text, response.Data.Text)
				suite.Equal(int64(0), response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_GetAd() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful get",
			args: args{
				id: 2009,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "id not found",
			args: args{
				err: app.ErrAdNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  2009,
				err: ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("GetAd",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id,
				).
					Return(&ads.Ad{ID: tc.args.id}, tc.args.err).
					Once()
			}
			var (
				response adResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.getAd("hi")
			} else {
				response, err = suite.Client.getAd(tc.args.id)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_UpdateAd() {
	type args struct {
		badId   bool
		badBody bool
		id      int64
		title   string
		text    string
		uid     int64
		err     error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful update",
			args: args{
				id:    2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "validation error",
			args: args{
				id:    2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
				err:   validator.ValidationErrors{},
			},
			wantErr:  true,
			needMock: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "id not found",
			args: args{
				id:    2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
				err:   app.ErrAdNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "forbidden",
			args: args{
				id:    2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
				err:   app.ErrForbidden,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrForbidden)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:    2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
				err:   ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: id not int",
			args: args{
				badId: true,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				uid:   13,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: unable to bind data",
			args: args{
				id:      2009,
				title:   "DAMN.",
				text:    "by Kendrick Lamar",
				uid:     13,
				badBody: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("UpdateAd",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id, tc.args.uid, tc.args.title, tc.args.text,
				).
					Return(&ads.Ad{ID: tc.args.id, AuthorID: tc.args.uid, Title: tc.args.title, Text: tc.args.text}, tc.args.err).
					Once()
			}
			var (
				response adResponse
				err      error
			)
			if tc.args.badId {
				response, err = suite.Client.updateAd(tc.args.uid, "oops", tc.args.title, tc.args.text)
			} else if tc.args.badBody {
				response, err = suite.Client.updateAd(tc.args.uid, tc.args.id, 13, tc.args.text)
			} else {
				response, err = suite.Client.updateAd(tc.args.uid, tc.args.id, tc.args.title, tc.args.text)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.title, response.Data.Title)
				suite.Equal(tc.args.text, response.Data.Text)
				suite.Equal(tc.args.uid, response.Data.AuthorID)
				suite.Equal(tc.args.id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_ChangeAdStatus() {
	type args struct {
		badId     bool
		badBody   bool
		id        int64
		published bool
		uid       int64
		err       error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful status change",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "id not found",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				err:       app.ErrAdNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "forbidden",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				err:       app.ErrForbidden,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrForbidden)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				err:       ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: id not int",
			args: args{
				badId:     true,
				published: true,
				uid:       13,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: unable to bind data",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				badBody:   true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("ChangeAdStatus",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id, tc.args.uid, tc.args.published,
				).
					Return(&ads.Ad{ID: tc.args.id, AuthorID: tc.args.uid, Published: tc.args.published}, tc.args.err).
					Once()
			}
			var (
				response adResponse
				err      error
			)
			if tc.args.badId {
				response, err = suite.Client.changeAdStatus(tc.args.uid, "oops", tc.args.published)
			} else if tc.args.badBody {
				response, err = suite.Client.changeAdStatus(tc.args.uid, tc.args.id, "pop")
			} else {
				response, err = suite.Client.changeAdStatus(tc.args.uid, tc.args.id, tc.args.published)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.published, response.Data.Published)
				suite.Equal(tc.args.id, response.Data.ID)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_DeleteAd() {
	type args struct {
		badID  bool
		badUID bool
		noUID  bool
		id     int64
		uid    int64
		err    error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "successful delete",
			args: args{
				id:  2009,
				uid: 13,
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "id not found",
			args: args{
				id:  2009,
				uid: 13,
				err: app.ErrAdNotFound,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrNotFound)
				return true
			},
		},
		{
			name: "forbidden",
			args: args{
				id:  2009,
				uid: 13,
				err: app.ErrForbidden,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrForbidden)
				return true
			},
		},
		{
			name: "internal error",
			args: args{
				id:  2009,
				uid: 13,
				err: ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: id not int",
			args: args{
				id:    2009,
				badID: true,
				uid:   13,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: uid not int",
			args: args{
				id:     2009,
				badUID: true,
				uid:    13,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: no userId",
			args: args{
				id:    2009,
				noUID: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("DeleteAd",
					mock.AnythingOfType("*gin.Context"),
					tc.args.id, tc.args.uid,
				).
					Return(tc.args.err).
					Once()
			}
			var err error
			if tc.args.badID {
				_, err = suite.Client.deleteAd("hi", tc.args.uid)
			} else if tc.args.badUID {
				_, err = suite.Client.deleteAd(tc.args.id, "pop")
			} else if tc.args.noUID {
				_, err = suite.Client.badDeleteAd(tc.args.id)
			} else {
				_, err = suite.Client.deleteAd(tc.args.id, tc.args.uid)
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *HTTPMockSuite) TestHandler_Filter() {
	type args struct {
		badBody bool
		badDate bool
		err     error
	}
	tests := []struct {
		name     string
		args     args
		needMock bool
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name:     "successful filter",
			args:     args{},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "internal error",
			args: args{
				err: ErrMock,
			},
			needMock: true,
			wantErr:  true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrInternal)
				return true
			},
		},
		{
			name: "bad request: unable to bind",
			args: args{
				badBody: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
		{
			name: "bad request: wrong date format",
			args: args{
				badDate: true,
			},
			wantErr: true,
			checkErr: func(err error) bool {
				suite.ErrorIs(err, ErrBadRequest)
				return true
			},
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("ListAds",
					mock.AnythingOfType("*gin.Context"),
					mock.AnythingOfType("app.ListAdsParams"),
				).
					Return(&ads.AdList{}, tc.args.err).
					Once()
			}
			var err error
			if tc.args.badBody {
				_, err = suite.Client.listAdsByStatus(1)
			} else if tc.args.badDate {
				_, err = suite.Client.listAdsByDate("01.05.2023")
			} else {
				_, err = suite.Client.listAdsByOptions(int64(777), "2018-08-03", true, "Swimming")
			}
			if tc.wantErr {
				suite.Error(err)
				suite.True(tc.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func TestHTTPHandlers(t *testing.T) {
	suite.Run(t, new(HTTPMockSuite))
}

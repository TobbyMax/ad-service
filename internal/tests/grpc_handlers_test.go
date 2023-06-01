package tests

import (
	"context"
	"github.com/TobbyMax/validator"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/ads"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
	"homework10/internal/tests/mocks"
	"homework10/internal/user"
	"log"
	"net"
	"testing"
	"time"
)

type GRPCMockSuite struct {
	suite.Suite
	App     *mocks.App
	Client  grpcPort.AdServiceClient
	Conn    *grpc.ClientConn
	Context context.Context
	Cancel  context.CancelFunc
	Server  *grpc.Server
	Lis     *bufconn.Listener
}

func (suite *GRPCMockSuite) SetupSuite() {
	log.Println("Setting Up Suite")

	suite.App = mocks.NewApp(suite.T())

	suite.Lis = bufconn.Listen(1024 * 1024)
	suite.Server = grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcPort.UnaryLoggerInterceptor,
		grpcPort.UnaryRecoveryInterceptor(),
	))

	svc := grpcPort.NewService(suite.App)
	grpcPort.RegisterAdServiceServer(suite.Server, svc)
	go func() {
		suite.NoError(suite.Server.Serve(suite.Lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return suite.Lis.Dial()
	}

	suite.Context, suite.Cancel = context.WithTimeout(context.Background(), 30*time.Second)

	conn, err := grpc.DialContext(suite.Context, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.NoError(err, "grpc.DialContext")
	suite.Conn = conn
}

func (suite *GRPCMockSuite) TearDownSuite() {
	log.Println("Tearing Down Suite")

	err := suite.Conn.Close()
	if err != nil {
		log.Println("Error closing connection")
	}
	suite.Cancel()
	suite.Server.Stop()
	err = suite.Lis.Close()
	if err != nil {
		log.Println("Error closing listener")
	}
}

func (suite *GRPCMockSuite) SetupTest() {
	log.Println("Setting Up Test")
	suite.Client = grpcPort.NewAdServiceClient(suite.Conn)
}

func (suite *GRPCMockSuite) TestHandler_CreateUser() {
	type args struct {
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			name: "invalid email",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr:       true,
			expectedError: ErrInvalidEmail,
		},
		{
			name: "validation error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      validator.ValidationErrors{},
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrValidationMock,
		},
		{
			name: "internal error",
			args: args{
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("CreateUser",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.nickname, tc.args.email,
				).
					Return(&user.User{Nickname: tc.args.nickname, Email: tc.args.email}, tc.args.err).
					Once()
			}
			response, err := suite.Client.CreateUser(suite.Context,
				&grpcPort.CreateUserRequest{Name: tc.args.nickname, Email: tc.args.email})
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.nickname, response.GetName())
				suite.Equal(tc.args.email, response.GetEmail())
				suite.Equal(int64(0), response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_GetUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("GetUser",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id,
				).
					Return(&user.User{ID: tc.args.id, Nickname: "Mac Miller", Email: "swimming@circles.com"}, tc.args.err).
					Once()
			}
			var (
				response *grpcPort.UserResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{})
			} else {
				response, err = suite.Client.GetUser(suite.Context, &grpcPort.GetUserRequest{Id: &tc.args.id})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal("Mac Miller", response.GetName())
				suite.Equal("swimming@circles.com", response.GetEmail())
				suite.Equal(tc.args.id, response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_UpdateUser() {
	type args struct {
		badReq   bool
		id       int64
		nickname string
		email    string
		err      error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrValidationMock,
		},
		{
			name: "id not found",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      app.ErrUserNotFound,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
				err:      ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
		{
			name: "bad request",
			args: args{
				badReq:   true,
				nickname: "Mac Miller",
				email:    "swimming@circles.com",
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "invalid email",
			args: args{
				id:       2009,
				nickname: "Mac Miller",
				email:    "swimming.com",
			},
			wantErr:       true,
			expectedError: ErrInvalidEmail,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			id, name, email, e := tc.args.id, tc.args.nickname, tc.args.email, tc.args.err
			if tc.needMock {
				suite.App.On("UpdateUser",
					mock.AnythingOfType("*context.valueCtx"),
					id, name, email,
				).
					Return(&user.User{ID: id, Nickname: name, Email: email}, e).
					Once()
			}
			var (
				response *grpcPort.UserResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.UpdateUser(suite.Context, &grpcPort.UpdateUserRequest{Name: name, Email: email})
			} else {
				response, err = suite.Client.UpdateUser(suite.Context,
					&grpcPort.UpdateUserRequest{Id: &id, Name: name, Email: email})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(name, response.GetName())
				suite.Equal(email, response.GetEmail())
				suite.Equal(id, response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_DeleteUser() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "id not found",
			args: args{
				id:  1,
				err: app.ErrUserNotFound,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrUserNotFound,
		},
		{
			name: "internal error",
			args: args{
				id:  1,
				err: ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("DeleteUser",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id,
				).
					Return(tc.args.err).
					Once()
			}
			var err error
			if tc.args.badReq {
				_, err = suite.Client.DeleteUser(suite.Context, &grpcPort.DeleteUserRequest{})
			} else {
				_, err = suite.Client.DeleteUser(suite.Context, &grpcPort.DeleteUserRequest{Id: &tc.args.id})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_CreateAd() {
	type args struct {
		title  string
		text   string
		uid    int64
		err    error
		badReq bool
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful create",
			args: args{
				uid:   2009,
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrUserNotFound,
		},
		{
			name: "validation error",
			args: args{
				uid:   2009,
				title: "",
				text:  "by Kendrick Lamar",
				err:   validator.ValidationErrors{},
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrValidationMock,
		},
		{
			name: "bad request: no userId",
			args: args{
				title:  "DAMN.",
				badReq: true,
				text:   "by Kendrick Lamar",
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "internal error",
			args: args{
				uid:   2009,
				title: "DAMN.",
				text:  "by Kendrick Lamar",
				err:   ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("CreateAd",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.title, tc.args.text, tc.args.uid,
				).
					Return(&ads.Ad{Title: tc.args.title, Text: tc.args.text}, tc.args.err).
					Once()
			}
			var (
				response *grpcPort.AdResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.CreateAd(suite.Context,
					&grpcPort.CreateAdRequest{Title: tc.args.title, Text: tc.args.text})
			} else {
				response, err = suite.Client.CreateAd(suite.Context,
					&grpcPort.CreateAdRequest{UserId: &tc.args.uid, Title: tc.args.title, Text: tc.args.text})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.title, response.GetTitle())
				suite.Equal(tc.args.text, response.GetText())
				suite.Equal(int64(0), response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_GetAd() {
	type args struct {
		badReq bool
		id     int64
		err    error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrAdNotFound,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "internal error",
			args: args{
				id:  2009,
				err: ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("GetAd",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id,
				).
					Return(&ads.Ad{ID: tc.args.id}, tc.args.err).
					Once()
			}
			var (
				response *grpcPort.AdResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.GetAd(suite.Context, &grpcPort.GetAdRequest{})
			} else {
				response, err = suite.Client.GetAd(suite.Context, &grpcPort.GetAdRequest{AdId: &tc.args.id})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.id, response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_UpdateAd() {
	type args struct {
		badReq bool
		id     int64
		title  string
		text   string
		uid    int64
		err    error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrValidationMock,
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrAdNotFound,
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrGRPCForbidden,
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
		{
			name: "bad request",
			args: args{
				badReq: true,
				title:  "DAMN.",
				text:   "by Kendrick Lamar",
				uid:    13,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("UpdateAd",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id, tc.args.uid, tc.args.title, tc.args.text,
				).
					Return(&ads.Ad{ID: tc.args.id, AuthorID: tc.args.uid, Title: tc.args.title, Text: tc.args.text}, tc.args.err).
					Once()
			}
			var (
				response *grpcPort.AdResponse
				err      error
			)
			if tc.args.badReq {
				response, err = suite.Client.UpdateAd(suite.Context,
					&grpcPort.UpdateAdRequest{AdId: &tc.args.id, Title: tc.args.title, Text: tc.args.text})
			} else {
				response, err = suite.Client.UpdateAd(suite.Context,
					&grpcPort.UpdateAdRequest{
						AdId:   &tc.args.id,
						UserId: &tc.args.uid,
						Title:  tc.args.title,
						Text:   tc.args.text,
					})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.title, response.GetTitle())
				suite.Equal(tc.args.text, response.GetText())
				suite.Equal(tc.args.uid, response.GetAuthorId())
				suite.Equal(tc.args.id, response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_ChangeAdStatus() {
	type args struct {
		badId     bool
		badUid    bool
		id        int64
		published bool
		uid       int64
		err       error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrAdNotFound,
		},
		{
			name: "forbidden",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				err:       app.ErrForbidden,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrGRPCForbidden,
		},
		{
			name: "internal error",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				err:       ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
		{
			name: "bad request: no adID",
			args: args{
				badId:     true,
				published: true,
				uid:       13,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "bad request: no userID",
			args: args{
				id:        2009,
				published: true,
				uid:       13,
				badUid:    true,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("ChangeAdStatus",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id, tc.args.uid, tc.args.published,
				).
					Return(&ads.Ad{ID: tc.args.id, AuthorID: tc.args.uid, Published: tc.args.published}, tc.args.err).
					Once()
			}
			var (
				response *grpcPort.AdResponse
				err      error
			)
			if tc.args.badId {
				response, err = suite.Client.ChangeAdStatus(suite.Context,
					&grpcPort.ChangeAdStatusRequest{UserId: &tc.args.uid, Published: tc.args.published})
			} else if tc.args.badUid {
				response, err = suite.Client.ChangeAdStatus(suite.Context,
					&grpcPort.ChangeAdStatusRequest{AdId: &tc.args.id, Published: tc.args.published})
			} else {
				response, err = suite.Client.ChangeAdStatus(suite.Context,
					&grpcPort.ChangeAdStatusRequest{AdId: &tc.args.id, UserId: &tc.args.uid, Published: tc.args.published})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
				suite.Equal(tc.args.published, response.GetPublished())
				suite.Equal(tc.args.id, response.GetId())
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_DeleteAd() {
	type args struct {
		badId  bool
		badUid bool
		id     int64
		uid    int64
		err    error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
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
			needMock:      true,
			wantErr:       true,
			expectedError: ErrAdNotFound,
		},
		{
			name: "forbidden",
			args: args{
				id:  2009,
				uid: 13,
				err: app.ErrForbidden,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrGRPCForbidden,
		},
		{
			name: "internal error",
			args: args{
				id:  2009,
				uid: 13,
				err: ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
		{
			name: "bad request: no adId",
			args: args{
				id:    2009,
				badId: true,
				uid:   13,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
		{
			name: "bad request: no userId",
			args: args{
				id:     2009,
				badUid: true,
				uid:    13,
			},
			wantErr:       true,
			expectedError: ErrMissingArgument,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("DeleteAd",
					mock.AnythingOfType("*context.valueCtx"),
					tc.args.id, tc.args.uid,
				).
					Return(tc.args.err).
					Once()
			}
			var err error
			if tc.args.badId {
				_, err = suite.Client.DeleteAd(suite.Context,
					&grpcPort.DeleteAdRequest{AuthorId: &tc.args.uid})
			} else if tc.args.badUid {
				_, err = suite.Client.DeleteAd(suite.Context,
					&grpcPort.DeleteAdRequest{AdId: &tc.args.id})
			} else {
				_, err = suite.Client.DeleteAd(suite.Context,
					&grpcPort.DeleteAdRequest{AdId: &tc.args.id, AuthorId: &tc.args.uid})
			}
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GRPCMockSuite) TestHandler_Filter() {
	type args struct {
		date string
		err  error
	}
	tests := []struct {
		name          string
		args          args
		needMock      bool
		wantErr       bool
		expectedError error
	}{
		{
			name: "successful filter",
			args: args{
				date: "2018-08-03",
			},
			needMock: true,
			wantErr:  false,
		},
		{
			name: "internal error",
			args: args{
				date: "2018-08-03",
				err:  ErrMock,
			},
			needMock:      true,
			wantErr:       true,
			expectedError: ErrMockInternal,
		},
		{
			name: "bad request: wrong date format",
			args: args{
				date: "abc",
			},
			wantErr:       true,
			expectedError: ErrDateMock,
		},
	}
	for _, tc := range tests {
		suite.Run(tc.name, func() {
			if tc.needMock {
				suite.App.On("ListAds",
					mock.AnythingOfType("*context.valueCtx"),
					mock.AnythingOfType("app.ListAdsParams"),
				).
					Return(&ads.AdList{}, tc.args.err).
					Once()
			}
			_, err := suite.Client.ListAds(suite.Context,
				&grpcPort.ListAdRequest{Date: &tc.args.date})
			if tc.wantErr {
				suite.Error(err)
				suite.Equal(tc.expectedError.Error(), err.Error(), "test expect an error, but got wrong error type")
			} else {
				suite.NoError(err)
			}
		})
	}
}

func (suite *GRPCMockSuite) TestPanic() {
	*suite.App = mocks.App{}
	suite.NotPanics(func() {
		_, err := suite.Client.ListAds(suite.Context,
			&grpcPort.ListAdRequest{})
		if err != nil {
			log.Printf("Function returned error: %v", err)
		}
	})
}

func TestGRPCHandlers(t *testing.T) {
	suite.Run(t, new(GRPCMockSuite))
}

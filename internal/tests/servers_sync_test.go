package tests

import (
	"context"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	"homework10/internal/graceful"
	grpcSvc "homework10/internal/ports/grpc"
	"homework10/internal/ports/httpgin"
	"homework10/internal/tests/mocks"
	"log"
	"net"
	"net/http/httptest"
	"os"
	"syscall"
	"testing"
	"time"
)

type ServerSuite struct {
	suite.Suite
	App        *mocks.App
	ClientHTTP *testClient
	ClientGRPC grpcSvc.AdServiceClient
	Lis        *bufconn.Listener
	SigQuit    chan os.Signal

	CtxClient    context.Context
	CancelClient context.CancelFunc
	ConnClient   *grpc.ClientConn
}

func (suite *ServerSuite) DialGRPC() {
	dialer := func(context.Context, string) (net.Conn, error) {
		return suite.Lis.Dial()
	}

	suite.CtxClient, suite.CancelClient = context.WithTimeout(context.Background(), 30*time.Second)

	conn, err := grpc.DialContext(suite.CtxClient, "", grpc.WithContextDialer(dialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.NoError(err, "grpc.DialContext")
	suite.ConnClient = conn

	suite.ClientGRPC = grpcSvc.NewAdServiceClient(suite.ConnClient)
}

func (suite *ServerSuite) SetupTest() {
	log.Println("Setting Up Test")
	appSvc := app.NewApp(adrepo.New())
	suite.Lis = bufconn.Listen(1024 * 1024)
	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterAdServiceServer(grpcServer, svc)

	httpServer := httpgin.NewHTTPServer(":18080", appSvc)
	suite.SigQuit = make(chan os.Signal, 1)

	eg, ctx := errgroup.WithContext(context.Background())

	eg.Go(graceful.CaptureSignal(ctx, suite.SigQuit))
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, suite.Lis, grpcServer))
	eg.Go(httpgin.RunHTTPServerGracefully(ctx, httpServer))
	go func() {
		if err := eg.Wait(); err != nil {
			log.Printf("gracefully shutting down the servers: %s\n", err.Error())
		}
		log.Println("servers were successfully shutdown")
	}()

	testServer := httptest.NewServer(httpServer.Handler)

	suite.ClientHTTP = &testClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
	suite.DialGRPC()
}

func (suite *ServerSuite) TestUser() {
	for id := 0; id < 100; id++ {
		name := "Mac Miller"
		email := "swimming@circles.com"
		if id%2 == 0 {
			resp, err := suite.ClientHTTP.createUser(name, email)

			suite.NoError(err)
			suite.Equal(int64(id), resp.Data.ID)
			suite.Equal(name, resp.Data.Nickname)
			suite.Equal(email, resp.Data.Email)
		} else {
			resp, err := suite.ClientGRPC.CreateUser(suite.CtxClient,
				&grpcSvc.CreateUserRequest{Name: name, Email: email})
			suite.NoError(err)
			suite.Equal(name, resp.GetName())
			suite.Equal(email, resp.GetEmail())
			suite.Equal(int64(id), resp.GetId())
		}
	}
}

func (suite *ServerSuite) TestAds() {
	for id := 0; id < 100; id++ {
		name := "Mac Miller"
		email := "swimming@circles.com"

		title := "Spins"
		text := "K.I.D.S."
		if id%2 == 0 {
			u, err := suite.ClientHTTP.createUser(name, email)
			suite.NoError(err)

			ad, err := suite.ClientHTTP.createAd(u.Data.ID, title, text)
			suite.NoError(err)
			suite.Equal(int64(id), ad.Data.ID)
			suite.Equal(title, ad.Data.Title)
			suite.Equal(text, ad.Data.Text)
		} else {
			u, err := suite.ClientGRPC.CreateUser(suite.CtxClient,
				&grpcSvc.CreateUserRequest{Name: name, Email: email})
			suite.NoError(err)
			ad, err := suite.ClientGRPC.CreateAd(suite.CtxClient,
				&grpcSvc.CreateAdRequest{Title: title, Text: text, UserId: &u.Id})
			suite.NoError(err)
			suite.Equal(title, ad.GetTitle())
			suite.Equal(text, ad.GetText())
			suite.Equal(int64(id), ad.GetId())
		}
	}
}

func (suite *ServerSuite) TearDownTest() {
	log.Println("Tearing Down Test")
	err := suite.ConnClient.Close()
	if err != nil {
		log.Println("Error closing connection")
	}
	suite.CancelClient()
	suite.SigQuit <- syscall.SIGINT
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerSuite))
}

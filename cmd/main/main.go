package main

import (
	"context"
	"github.com/TobbyMax/ad-service.git/internal/adapters/adrepo"
	"github.com/TobbyMax/ad-service.git/internal/app"
	"github.com/TobbyMax/ad-service.git/internal/graceful"
	grpcSvc "github.com/TobbyMax/ad-service.git/internal/ports/grpc"
	"github.com/TobbyMax/ad-service.git/internal/ports/httpgin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"os"

	"log"
	"net"
)

const (
	grpcPort = ":8080"
	httpPort = ":18080"
)

func CreateDB(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, os.Getenv("DB_CONNECT_STRING"))
}

func main() {
	appSvc := app.NewApp(adrepo.New())

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := grpcSvc.NewService(appSvc)
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcSvc.UnaryLoggerInterceptor,
		grpcSvc.UnaryRecoveryInterceptor(),
	))
	grpcSvc.RegisterAdServiceServer(grpcServer, svc)

	httpServer := httpgin.NewHTTPServer(httpPort, appSvc)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	eg.Go(graceful.CaptureSignal(ctx, sigQuit))
	// run grpc server
	eg.Go(grpcSvc.RunGRPCServerGracefully(ctx, lis, grpcServer))
	// run http server
	eg.Go(httpgin.RunHTTPServerGracefully(ctx, httpServer))

	if err := eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}
	log.Println("servers were successfully shutdown")
}

package grpc

import (
	"context"
	"fmt"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework10/internal/app"
	"log"
	"net"
	"runtime/debug"
	"time"
)

type AdService struct {
	app app.App
}

func NewService(a app.App) AdServiceServer {
	service := &AdService{app: a}
	return service
}

func UnaryLoggerInterceptor(ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	start := time.Now()
	log.Printf("-- received request -- | protocol: GRPC | method: %s", info.FullMethod)

	h, err := handler(ctx, req)

	latency := time.Since(start)
	log.Printf("-- handled request -- | protocol: GRPC | latency: %+v | method: %s | error: (%v)\n",
		latency, info.FullMethod, err)

	return h, err
}

func UnaryRecoveryInterceptor() grpc.UnaryServerInterceptor {
	stackTraceLogger := grpcRecovery.WithRecoveryHandlerContext(
		func(ctx context.Context, p interface{}) error {
			fmt.Print("\n\n")
			log.Printf("[PANIC] %s\n%s\n", p, string(debug.Stack()))
			return status.Errorf(codes.Internal, "%s", p)
		},
	)
	return grpcRecovery.UnaryServerInterceptor(stackTraceLogger)
}

func RunGRPCServerGracefully(ctx context.Context, lis net.Listener, server *grpc.Server) func() error {
	return func() error {
		log.Printf("starting grpc server, listening on %s\n", lis.Addr())
		defer log.Printf("close grpc server listening on %s\n", lis.Addr())

		errCh := make(chan error)

		defer func() {
			server.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err := server.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	}
}

package httpgin

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"homework10/internal/app"
)

func LoggerMiddleWare(c *gin.Context) {
	start := time.Now()

	log.Printf("-- received request -- | protocol: HTTP | method: %s | path: %s\n", c.Request.Method, c.Request.URL.Path)

	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()

	log.Printf("-- handled request -- | protocol: HTTP | status: %d | latency: %+v | method: %s | path: %s\n", status, latency, c.Request.Method, c.Request.URL.Path)
}

func NewHTTPServer(port string, a app.App) *http.Server {
	gin.SetMode(gin.ReleaseMode)
	handler := gin.New()
	s := &http.Server{Addr: port, Handler: handler}

	// todo: add your own logic

	api := handler.Group("/api/v1")

	// MiddleWare для логирования и паник
	api.Use(gin.Logger())
	api.Use(gin.Recovery())

	api.Use(LoggerMiddleWare)

	AppRouter(api, a)
	return s
}

func RunHTTPServerGracefully(ctx context.Context, server *http.Server) func() error {
	return func() error {
		log.Printf("starting http server, listening on %s\n", server.Addr)
		defer log.Printf("close http server listening on %s\n", server.Addr)

		errCh := make(chan error)

		defer func() {
			shCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := server.Shutdown(shCtx); err != nil {
				log.Printf("can't close http server listening on %s: %s", server.Addr, err.Error())
			}

			close(errCh)
		}()

		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return fmt.Errorf("http server can't listen and serve requests: %w", err)
		}
	}
}

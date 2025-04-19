package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"resource-service/internal/di"
	"time"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	shutdownTimeout time.Duration
}

func NewServer(ctx context.Context, host, port string, shutdownTimeout time.Duration) (context.Context, Server) {
	srv := Server{
		engine:          gin.New(),
		httpAddr:        fmt.Sprintf("%s:%s", host, port),
		shutdownTimeout: shutdownTimeout,
	}
	return serverContext(ctx), srv
}

func (s *Server) Run(ctx context.Context, container *di.Container) error {
	log.Println("server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}
	registerRoutes(s, container)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

// Gracefully shutdown
func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()
	return ctx
}

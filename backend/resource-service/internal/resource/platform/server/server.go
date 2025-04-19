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
	"syscall"
	"time"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	shutdownTimeout time.Duration
}

func NewServer(ctx context.Context, host, port string, shutdownTimeout time.Duration) (context.Context,
	context.CancelFunc, Server) {
	srv := Server{
		engine:          gin.New(),
		httpAddr:        fmt.Sprintf("%s:%s", host, port),
		shutdownTimeout: shutdownTimeout,
	}
	srv.engine.Use(gin.Logger(), gin.Recovery())
	c, cancel := serverContext(ctx)
	return c, cancel, srv
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
			log.Printf("server error: %v", err) // Cambiado de Fatal a Printf
		}
	}()

	<-ctx.Done()
	log.Println("server received shutdown signal")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	log.Println("server stopped gracefully")
	return nil
}

// Gracefully shutdown (mejorada)
func serverContext(ctx context.Context) (context.Context, context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Añadido SIGTERM

	ctx, cancel := context.WithCancel(ctx)

	// Almacenamos la función cancel en el contexto para poder usarla desde fuera
	ctx = context.WithValue(ctx, "cancelFunc", cancel)

	go func() {
		<-c
		log.Println("received interrupt signal")
		cancel()
	}()

	return ctx, cancel
}

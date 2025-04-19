package server

import (
	"github.com/gin-gonic/gin"
	"resource-service/internal/di"
)

func registerRoutes(s *Server, container *di.Container) {
	s.engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})
}

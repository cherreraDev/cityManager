package server

import (
	"github.com/gin-gonic/gin"
	"resource-service/internal/di"
	"resource-service/internal/resource/platform/server/handler"
)

func registerRoutes(s *Server, container *di.Container) {
	s.engine.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	resourceHandler := handler.NewResourceHandler(container.ResourceService)
	s.engine.GET("/cities/:city_id/resources", resourceHandler.GetResources)
	s.engine.POST("/cities/:city_id/resources/consume", resourceHandler.ConsumeResources)
}

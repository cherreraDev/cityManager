package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"resource-service/internal/resource/domain"
)

type ResourceHandler struct {
	service domain.ResourceService
}

func NewResourceHandler(service domain.ResourceService) *ResourceHandler {
	return &ResourceHandler{service: service}
}

type ConsumeResourcesDTO struct {
	Resources map[string]float64 `json:"resources"`
}

func (rh *ResourceHandler) GetResources(ctx *gin.Context) {
	cityIDStr := ctx.Param("city_id")
	cityID, err := uuid.Parse(cityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid city_id"})
		return
	}

	resources, err := rh.service.GetResources(cityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, resources)
}

func (rh *ResourceHandler) ConsumeResources(ctx *gin.Context) {
	cityIDStr := ctx.Param("city_id")
	cityID, err := uuid.Parse(cityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid city_id"})
		return
	}

	var dto ConsumeResourcesDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	if err := rh.service.ConsumeResources(context.Background(), cityID, dto.Resources); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "resources consumed successfully"})
}

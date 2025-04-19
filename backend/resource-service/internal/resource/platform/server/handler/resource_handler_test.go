package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"resource-service/internal/resource/domain"
	"resource-service/internal/resource/platform/server/handler"
	"resource-service/test/mocks"
	"testing"
)

func setupRouter(handler *handler.ResourceHandler) *gin.Engine {
	r := gin.Default()
	r.GET("/cities/:city_id/resources", handler.GetResources)
	r.POST("/cities/:city_id/resources/consume", handler.ConsumeResources)
	return r
}

func TestGetResources(t *testing.T) {
	mockService := &mocks.MockResourceService{
		GetResourcesFn: func(cityID uuid.UUID) ([]domain.Resource, error) {
			return []domain.Resource{
				{ID: uuid.New(), CityID: cityID, Type: "wood", Amount: 100},
			}, nil
		},
	}
	controller := handler.NewResourceHandler(mockService)
	router := setupRouter(controller)

	cityID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/cities/"+cityID.String()+"/resources", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "wood")
}

func TestConsumeResources(t *testing.T) {
	mockService := &mocks.MockResourceService{
		ConsumeResourcesFn: func(cityID uuid.UUID, resources map[string]float64) error {
			assert.Equal(t, float64(10), resources["wood"])
			return nil
		},
	}
	controller := handler.NewResourceHandler(mockService)
	router := setupRouter(controller)

	cityID := uuid.New()
	body := map[string]interface{}{
		"resources": map[string]float64{
			"wood": 10,
		},
	}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/cities/"+cityID.String()+"/resources/consume", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "resources consumed")
}

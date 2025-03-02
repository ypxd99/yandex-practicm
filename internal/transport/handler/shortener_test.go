package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/mocks"
	"github.com/ypxd99/yandex-practicm/internal/transport/handler"
	"github.com/ypxd99/yandex-practicm/util"
)

func setupRouter(service *mocks.MockLinkService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	h := handler.InitHandler(service)
	h.InitRoutes(r)
	return r
}

func TestShorterLinkHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	t.Run("successful request", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		url := "https://yandex.ru"
		mockService.On("ShorterLink", mock.Anything, url).
			Return("abc123", nil).
			Once()

		req := httptest.NewRequest("POST", "/", bytes.NewBufferString(url))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)
		assert.Equal(t, "abc123", resp.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		url := "https://error.com"
		mockService.On("ShorterLink", mock.Anything, url).
			Return("", errors.New("service error")).
			Once()

		req := httptest.NewRequest("POST", "/", bytes.NewBufferString(url))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetLinkByIDHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	t.Run("successful redirect", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		id := "abc123"
		target := "https://yandex.ru"
		mockService.On("FindLink", mock.Anything, id).
			Return(target, nil).
			Once()

		req := httptest.NewRequest("GET", "/"+id, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
		assert.Equal(t, target, resp.Header().Get("Location"))
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		id := "invalid"
		mockService.On("FindLink", mock.Anything, id).
			Return("", errors.New("not found")).
			Once()

		req := httptest.NewRequest("GET", "/"+id, nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

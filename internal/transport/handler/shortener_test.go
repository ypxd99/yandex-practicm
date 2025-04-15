package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/mocks"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/transport/handler"
	"github.com/ypxd99/yandex-practicm/util"
)

func setupRouter(service *mocks.MockLinkService) *gin.Engine {
	cfg := util.GetConfig()
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		testUserID := uuid.New()
		c.Set(cfg.Auth.CookieName, testUserID)
		c.Next()
	})

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
		mockService.On("ShorterLink", mock.Anything, url, mock.AnythingOfType("uuid.UUID")).
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
		mockService.On("ShorterLink", mock.Anything, url, mock.AnythingOfType("uuid.UUID")).
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

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestShortenHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)

	t.Run("successful request", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		input := model.ShortenRequest{
			URL: "https://yandex.ru",
		}
		mockService.On("ShorterLink", mock.Anything, input.URL, mock.AnythingOfType("uuid.UUID")).
			Return("abc123", nil).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response model.ShortenResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "abc123", response.Result)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid URL", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		input := struct {
			URL string `json:"url"`
		}{
			URL: "",
		}

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		input := model.ShortenRequest{
			URL: "https://yandex.ru",
		}
		mockService.On("ShorterLink", mock.Anything, input.URL, mock.AnythingOfType("uuid.UUID")).
			Return("", errors.New("service error")).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusInternalServerError, resp.Code)
		mockService.AssertExpectations(t)
	})
}

func TestBatchShortenHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)

	t.Run("successful batch", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		input := []model.BatchRequest{
			{CorrelationID: "1", OriginalURL: "https://example.com"},
			{CorrelationID: "2", OriginalURL: "https://yandex.ru"},
		}

		output := []model.BatchResponse{
			{CorrelationID: "1", ShortURL: "http://localhost:8080/abc123"},
			{CorrelationID: "2", ShortURL: "http://localhost:8080/def456"},
		}

		mockService.On("BatchShorten", mock.Anything, input, mock.AnythingOfType("uuid.UUID")).
			Return(output, nil).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/api/shorten/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusCreated, resp.Code)

		var response []model.BatchResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		mockService.AssertExpectations(t)
	})
}

func TestGetUserURLsHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)

	t.Run("successful get user urls", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		output := []model.UserURLResponse{
			{ShortURL: "http://localhost:8080/abc123", OriginalURL: "https://example.com"},
			{ShortURL: "http://localhost:8080/def456", OriginalURL: "https://yandex.ru"},
		}

		mockService.On("GetUserURLs", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(output, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)

		var response []model.UserURLResponse
		err := json.Unmarshal(resp.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("no content", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		var emptyOutput []model.UserURLResponse

		mockService.On("GetUserURLs", mock.Anything, mock.AnythingOfType("uuid.UUID")).
			Return(emptyOutput, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusNoContent, resp.Code)

		mockService.AssertExpectations(t)
	})
}

func TestDeleteURLsHandler(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)

	t.Run("successful delete", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		input := []string{"abc123", "def456"}

		mockService.On("DeleteURLs", mock.Anything, input, mock.AnythingOfType("uuid.UUID")).
			Return(2, nil).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusAccepted, resp.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("empty url list", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		body, _ := json.Marshal([]string{})
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockService := new(mocks.MockLinkService)
		router := setupRouter(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")

		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusBadRequest, resp.Code)
	})
}

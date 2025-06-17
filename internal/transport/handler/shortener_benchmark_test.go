package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/mocks"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
)

func setupTestRouter() (*gin.Engine, *mocks.MockLinkRepository) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockRepo := &mocks.MockLinkRepository{}
	svc := service.InitService(mockRepo)
	h := InitHandler(svc)

	h.InitRoutes(router)
	return router, mockRepo
}

func BenchmarkShortenLink(b *testing.B) {
	router, mockRepo := setupTestRouter()
	userID := uuid.New()

	req := model.ShortenRequest{
		URL: "https://example.com",
	}
	body, _ := json.Marshal(req)

	// Настраиваем мок
	mockRepo.On("CreateLink", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&model.Link{
		ID:     "test123",
		Link:   req.URL,
		UserID: userID,
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/shorten", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userID.String())

		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetLink(b *testing.B) {
	router, mockRepo := setupTestRouter()

	// Настраиваем мок
	mockRepo.On("FindLink", mock.Anything, "id123").Return(&model.Link{
		ID:     "id123",
		Link:   "https://example.com",
		UserID: uuid.New(),
	}, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/id123", nil)

		router.ServeHTTP(w, req)
	}
}

func BenchmarkBatchShorten(b *testing.B) {
	router, mockRepo := setupTestRouter()
	userID := uuid.New()

	batch := []model.BatchRequest{
		{CorrelationID: "1", OriginalURL: "https://example1.com"},
		{CorrelationID: "2", OriginalURL: "https://example2.com"},
		{CorrelationID: "3", OriginalURL: "https://example3.com"},
	}
	body, _ := json.Marshal(batch)

	// Настраиваем мок
	mockRepo.On("BatchCreate", mock.Anything, mock.Anything).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/shorten/batch", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+userID.String())

		router.ServeHTTP(w, req)
	}
}

func BenchmarkGetUserURLs(b *testing.B) {
	router, mockRepo := setupTestRouter()
	mockRepo.On("FindUserLinks", mock.Anything, mock.Anything).Return([]model.Link{
		{ID: "test1", Link: "http://test1.com", UserID: uuid.New(), IsDeleted: false},
		{ID: "test2", Link: "http://test2.com", UserID: uuid.New(), IsDeleted: false},
	}, nil)

	req := httptest.NewRequest("GET", "/api/user/urls", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "user_id", Value: uuid.New().String()})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

func BenchmarkDeleteURLs(b *testing.B) {
	router, mockRepo := setupTestRouter()
	mockRepo.On("MarkDeletedURLs", mock.Anything, mock.Anything, mock.Anything).Return(0, nil)

	deleteReq := []string{"id1", "id2", "id3"}
	body, _ := json.Marshal(deleteReq)
	req := httptest.NewRequest("DELETE", "/api/user/urls", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "user_id", Value: uuid.New().String()})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

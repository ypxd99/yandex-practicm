package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
)

// MockLinkService представляет мок-сервис для тестирования работы с сокращенными ссылками.
// Реализует интерфейс service.LinkService.
type MockLinkService struct {
	mock.Mock
}

// ShorterLink создает сокращенную ссылку для указанного URL.
// Принимает контекст, оригинальный URL и ID пользователя.
// Возвращает сокращенную ссылку или ошибку.
func (m *MockLinkService) ShorterLink(ctx context.Context, url string, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, url, userID)
	return args.String(0), args.Error(1)
}

// FindLink ищет оригинальный URL по сокращенному идентификатору.
// Принимает контекст и идентификатор сокращенной ссылки.
// Возвращает оригинальный URL или ошибку.
func (m *MockLinkService) FindLink(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

// StorageStatus проверяет доступность хранилища.
// Принимает контекст.
// Возвращает статус доступности и ошибку.
func (m *MockLinkService) StorageStatus(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

// BatchShorten создает сокращенные ссылки для пакета URL.
// Принимает контекст, список запросов на сокращение и ID пользователя.
// Возвращает список сокращенных ссылок или ошибку.
func (m *MockLinkService) BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error) {
	args := m.Called(ctx, batch, userID)
	return args.Get(0).([]model.BatchResponse), args.Error(1)
}

// GetUserURLs возвращает все сокращенные ссылки пользователя.
// Принимает контекст и ID пользователя.
// Возвращает список сокращенных ссылок или ошибку.
func (m *MockLinkService) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.UserURLResponse), args.Error(1)
}

// DeleteURLs помечает указанные сокращенные ссылки как удаленные.
// Принимает контекст, список идентификаторов ссылок и ID пользователя.
// Возвращает количество удаленных ссылок и ошибку.
func (m *MockLinkService) DeleteURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, ids, userID)
	return args.Int(0), args.Error(1)
}

var _ service.LinkService = (*MockLinkService)(nil)

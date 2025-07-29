package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

// MockLinkRepository представляет мок-репозиторий для тестирования работы с ссылками.
// Реализует интерфейс repository.LinkRepository.
type MockLinkRepository struct {
	mock.Mock
}

// CreateLink создает новую запись о сокращенной ссылке.
// Принимает контекст, идентификатор ссылки, оригинальный URL и ID пользователя.
// Возвращает созданную запись или ошибку.
func (m *MockLinkRepository) CreateLink(ctx context.Context, id, url string, userID uuid.UUID) (*model.Link, error) {
	args := m.Called(ctx, id, url, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

// FindLink ищет запись о сокращенной ссылке по её идентификатору.
// Принимает контекст и идентификатор ссылки.
// Возвращает найденную запись или ошибку.
func (m *MockLinkRepository) FindLink(ctx context.Context, id string) (*model.Link, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

// FindUserLinks возвращает все ссылки, созданные указанным пользователем.
// Принимает контекст и ID пользователя.
// Возвращает список ссылок или ошибку.
func (m *MockLinkRepository) FindUserLinks(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Link), args.Error(1)
}

// BatchCreate создает несколько записей о сокращенных ссылках.
// Принимает контекст и список ссылок для создания.
// Возвращает ошибку в случае неудачи.
func (m *MockLinkRepository) BatchCreate(ctx context.Context, links []model.Link) error {
	args := m.Called(ctx, links)
	return args.Error(0)
}

// MarkDeletedURLs помечает указанные ссылки как удаленные.
// Принимает контекст, список идентификаторов ссылок и ID пользователя.
// Возвращает количество помеченных ссылок и ошибку.
func (m *MockLinkRepository) MarkDeletedURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error) {
	args := m.Called(ctx, ids, userID)
	return args.Int(0), args.Error(1)
}

// Close закрывает соединение с хранилищем.
// Возвращает ошибку в случае неудачи.
func (m *MockLinkRepository) Close() error {
	return nil
}

// Status проверяет доступность хранилища.
// Принимает контекст.
// Возвращает статус доступности и ошибку.
func (m *MockLinkRepository) Status(ctx context.Context) (bool, error) {
	return true, nil
}

// GetStats возвращает статистику сервиса.
// Принимает контекст.
// Возвращает количество URL, количество пользователей и ошибку.
func (m *MockLinkRepository) GetStats(ctx context.Context) (int64, int64, error) {
	args := m.Called(ctx)
	return args.Get(0).(int64), args.Get(1).(int64), args.Error(2)
}

var _ repository.LinkRepository = (*MockLinkRepository)(nil)

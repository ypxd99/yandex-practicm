package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

// Service представляет сервисный слой для работы с сокращенными URL.
// Обеспечивает бизнес-логику для операций с URL и взаимодействует с репозиторием.
type Service struct {
	repo repository.LinkRepository
}

// LinkService определяет интерфейс для работы с сокращенными URL.
// Предоставляет методы для создания, поиска и управления URL.
type LinkService interface {
	// ShorterLink создает сокращенный URL для заданного длинного URL.
	// Возвращает сокращенный URL и ошибку, если операция не удалась.
	ShorterLink(ctx context.Context, url string, userID uuid.UUID) (string, error)

	// FindLink находит оригинальный URL по его сокращенному идентификатору.
	// Возвращает оригинальный URL и ошибку, если URL не найден.
	FindLink(ctx context.Context, id string) (string, error)

	// StorageStatus проверяет доступность хранилища.
	// Возвращает true, если хранилище доступно, и ошибку в противном случае.
	StorageStatus(ctx context.Context) (bool, error)

	// BatchShorten создает сокращенные URL для пакета длинных URL.
	// Возвращает массив сокращенных URL и ошибку, если операция не удалась.
	BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error)

	// GetUserURLs возвращает список всех URL, созданных пользователем.
	// Возвращает массив URL и ошибку, если операция не удалась.
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error)

	// DeleteURLs помечает указанные URL как удаленные.
	// Возвращает количество удаленных URL и ошибку, если операция не удалась.
	DeleteURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error)
}

// InitService создает и возвращает новый экземпляр Service с предоставленным репозиторием.
// Инициализирует сервис с заданной реализацией LinkRepository.
func InitService(repo repository.LinkRepository) *Service {
	return &Service{repo: repo}
}

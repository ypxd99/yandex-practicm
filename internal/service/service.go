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
	// Принимает контекст, оригинальный URL и идентификатор пользователя.
	// Возвращает сокращенный URL и ошибку, если операция не удалась.
	ShorterLink(ctx context.Context, url string, userID uuid.UUID) (string, error)

	// FindLink находит оригинальный URL по его сокращенному идентификатору.
	// Принимает контекст и идентификатор сокращенного URL.
	// Возвращает оригинальный URL и ошибку, если URL не найден.
	FindLink(ctx context.Context, id string) (string, error)

	// StorageStatus проверяет доступность хранилища.
	// Принимает контекст.
	// Возвращает true, если хранилище доступно, и ошибку в противном случае.
	StorageStatus(ctx context.Context) (bool, error)

	// BatchShorten создает сокращенные URL для пакета длинных URL.
	// Принимает контекст, массив запросов на сокращение и идентификатор пользователя.
	// Возвращает массив сокращенных URL и ошибку, если операция не удалась.
	BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error)

	// GetUserURLs возвращает список всех URL, созданных пользователем.
	// Принимает контекст и идентификатор пользователя.
	// Возвращает массив URL и ошибку, если операция не удалась.
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error)

	// DeleteURLs помечает указанные URL как удаленные.
	// Принимает контекст, массив идентификаторов URL и идентификатор пользователя.
	// Возвращает количество удаленных URL и ошибку, если операция не удалась.
	DeleteURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error)
}

// InitService создает и возвращает новый экземпляр Service с предоставленным репозиторием.
// Принимает реализацию интерфейса LinkRepository.
// Возвращает инициализированный сервис.
func InitService(repo repository.LinkRepository) *Service {
	return &Service{repo: repo}
}

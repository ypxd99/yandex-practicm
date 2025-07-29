package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
)

// LinkRepository определяет интерфейс для работы с хранилищем сокращенных URL.
// Предоставляет методы для создания, поиска и управления URL в базе данных.
type LinkRepository interface {
	// CreateLink создает новую запись сокращенного URL в хранилище.
	// Возвращает созданную запись и ошибку, если операция не удалась.
	CreateLink(ctx context.Context, id, url string, userID uuid.UUID) (*model.Link, error)

	// FindLink находит запись сокращенного URL по его идентификатору.
	// Возвращает найденную запись и ошибку, если URL не найден.
	FindLink(ctx context.Context, id string) (*model.Link, error)

	// FindUserLinks возвращает все URL, созданные указанным пользователем.
	// Возвращает массив URL и ошибку, если операция не удалась.
	FindUserLinks(ctx context.Context, userID uuid.UUID) ([]model.Link, error)

	// BatchCreate создает несколько записей сокращенных URL в хранилище.
	// Возвращает ошибку, если операция не удалась.
	BatchCreate(ctx context.Context, links []model.Link) error

	// MarkDeletedURLs помечает указанные URL как удаленные.
	// Возвращает количество удаленных URL и ошибку, если операция не удалась.
	MarkDeletedURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error)

	// Status проверяет доступность хранилища.
	// Возвращает true, если хранилище доступно, и ошибку в противном случае.
	Status(ctx context.Context) (bool, error)

	// Close закрывает соединение с хранилищем.
	// Возвращает ошибку, если закрытие не удалось.
	Close() error

	// GetStats возвращает статистику сервиса.
	// Возвращает количество URL и пользователей, а также ошибку, если операция не удалась.
	GetStats(ctx context.Context) (int64, int64, error)
}

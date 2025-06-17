package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
)

// ErrIDExists ошибка, возникающая при попытке создать ссылку с уже существующим ID
// ErrNotFound ошибка, возникающая при попытке найти несуществующую ссылку
// ErrStorageAccess ошибка, возникающая при проблемах с доступом к хранилищу
var (
	ErrIDExists = errors.New("ID already exists")
	ErrNotFound = errors.New("link not found")
	ErrStorageAccess = errors.New("storage access error")
)

// LocalStorage представляет локальное хранилище для сокращенных URL.
// Использует файловую систему для персистентного хранения данных.
type LocalStorage struct {
	mu       sync.RWMutex
	links    map[string]linkData
	filePath string
}

// linkData представляет структуру данных для хранения информации об URL.
type linkData struct {
	URL       string    // Оригинальный URL
	UserID    uuid.UUID // Идентификатор пользователя
	IsDeleted bool      // Флаг удаления
}

// fileLinks представляет структуру для сериализации данных в JSON.
type fileLinks struct {
	UUID        string    `json:"uuid"`         // Уникальный идентификатор записи
	ShortURL    string    `json:"short_url"`    // Сокращенный URL
	OriginalURL string    `json:"original_url"` // Оригинальный URL
	UserID      uuid.UUID `json:"user_id"`      // Идентификатор пользователя
	IsDeleted   bool      `json:"is_deleted"`   // Флаг удаления
}

// InitStorage создает и инициализирует новое локальное хранилище.
// Если указан путь к файлу, загружает данные из него.
// Возвращает инициализированное хранилище и ошибку, если инициализация не удалась.
func InitStorage(filePath string) (*LocalStorage, error) {
	s := &LocalStorage{
		links:    make(map[string]linkData),
		filePath: filePath,
	}

	if filePath != "" {
		if err := s.readFromFile(); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// CreateLink создает новую запись сокращенного URL в хранилище.
// Возвращает созданную запись и ошибку, если операция не удалась.
func (s *LocalStorage) CreateLink(ctx context.Context, id, url string, userID uuid.UUID) (*model.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.links[id]; exists {
		return nil, ErrIDExists
	}

	s.links[id] = linkData{
		URL:       url,
		UserID:    userID,
		IsDeleted: false,
	}

	if s.filePath != "" {
		if err := s.writeToFile(); err != nil {
			delete(s.links, id)
			return nil, err
		}
	}

	return &model.Link{ID: id, Link: url, UserID: userID, IsDeleted: false}, nil
}

// FindLink находит запись сокращенного URL по его идентификатору.
// Возвращает найденную запись и ошибку, если URL не найден.
func (s *LocalStorage) FindLink(ctx context.Context, id string) (*model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.links[id]
	if !exists {
		return nil, ErrNotFound
	}

	return &model.Link{ID: id, Link: data.URL, UserID: data.UserID, IsDeleted: data.IsDeleted}, nil
}

// FindUserLinks возвращает все URL, созданные указанным пользователем.
// Возвращает массив URL и ошибку, если операция не удалась.
func (s *LocalStorage) FindUserLinks(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []model.Link
	for id, data := range s.links {
		if data.UserID == userID && !data.IsDeleted {
			result = append(result, model.Link{
				ID:        id,
				Link:      data.URL,
				UserID:    userID,
				IsDeleted: data.IsDeleted,
			})
		}
	}

	return result, nil
}

// BatchCreate создает несколько записей сокращенных URL в хранилище.
// Возвращает ошибку, если операция не удалась.
func (s *LocalStorage) BatchCreate(ctx context.Context, links []model.Link) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, link := range links {
		s.links[link.ID] = linkData{
			URL:       link.Link,
			UserID:    link.UserID,
			IsDeleted: link.IsDeleted,
		}
	}

	if s.filePath != "" {
		if err := s.writeToFile(); err != nil {
			for _, link := range links {
				delete(s.links, link.ID)
			}
			return err
		}
	}
	return nil
}

// MarkDeletedURLs помечает указанные URL как удаленные.
// Возвращает количество удаленных URL и ошибку, если операция не удалась.
func (s *LocalStorage) MarkDeletedURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0
	for _, id := range ids {
		data, exists := s.links[id]
		if exists && data.UserID == userID && !data.IsDeleted {
			data.IsDeleted = true
			s.links[id] = data
			count++
		}
	}

	if count > 0 && s.filePath != "" {
		if err := s.writeToFile(); err != nil {
			return 0, err
		}
	}

	return count, nil
}

// Close закрывает хранилище и сохраняет данные в файл.
// Возвращает ошибку, если операция не удалась.
func (s *LocalStorage) Close() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.filePath != "" {
		return s.writeToFile()
	}
	return nil
}

// Status проверяет доступность хранилища.
// Всегда возвращает true, так как хранилище всегда доступно.
func (s *LocalStorage) Status(ctx context.Context) (bool, error) {
	return true, nil
}

// readFromFile загружает данные из файла в хранилище.
// Возвращает ошибку, если операция не удалась.
func (s *LocalStorage) readFromFile() error {
	if _, err := os.Stat(s.filePath); os.IsNotExist(err) {
		file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return ErrStorageAccess
		}
		file.Close()
		return nil
	}

	file, err := os.OpenFile(s.filePath, os.O_RDONLY, 0644)
	if err != nil {
		return ErrStorageAccess
	}
	defer file.Close()

	var links []fileLinks
	if err := json.NewDecoder(file).Decode(&links); err != nil {
		if err.Error() == "EOF" {
			return nil
		}
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	for _, link := range links {
		s.links[link.ShortURL] = linkData{
			URL:       link.OriginalURL,
			UserID:    link.UserID,
			IsDeleted: link.IsDeleted,
		}
	}

	return nil
}

// writeToFile сохраняет данные из хранилища в файл.
// Возвращает ошибку, если операция не удалась.
func (s *LocalStorage) writeToFile() error {
	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return ErrStorageAccess
	}
	defer file.Close()

	var links []fileLinks
	for shortURL, data := range s.links {
		links = append(links, fileLinks{
			UUID:        uuid.New().String(),
			ShortURL:    shortURL,
			OriginalURL: data.URL,
			UserID:      data.UserID,
			IsDeleted:   data.IsDeleted,
		})
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(links)
}

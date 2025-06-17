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

var (
	ErrIDExists      = errors.New("ID already exists")
	ErrNotFound      = errors.New("link not found")
	ErrStorageAccess = errors.New("storage access error")
)

type LocalStorage struct {
	mu       sync.RWMutex
	links    map[string]linkData
	filePath string
}

type linkData struct {
	URL       string
	UserID    uuid.UUID
	IsDeleted bool
}

type fileLinks struct {
	UUID        string    `json:"uuid"`
	ShortURL    string    `json:"short_url"`
	OriginalURL string    `json:"original_url"`
	UserID      uuid.UUID `json:"user_id"`
	IsDeleted   bool      `json:"is_deleted"`
}

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

func (s *LocalStorage) FindLink(ctx context.Context, id string) (*model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, exists := s.links[id]
	if !exists {
		return nil, ErrNotFound
	}

	return &model.Link{ID: id, Link: data.URL, UserID: data.UserID, IsDeleted: data.IsDeleted}, nil
}

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

func (s *LocalStorage) Close() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.filePath != "" {
		return s.writeToFile()
	}
	return nil
}

func (s *LocalStorage) Status(ctx context.Context) (bool, error) {
	return true, nil
}

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

package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/ypxd99/yandex-practicm/internal/model"
)

var (
	ErrIDExists = errors.New("ID already exists")
	ErrNotFound = errors.New("link not found")
)

type LocalStorage struct {
	mu    sync.RWMutex
	links map[string]string
}

func InitStorage() *LocalStorage {
	return &LocalStorage{
		links: make(map[string]string),
	}
}

func (s *LocalStorage) CreateLink(ctx context.Context, id, url string) (*model.Link, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.links[id]; exists {
		return nil, ErrIDExists
	}

	s.links[id] = url
	return &model.Link{ID: id, Link: url}, nil
}

func (s *LocalStorage) FindLink(ctx context.Context, id string) (*model.Link, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.links[id]
	if !exists {
		return nil, ErrNotFound
	}

	return &model.Link{ID: id, Link: url}, nil
}

func (s *LocalStorage) Close() error {
	return nil
}

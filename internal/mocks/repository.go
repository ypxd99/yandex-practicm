package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) CreateLink(ctx context.Context, id, url string, userID uuid.UUID) (*model.Link, error) {
	args := m.Called(ctx, id, url, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) FindLink(ctx context.Context, id string) (*model.Link, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) FindUserLinks(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Link), args.Error(1)
}

func (m *MockLinkRepository) BatchCreate(ctx context.Context, links []model.Link) error {
	args := m.Called(ctx, links)
	return args.Error(0)
}

func (m *MockLinkRepository) Close() error {
	return nil
}

func (m *MockLinkRepository) Status(ctx context.Context) (bool, error) {
	return true, nil
}

var _ repository.LinkRepository = (*MockLinkRepository)(nil)

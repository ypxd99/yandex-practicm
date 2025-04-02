package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

type MockLinkRepository struct {
	mock.Mock
}

func (m *MockLinkRepository) CreateLink(ctx context.Context, id, url string) (*model.Link, error) {
	args := m.Called(ctx, id, url)
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) FindLink(ctx context.Context, id string) (*model.Link, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model.Link), args.Error(1)
}

func (m *MockLinkRepository) Close() error {
	return nil
}

func (m *MockLinkRepository) Status(ctx context.Context) (bool, error) {
	return true, nil
}

var _ repository.LinkRepository = (*MockLinkRepository)(nil)

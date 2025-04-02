package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/service"
)

type MockLinkService struct {
	mock.Mock
}

func (m *MockLinkService) ShorterLink(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)
	return args.String(0), args.Error(1)
}

func (m *MockLinkService) FindLink(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

func (m *MockLinkService) StorageStatus(ctx context.Context) (bool, error) {
	args := m.Called(ctx)
	return args.Bool(0), args.Error(1)
}

var _ service.LinkService = (*MockLinkService)(nil)

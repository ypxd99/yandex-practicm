package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
)

type MockLinkService struct {
	mock.Mock
}

func (m *MockLinkService) ShorterLink(ctx context.Context, url string, userID uuid.UUID) (string, error) {
	args := m.Called(ctx, url, userID)
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

func (m *MockLinkService) BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error) {
	args := m.Called(ctx, batch, userID)
	return args.Get(0).([]model.BatchResponse), args.Error(1)
}

func (m *MockLinkService) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.UserURLResponse), args.Error(1)
}

var _ service.LinkService = (*MockLinkService)(nil)

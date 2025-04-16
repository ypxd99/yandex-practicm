package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

type Service struct {
	repo repository.LinkRepository
}

type LinkService interface {
	ShorterLink(ctx context.Context, url string, userID uuid.UUID) (string, error)
	FindLink(ctx context.Context, id string) (string, error)
	StorageStatus(ctx context.Context) (bool, error)
	BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error)
	GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error)
	DeleteURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error)
}

func InitService(repo repository.LinkRepository) *Service {
	return &Service{repo: repo}
}

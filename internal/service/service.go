package service

import (
	"context"

	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/repository"
)

type Service struct {
	repo repository.LinkRepository
}

type LinkService interface {
	ShorterLink(ctx context.Context, url string) (string, error)
	FindLink(ctx context.Context, id string) (string, error)
	StorageStatus(ctx context.Context) (bool, error)
	BatchShorten(ctx context.Context, batch []model.BatchRequest) ([]model.BatchResponse, error)
}

func InitService(repo repository.LinkRepository) *Service {
	return &Service{repo: repo}
}

package service

import (
	"context"

	"github.com/ypxd99/yandex-practicm/internal/repository"
)

type Service struct {
	repo repository.LinkRepository
}

type LinkService interface {
	ShorterLink(ctx context.Context, url string) (string, error)
	FindLink(ctx context.Context, id string) (string, error)
}

func InitService(repo repository.LinkRepository) *Service {
	return &Service{repo: repo}
}

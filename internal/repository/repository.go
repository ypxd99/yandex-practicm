package repository

import (
	"context"

	"github.com/ypxd99/yandex-practicm/internal/model"
)

type LinkRepository interface {
	CreateLink(ctx context.Context, id, url string) (*model.Link, error)
	FindLink(ctx context.Context, id string) (*model.Link, error)
	BatchCreate(ctx context.Context, links []model.Link) error
	Status(ctx context.Context) (bool, error)
	Close() error
}

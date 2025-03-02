package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ypxd99/yandex-practicm/internal/repository"
	"github.com/ypxd99/yandex-practicm/internal/repository/postgres"
	"github.com/ypxd99/yandex-practicm/internal/repository/storage"
	"github.com/ypxd99/yandex-practicm/util"
)

func TestCreateAndFindLink(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testID := "test123"
	testURL := "https://example.com"
	var repo repository.LinkRepository
	if cfg.Postgres.UsePostgres {
		postgresRepo, err := postgres.Connect(context.Background())
		assert.NoError(t, err)
		repo = postgresRepo
	} else {
		repo = storage.InitStorage()
	}
	defer repo.Close()

	createdLink, err := repo.CreateLink(ctx, testID, testURL)
	assert.NoError(t, err)
	assert.Equal(t, testID, createdLink.ID)
	assert.Equal(t, testURL, createdLink.Link)

	foundLink, err := repo.FindLink(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, testID, foundLink.ID)
	assert.Equal(t, testURL, foundLink.Link)

	_, err = repo.FindLink(ctx, "non-existent")
	assert.Error(t, sql.ErrNoRows, err)
}

package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
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
	testUserID := uuid.New()

	var (
		repo repository.LinkRepository
		err  error
	)
	if cfg.Postgres.UsePostgres {
		postgresRepo, err := postgres.Connect(context.Background())
		assert.NoError(t, err)
		repo = postgresRepo
	} else {
		repo, err = storage.InitStorage(cfg.FileStoragePath)
		if err != nil {
			util.GetLogger().Fatal(err)
			return
		}
	}
	defer repo.Close()

	createdLink, err := repo.CreateLink(ctx, testID, testURL, testUserID)
	assert.NoError(t, err)
	assert.Equal(t, testID, createdLink.ID)
	assert.Equal(t, testURL, createdLink.Link)
	assert.Equal(t, testUserID, createdLink.UserID)

	foundLink, err := repo.FindLink(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, testID, foundLink.ID)
	assert.Equal(t, testURL, foundLink.Link)
	assert.Equal(t, testUserID, foundLink.UserID)

	_, err = repo.FindLink(ctx, "non-existent")
	assert.Error(t, sql.ErrNoRows, err)

	userLinks, err := repo.FindUserLinks(ctx, testUserID)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(userLinks), 1)

	for _, link := range userLinks {
		assert.Equal(t, testUserID, link.UserID)
	}
}

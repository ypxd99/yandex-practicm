package postgres_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ypxd99/yandex-practicm/internal/repository/postgres"
)

func TestCreateAndFindLink(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	testID := "test123"
	testURL := "https://example.com"
	db, err := postgres.Connect(ctx)
	assert.NoError(t, err)

	createdLink, err := db.CreateLink(ctx, testID, testURL)
	assert.NoError(t, err)
	assert.Equal(t, testID, createdLink.ID)
	assert.Equal(t, testURL, createdLink.Link)

	foundLink, err := db.FindLink(ctx, testID)
	assert.NoError(t, err)
	assert.Equal(t, testID, foundLink.ID)
	assert.Equal(t, testURL, foundLink.Link)

	_, err = db.FindLink(ctx, "non-existent")
	assert.Error(t, sql.ErrNoRows, err)
}

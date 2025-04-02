package service_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/ypxd99/yandex-practicm/internal/mocks"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/util"
)

func TestShorterLink(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		expected := &model.Link{ID: "abc123", Link: "https://example.com"}
		mockRepo.On("CreateLink", ctx, mock.Anything, "https://example.com").
			Return(expected, nil).
			Once()

		id, err := svc.ShorterLink(ctx, "https://example.com")
		if err != nil {
			if !errors.Is(err, service.ErrURLExist) {
				assert.NoError(t, err)
			}
		}

		assert.NotEmpty(t, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		mockRepo.On("CreateLink", ctx, mock.Anything, "https://error.com").
			Return((*model.Link)(nil), errors.New("db error")).
			Once()

		_, err := svc.ShorterLink(ctx, "https://error.com")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindLink(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()

	t.Run("found existing link", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		expected := &model.Link{ID: "abc123", Link: "https://example.com"}
		mockRepo.On("FindLink", ctx, "abc123").
			Return(expected, nil).
			Once()

		url, err := svc.FindLink(ctx, "abc123")

		assert.NoError(t, err)
		assert.Equal(t, "https://example.com", url)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		mockRepo.On("FindLink", ctx, "notfound").
			Return((*model.Link)(nil), errors.New("not found")).
			Once()

		_, err := svc.FindLink(ctx, "notfound")

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBatchShorten(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()
	t.Run("save batch", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		mockRepo.On("BatchCreate", ctx, mock.AnythingOfType("[]model.Link")).
			Return(nil).
			Once()

		batch := []model.BatchRequest{
			{CorrelationID: "1", OriginalURL: "https://example.com"},
			{CorrelationID: "2", OriginalURL: "https://yandex.ru"},
		}

		result, err := svc.BatchShorten(context.Background(), batch)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

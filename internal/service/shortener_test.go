package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
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
	testUserID := uuid.New()

	t.Run("successful creation", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		expected := &model.Link{ID: "abc123", Link: "https://example.com", UserID: testUserID, IsDeleted: false}
		mockRepo.On("CreateLink", ctx, mock.Anything, "https://example.com", testUserID).
			Return(expected, nil).
			Once()

		id, err := svc.ShorterLink(ctx, "https://example.com", testUserID)
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

		mockRepo.On("CreateLink", ctx, mock.Anything, "https://error.com", testUserID).
			Return((*model.Link)(nil), errors.New("db error")).
			Once()

		_, err := svc.ShorterLink(ctx, "https://error.com", testUserID)

		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestFindLink(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()
	testUserID := uuid.New()

	t.Run("found existing link", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		expected := &model.Link{ID: "abc123", Link: "https://example.com", UserID: testUserID, IsDeleted: false}
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
	testUserID := uuid.New()

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

		result, err := svc.BatchShorten(ctx, batch, testUserID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetUserURLs(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()
	testUserID := uuid.New()

	t.Run("get user urls successful", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		links := []model.Link{
			{ID: "abc123", Link: "https://example.com", UserID: testUserID, IsDeleted: false},
			{ID: "def456", Link: "https://yandex.ru", UserID: testUserID, IsDeleted: false},
		}

		mockRepo.On("FindUserLinks", ctx, testUserID).
			Return(links, nil).
			Once()

		result, err := svc.GetUserURLs(ctx, testUserID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, "https://example.com", result[0].OriginalURL)
		mockRepo.AssertExpectations(t)
	})

	t.Run("get user urls empty", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		var emptyLinks []model.Link

		mockRepo.On("FindUserLinks", ctx, testUserID).
			Return(emptyLinks, nil).
			Once()

		result, err := svc.GetUserURLs(ctx, testUserID)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteURLs(t *testing.T) {
	cfg := util.GetConfig()
	util.InitLogger(cfg.Logger)
	ctx := context.Background()
	testUserID := uuid.New()

	t.Run("delete urls successful", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		ids := []string{"abc123", "def456"}

		mockRepo.On("MarkDeletedURLs", mock.AnythingOfType("*context.timerCtx"), ids, testUserID).
			Return(2, nil).
			Once()

		count, err := svc.DeleteURLs(ctx, ids, testUserID)

		assert.NoError(t, err)
		assert.Equal(t, len(ids), count)
		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty ids list", func(t *testing.T) {
		mockRepo := new(mocks.MockLinkRepository)
		svc := service.InitService(mockRepo)

		var emptyIDs []string

		mockRepo.On("MarkDeletedURLs", mock.AnythingOfType("*context.timerCtx"), emptyIDs, testUserID).
			Return(0, nil).
			Once()

		count, err := svc.DeleteURLs(ctx, emptyIDs, testUserID)

		assert.NoError(t, err)
		assert.Equal(t, 0, count)
		time.Sleep(100 * time.Millisecond)
		mockRepo.AssertExpectations(t)
	})
}

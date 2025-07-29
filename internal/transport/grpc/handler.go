package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/internal/service"
	"github.com/ypxd99/yandex-practicm/proto/shortener"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler представляет gRPC обработчик сервиса сокращения URL.
// Содержит бизнес-логику сервиса и предоставляет gRPC-эндпоинты
// для операций сокращения URL.
type GRPCHandler struct {
	service service.LinkService
	shortener.UnimplementedShortenerServiceServer
}

// NewGRPCHandler создает и возвращает новый экземпляр GRPCHandler с предоставленным сервисом.
// Принимает реализацию интерфейса LinkService.
// Возвращает инициализированный обработчик.
func NewGRPCHandler(service service.LinkService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

// ShortenLink создает сокращенную ссылку для указанного URL (POST /).
func (h *GRPCHandler) ShortenLink(ctx context.Context, req *shortener.ShortenLinkRequest) (*shortener.ShortenLinkResponse, error) {
	if req.Url == "" {
		return &shortener.ShortenLinkResponse{
			ShortUrl: "",
			Error:    "URL is required",
		}, status.Error(codes.InvalidArgument, "URL is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &shortener.ShortenLinkResponse{
			ShortUrl: "",
			Error:    "Invalid user ID",
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	shortURL, err := h.service.ShorterLink(ctx, req.Url, userID)
	if err != nil {
		if err == service.ErrURLExist {
			return &shortener.ShortenLinkResponse{
				ShortUrl: shortURL,
				Error:    "",
			}, status.Error(codes.AlreadyExists, "URL already exists")
		}
		return &shortener.ShortenLinkResponse{
			ShortUrl: "",
			Error:    err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &shortener.ShortenLinkResponse{
		ShortUrl: shortURL,
		Error:    "",
	}, nil
}

// Shorten создает сокращенную ссылку через API (POST /api/shorten).
func (h *GRPCHandler) Shorten(ctx context.Context, req *shortener.ShortenRequest) (*shortener.ShortenResponse, error) {
	if req.Url == "" {
		return &shortener.ShortenResponse{
			ShortUrl: "",
			Error:    "URL is required",
		}, status.Error(codes.InvalidArgument, "URL is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &shortener.ShortenResponse{
			ShortUrl: "",
			Error:    "Invalid user ID",
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	shortURL, err := h.service.ShorterLink(ctx, req.Url, userID)
	if err != nil {
		if err == service.ErrURLExist {
			return &shortener.ShortenResponse{
				ShortUrl: shortURL,
				Error:    "",
			}, status.Error(codes.AlreadyExists, "URL already exists")
		}
		return &shortener.ShortenResponse{
			ShortUrl: "",
			Error:    err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &shortener.ShortenResponse{
		ShortUrl: shortURL,
		Error:    "",
	}, nil
}

// GetOriginalURL получает оригинальный URL по сокращенному идентификатору (GET /:id).
func (h *GRPCHandler) GetOriginalURL(ctx context.Context, req *shortener.GetOriginalURLRequest) (*shortener.GetOriginalURLResponse, error) {
	if req.ShortId == "" {
		return &shortener.GetOriginalURLResponse{
			OriginalUrl: "",
			Error:       "Short ID is required",
		}, status.Error(codes.InvalidArgument, "Short ID is required")
	}

	originalURL, err := h.service.FindLink(ctx, req.ShortId)
	if err != nil {
		if err == service.ErrURLDeleted {
			return &shortener.GetOriginalURLResponse{
				OriginalUrl: "",
				Error:       "URL was deleted",
			}, status.Error(codes.NotFound, "URL was deleted")
		}
		return &shortener.GetOriginalURLResponse{
			OriginalUrl: "",
			Error:       err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &shortener.GetOriginalURLResponse{
		OriginalUrl: originalURL,
		Error:       "",
	}, nil
}

// BatchShorten создает сокращенные ссылки для пакета URL (POST /api/shorten/batch).
func (h *GRPCHandler) BatchShorten(ctx context.Context, req *shortener.BatchShortenRequest) (*shortener.BatchShortenResponse, error) {
	if len(req.Items) == 0 {
		return &shortener.BatchShortenResponse{
			Items: nil,
			Error: "Batch items are required",
		}, status.Error(codes.InvalidArgument, "Batch items are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &shortener.BatchShortenResponse{
			Items: nil,
			Error: "Invalid user ID",
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	// Конвертируем gRPC запрос в модель
	batchItems := make([]model.BatchRequest, len(req.Items))
	for i, item := range req.Items {
		batchItems[i] = model.BatchRequest{
			CorrelationID: item.CorrelationId,
			OriginalURL:   item.OriginalUrl,
		}
	}

	// Вызываем сервис
	batchResponse, err := h.service.BatchShorten(ctx, batchItems, userID)
	if err != nil {
		return &shortener.BatchShortenResponse{
			Items: nil,
			Error: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// Конвертируем ответ в gRPC формат
	responseItems := make([]*shortener.BatchResponseItem, len(batchResponse))
	for i, item := range batchResponse {
		responseItems[i] = &shortener.BatchResponseItem{
			CorrelationId: item.CorrelationID,
			ShortUrl:      item.ShortURL,
		}
	}

	return &shortener.BatchShortenResponse{
		Items: responseItems,
		Error: "",
	}, nil
}

// GetUserURLs возвращает все URL, созданные пользователем (GET /api/user/urls).
func (h *GRPCHandler) GetUserURLs(ctx context.Context, req *shortener.GetUserURLsRequest) (*shortener.GetUserURLsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &shortener.GetUserURLsResponse{
			Items: nil,
			Error: "Invalid user ID",
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	userURLs, err := h.service.GetUserURLs(ctx, userID)
	if err != nil {
		return &shortener.GetUserURLsResponse{
			Items: nil,
			Error: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	// Конвертируем ответ в gRPC формат
	responseItems := make([]*shortener.UserURLItem, len(userURLs))
	for i, item := range userURLs {
		responseItems[i] = &shortener.UserURLItem{
			ShortUrl:    item.ShortURL,
			OriginalUrl: item.OriginalURL,
		}
	}

	return &shortener.GetUserURLsResponse{
		Items: responseItems,
		Error: "",
	}, nil
}

// DeleteURLs помечает указанные URL как удаленные (DELETE /api/user/urls).
func (h *GRPCHandler) DeleteURLs(ctx context.Context, req *shortener.DeleteURLsRequest) (*shortener.DeleteURLsResponse, error) {
	if len(req.ShortIds) == 0 {
		return &shortener.DeleteURLsResponse{
			DeletedCount: 0,
			Error:        "Short IDs are required",
		}, status.Error(codes.InvalidArgument, "Short IDs are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return &shortener.DeleteURLsResponse{
			DeletedCount: 0,
			Error:        "Invalid user ID",
		}, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	deletedCount, err := h.service.DeleteURLs(ctx, req.ShortIds, userID)
	if err != nil {
		return &shortener.DeleteURLsResponse{
			DeletedCount: 0,
			Error:        err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &shortener.DeleteURLsResponse{
		DeletedCount: int32(deletedCount),
		Error:        "",
	}, nil
}

// Ping проверяет статус хранилища (GET /ping).
func (h *GRPCHandler) Ping(ctx context.Context, req *shortener.PingRequest) (*shortener.PingResponse, error) {
	storageStatus, err := h.service.StorageStatus(ctx)
	if err != nil {
		return &shortener.PingResponse{
			Status: false,
			Error:  err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &shortener.PingResponse{
		Status: storageStatus,
		Error:  "",
	}, nil
}

// GetStats возвращает статистику сервиса (GET /api/internal/stats).
func (h *GRPCHandler) GetStats(ctx context.Context, req *shortener.GetStatsRequest) (*shortener.GetStatsResponse, error) {
	urls, users, err := h.service.GetStats(ctx)
	if err != nil {
		return &shortener.GetStatsResponse{
			Urls:  0,
			Users: 0,
			Error: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &shortener.GetStatsResponse{
		Urls:  urls,
		Users: users,
		Error: "",
	}, nil
}

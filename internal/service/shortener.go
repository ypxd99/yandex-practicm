package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/util"
)

// ErrURLExist ошибка, возникающая при попытке создать уже существующий URL
// ErrURLDeleted ошибка, возникающая при попытке получить доступ к удаленному URL
var (
	ErrURLExist = errors.New("url already exists")
	ErrURLDeleted = errors.New("url is deleted")
)

// normalizeQuery нормализует запрос, удаляя SQL-ключевые слова.
// Принимает строку запроса и возвращает нормализованную строку.
func normalizeQuery(data string) string {
	keyWords := util.GetConfig().Postgres.SQLKeyWords
	res := ""
	if len(keyWords) > 0 {
		for _, s := range strings.Split(data, " ") {
			clean := true
			for _, v := range keyWords {
				if strings.Contains(strings.ToUpper(s), strings.ToUpper(v)) {
					clean = false
					break
				}
			}
			if clean {
				if len(res) > 0 {
					res += " "
				}
				res += s
			}
		}
	}

	return res
}

// generateShortID генерирует короткий идентификатор для URL.
// Возвращает строку с идентификатором и ошибку, если генерация не удалась.
func (s *Service) generateShortID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while reading rand")
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

// ShorterLink создает сокращенную версию URL.
// Принимает контекст, оригинальный URL и идентификатор пользователя.
// Возвращает сокращенный URL и ошибку, если операция не удалась.
func (s *Service) ShorterLink(ctx context.Context, req string, userID uuid.UUID) (string, error) {
	id, err := s.generateShortID()
	if err != nil {
		return "", err
	}
	link, err := s.repo.CreateLink(ctx, id, normalizeQuery(req), userID)
	if err != nil {
		return "", err
	}

	baseURL := util.GetConfig().Server.BaseURL
	if link.ID != id {
		var res strings.Builder
		res.WriteString(baseURL)
		res.WriteString("/")
		res.WriteString(link.ID)
		return res.String(), ErrURLExist
	}
	var res strings.Builder
	res.WriteString(baseURL)
	res.WriteString("/")
	res.WriteString(link.ID)
	return res.String(), nil
}

// FindLink находит оригинальный URL по его сокращенной версии.
// Принимает контекст и сокращенный URL.
// Возвращает оригинальный URL и ошибку, если URL не найден или удален.
func (s *Service) FindLink(ctx context.Context, req string) (string, error) {
	str := normalizeQuery(req)
	link, err := s.repo.FindLink(ctx, str)
	if err != nil {
		return "", err
	}

	if link.IsDeleted {
		return "", ErrURLDeleted
	}

	return link.Link, nil
}

// StorageStatus проверяет доступность хранилища.
// Принимает контекст.
// Возвращает статус доступности и ошибку, если проверка не удалась.
func (s *Service) StorageStatus(ctx context.Context) (bool, error) {
	return s.repo.Status(ctx)
}

// BatchShorten создает сокращенные версии для нескольких URL.
// Принимает контекст, массив запросов на сокращение и идентификатор пользователя.
// Возвращает массив ответов с сокращенными URL и ошибку, если операция не удалась.
func (s *Service) BatchShorten(ctx context.Context, batch []model.BatchRequest, userID uuid.UUID) ([]model.BatchResponse, error) {
	resp := make([]model.BatchResponse, len(batch))
	links := make([]model.Link, len(batch))
	baseURL := util.GetConfig().Server.BaseURL

	for i, item := range batch {
		shortURL, err := s.generateShortID()
		if err != nil {
			return nil, err
		}
		links[i] = model.Link{
			ID:     shortURL,
			Link:   item.OriginalURL,
			UserID: userID,
		}
		var res strings.Builder
		res.WriteString(baseURL)
		res.WriteString("/")
		res.WriteString(shortURL)
		resp[i] = model.BatchResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      res.String(),
		}
	}

	err := s.repo.BatchCreate(ctx, links)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetUserURLs возвращает все URL, созданные указанным пользователем.
// Принимает контекст и идентификатор пользователя.
// Возвращает массив URL и ошибку, если операция не удалась.
func (s *Service) GetUserURLs(ctx context.Context, userID uuid.UUID) ([]model.UserURLResponse, error) {
	links, err := s.repo.FindUserLinks(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return []model.UserURLResponse{}, nil
	}

	result := make([]model.UserURLResponse, len(links))
	baseURL := util.GetConfig().Server.BaseURL
	for i, link := range links {
		var res strings.Builder
		res.WriteString(baseURL)
		res.WriteString("/")
		res.WriteString(link.ID)
		result[i] = model.UserURLResponse{
			ShortURL:    res.String(),
			OriginalURL: link.Link,
		}
	}

	return result, nil
}

// DeleteURLs помечает указанные URL как удаленные.
// Принимает контекст, массив идентификаторов URL и идентификатор пользователя.
// Возвращает количество удаленных URL и ошибку, если операция не удалась.
func (s *Service) DeleteURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error) {
	count, err := s.repo.MarkDeletedURLs(ctx, ids, userID)
	if err != nil {
		util.GetLogger().Errorf("failed to mark URLs as deleted: %v", err)
		return 0, err
	}

	util.GetLogger().Infof("marked %d URLs as deleted", count)
	return count, nil
}

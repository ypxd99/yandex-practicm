package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/internal/model"
	"github.com/ypxd99/yandex-practicm/util"
)

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

func (s *Service) generateShortID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while reading rand")
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func (s *Service) ShorterLink(ctx context.Context, req string) (string, error) {
	str := normalizeQuery(req)
	id, err := s.generateShortID()
	if err != nil {
		return "", err
	}
	link, err := s.repo.CreateLink(ctx, id, str)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", util.GetConfig().Server.BaseURL, link.ID), nil
}

func (s *Service) FindLink(ctx context.Context, req string) (string, error) {
	str := normalizeQuery(req)
	link, err := s.repo.FindLink(ctx, str)
	if err != nil {
		return "", err
	}

	return link.Link, nil
}

func (s *Service) StorageStatus(ctx context.Context) (bool, error) {
	return s.repo.Status(ctx)
}

func (s *Service) BatchShorten(ctx context.Context, batch []model.BatchRequest) ([]model.BatchResponse, error) {
	resp := make([]model.BatchResponse, 0, len(batch))
	links := make([]model.Link, 0, len(batch))

	for _, item := range batch {
		shortURL, err := s.generateShortID()
		if err != nil {
			return nil, err
		}

		links = append(links, model.Link{
			ID:   shortURL,
			Link: item.OriginalURL,
		})

		resp = append(resp, model.BatchResponse{
			CorrelationID: item.CorrelationID,
			ShortURL:      fmt.Sprintf("%s/%s", util.GetConfig().Server.BaseURL, shortURL),
		})
	}

	err := s.repo.BatchCreate(ctx, links)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

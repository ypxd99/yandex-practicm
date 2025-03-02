package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/util"
)

func (s *Service) generateShortID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while reading rand")
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func (s *Service) ShorterLink(ctx context.Context, req string) (string, error) {
	id, err := s.generateShortID()
	if err != nil {
		return "", err
	}
	link, err := s.repo.CreateLink(ctx, id, req)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s:%d/%s", util.GetConfig().Server.Address, util.GetConfig().Server.Port, link.ID), nil
}

func (s *Service) FindLink(ctx context.Context, req string) (string, error) {
	link, err := s.repo.FindLink(ctx, req)
	if err != nil {
		return "", err
	}

	return link.Link, nil
}

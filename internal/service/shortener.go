package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/pkg/errors"
	"github.com/ypxd99/yandex-practicm/internal/repository/postgres"
)

func generateShortID() (string, error) {
	b := make([]byte, 6)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.WithMessage(err, "error occurred while reading rand")
	}
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "="), nil
}

func ShorterLink(ctx context.Context, req string) (string, error) {
	id, err := generateShortID()
	if err != nil {
		return "", err
	}
	link, err := postgres.CreateLink(ctx, id, req)
	if err != nil {
		return "", err
	}

	return link.ID, nil
}

func FindLink(ctx context.Context, req string) (string, error) {
	link, err := postgres.FindLink(ctx, req)
	if err != nil {
		return "", err
	}

	return link.Link, nil
}
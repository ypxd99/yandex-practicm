package postgres

import (
	"context"

	"github.com/ypxd99/yandex-practicm/internal/model"
)

func CreateLink(ctx context.Context, id, link string) (*model.Link, error) {
	var newLink model.Link

	db, err := Connect()
	if err != nil {
		return nil, err
	}

	query := `
		INSERT INTO shortener.links (id, link)
		VALUES (?, ?)
		RETURNING id, link;
	`

	err = db.NewRaw(query, id, link).Scan(ctx, &newLink)
	if err != nil {
		return nil, err
	}

	return &newLink, nil
}

func FindLink(ctx context.Context, id string) (*model.Link, error) {
	var (
		link  model.Link
		query = `
				SELECT id, link
				FROM shortener.links
				WHERE id = ?
				LIMIT 1;
			`
	)

	db, err := Connect()
	if err != nil {
		return nil, err
	}

	err = db.NewRaw(query, id).Scan(ctx, &link)
	if err != nil {
		return nil, err
	}

	return &link, err
}

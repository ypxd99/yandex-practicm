package postgres

import (
	"context"

	"github.com/ypxd99/yandex-practicm/internal/model"
)

func (p *Postgres) CreateLink(ctx context.Context, id, link string) (*model.Link, error) {
	var newLink model.Link

	query := `
		INSERT INTO shortener.links (id, link)
		VALUES (?, ?)
		RETURNING id, link;
	`

	err := p.db.NewRaw(query, id, link).Scan(ctx, &newLink)
	if err != nil {
		return nil, err
	}

	return &newLink, nil
}

func (p *Postgres)  FindLink(ctx context.Context, id string) (*model.Link, error) {
	var (
		link  model.Link
		query = `
				SELECT id, link
				FROM shortener.links
				WHERE id = ?
				LIMIT 1;
			`
	)

	err := p.db.NewRaw(query, id).Scan(ctx, &link)
	if err != nil {
		return nil, err
	}

	return &link, err
}

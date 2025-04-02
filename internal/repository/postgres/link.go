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
        ON CONFLICT (link) DO UPDATE SET link = EXCLUDED.link
        RETURNING id, link;
	`

	err := p.db.NewRaw(query, id, link).Scan(ctx, &newLink)
	if err != nil {
		return nil, err
	}

	return &newLink, nil
}

func (p *Postgres) FindLink(ctx context.Context, id string) (*model.Link, error) {
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

func (p *Postgres) BatchCreate(ctx context.Context, links []model.Link) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = tx.NewInsert().
		Model(&links).
		Exec(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

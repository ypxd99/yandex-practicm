package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
)

func (p *Postgres) CreateLink(ctx context.Context, id, link string, userID uuid.UUID) (*model.Link, error) {
	var newLink model.Link

	query := `
		INSERT INTO shortener.links (id, link, user_id)
        VALUES (?, ?, ?)
        ON CONFLICT (link) DO UPDATE SET link = EXCLUDED.link
        RETURNING id, link, user_id;
	`

	err := p.db.NewRaw(query, id, link, userID).Scan(ctx, &newLink)
	if err != nil {
		return nil, err
	}

	return &newLink, nil
}

func (p *Postgres) FindLink(ctx context.Context, id string) (*model.Link, error) {
	var (
		link  model.Link
		query = `
				SELECT id, link, user_id
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

func (p *Postgres) FindUserLinks(ctx context.Context, userID uuid.UUID) ([]model.Link, error) {
	var (
		links []model.Link
		query = `
				SELECT id, link, user_id
				FROM shortener.links
				WHERE user_id = ?
				ORDER BY id;
			`
	)

	err := p.db.NewRaw(query, userID).Scan(ctx, &links)
	if err != nil {
		return nil, err
	}

	return links, nil
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

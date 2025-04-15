package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/ypxd99/yandex-practicm/internal/model"
)

func (p *Postgres) CreateLink(ctx context.Context, id, link string, userID uuid.UUID) (*model.Link, error) {
	var newLink model.Link

	query := `
		INSERT INTO shortener.links (id, link, user_id, is_deleted)
        VALUES (?, ?, ?, false)
        ON CONFLICT (link) DO UPDATE SET link = EXCLUDED.link
        RETURNING id, link, user_id, is_deleted;
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
				SELECT id, link, user_id, is_deleted
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
				SELECT id, link, user_id, is_deleted
				FROM shortener.links
				WHERE user_id = ? AND is_deleted = false
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

func (p *Postgres) MarkDeletedURLs(ctx context.Context, ids []string, userID uuid.UUID) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	result, err := p.db.NewUpdate().
		Table("shortener.links").
		Set("is_deleted = true").
		Where("id IN (?) AND user_id = ? AND is_deleted = false", ids, userID).
		Exec(ctx)

	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

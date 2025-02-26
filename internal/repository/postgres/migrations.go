package postgres

import (
	"context"

	"github.com/pressly/goose/v3"
)

func MigrateDBUp(ctx context.Context) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	return goose.UpContext(ctx, db.DB, "./migration")
}

func MigrateDBDown(ctx context.Context) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	return goose.DownContext(ctx, db.DB, "./migration")
}

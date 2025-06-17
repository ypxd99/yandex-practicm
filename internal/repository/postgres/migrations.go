package postgres

import (
	"context"

	"github.com/pressly/goose/v3"
)

// MigrateDBUp применяет миграции базы данных вверх.
// Принимает контекст для управления временем жизни операции.
// Возвращает ошибку в случае неудачи.
func MigrateDBUp(ctx context.Context) error {
	db, err := Connect(ctx)
	if err != nil {
		return err
	}
	return goose.UpContext(ctx, db.db.DB, "./migration")
}

// MigrateDBDown откатывает миграции базы данных вниз.
// Принимает контекст для управления временем жизни операции.
// Возвращает ошибку в случае неудачи.
func MigrateDBDown(ctx context.Context) error {
	db, err := Connect(ctx)
	if err != nil {
		return err
	}
	return goose.DownContext(ctx, db.db.DB, "./migration")
}

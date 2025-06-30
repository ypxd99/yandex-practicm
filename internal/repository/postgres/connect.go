package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/ypxd99/yandex-practicm/util"
)

// Postgres представляет структуру для работы с PostgreSQL.
// Содержит соединение с базой данных.
type Postgres struct {
	db *bun.DB
}

// Connect устанавливает соединение с PostgreSQL.
// Принимает контекст для управления временем жизни соединения.
// Возвращает экземпляр Postgres и ошибку.
func Connect(ctx context.Context) (*Postgres, error) {
	cfg := util.GetConfig().Postgres
	// connStr := fmt.Sprintf(dbConnStr, cfg.User, cfg.Password, cfg.Address, cfg.DBName)

	// sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.ConnString)))
	db := bun.NewDB(sqlDB, pgdialect.New())

	err := db.Ping()
	if err != nil {
		return nil, err
	}
	db.DB.SetMaxOpenConns(cfg.MaxConn)
	db.DB.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifeTime) * time.Second)
	db.Exec(`SET search_path TO shortener, public;`)

	return &Postgres{db: db}, nil
}

// Close закрывает соединение с базой данных.
// Возвращает ошибку в случае неудачи.
func (p *Postgres) Close() error {
	return p.db.Close()
}

// Status проверяет доступность базы данных.
// Принимает контекст для управления временем жизни запроса.
// Возвращает статус доступности и ошибку.
func (p *Postgres) Status(ctx context.Context) (bool, error) {
	err := p.db.PingContext(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

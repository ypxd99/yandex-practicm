package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/ypxd99/yandex-practicm/util"
)

const dbConnStr = "postgres://%s:%s@%s/%s?sslmode=disable"

type Postgres struct {
	db *bun.DB
}

func Connect(ctx context.Context) (*Postgres, error) {
	cfg := util.GetConfig().Postgres
	connStr := fmt.Sprintf(dbConnStr, cfg.User, cfg.Password, cfg.Address, cfg.DBName)

	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	db := bun.NewDB(sqlDB, pgdialect.New())

	if cfg.Trace {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	db.DB.SetMaxOpenConns(cfg.MaxConn)
	db.DB.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifeTime) * time.Second)
	db.Exec(`SET search_path TO shortener, public;`)

	return &Postgres{db: db}, db.PingContext(ctx)
}

func (p *Postgres) Close() error {
	return p.db.Close()
}

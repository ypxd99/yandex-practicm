package postgres

import (
	"context"
	"database/sql"
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
	//connStr := fmt.Sprintf(dbConnStr, cfg.User, cfg.Password, cfg.Address, cfg.DBName)

	//sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	sqlDB := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.ConnString)))
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

func (p *Postgres) Status(ctx context.Context) (bool, error) {
	err := p.db.PingContext(ctx)
	if err != nil {
		return false, err
	}

	return true, nil
}

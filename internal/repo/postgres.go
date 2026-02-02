package repo

import (
	"context"
	"database/sql"
	"time"
)

const (
	dbDriverName    = "pgx"
	dbPingTimeout   = 3 * time.Second
	defaultMaxOpen  = 10
	defaultMaxIdle  = 10
	defaultConnLife = 30 * time.Minute
)

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open(dbDriverName, dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(defaultMaxOpen)
	db.SetMaxIdleConns(defaultMaxIdle)
	db.SetConnMaxLifetime(defaultConnLife)

	ctx, cancel := context.WithTimeout(context.Background(), dbPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

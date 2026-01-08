package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	*pgxpool.Pool
}

func New(ctx context.Context, dsn string, minConns, maxConns int32, connMaxIdle, connMaxLife time.Duration) (*Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	cfg.MinConns = minConns
	cfg.MaxConns = maxConns
	cfg.MaxConnIdleTime = connMaxIdle
	cfg.MaxConnLifetime = connMaxLife

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	return &Pool{pool}, nil
}

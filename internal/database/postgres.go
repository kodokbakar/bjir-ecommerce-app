package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func NewPostgresPool(ctx context.Context, cfg config.DatabaseConfig) (*pgxpool.Pool, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create postgres pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}

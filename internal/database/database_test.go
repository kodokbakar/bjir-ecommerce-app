package database_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
	"github.com/kodokbakar/go-ecommerce-api/internal/database"
)

func testDatabaseConfig() config.DatabaseConfig {
	return config.DatabaseConfig{
		Host:              "localhost",
		Port:              5432,
		User:              "postgres",
		Password:          "postgres",
		Name:              "go_ecommerce_api",
		SSLMode:           "disable",
		MaxConns:          10,
		MinConns:          1,
		MaxConnLifetime:   time.Hour,
		MaxConnIdleTime:   30 * time.Minute,
		HealthCheckPeriod: time.Minute,
	}
}

func TestNewPostgresPoolConfig(t *testing.T) {
	cfg := testDatabaseConfig()

	poolConfig, err := database.NewPostgresPoolConfig(cfg)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if poolConfig.MaxConns != cfg.MaxConns {
		t.Fatalf("expected max conns %d, got %d", cfg.MaxConns, poolConfig.MaxConns)
	}

	if poolConfig.MinConns != cfg.MinConns {
		t.Fatalf("expected min conns %d, got %d", cfg.MinConns, poolConfig.MinConns)
	}

	if poolConfig.MaxConnLifetime != cfg.MaxConnLifetime {
		t.Fatalf("expected max conn lifetime %s, got %s", cfg.MaxConnLifetime, poolConfig.MaxConnLifetime)
	}
}

func TestNewRedisOptions(t *testing.T) {
	cfg := config.RedisConfig{
		Host:         "redis",
		Port:         6379,
		Password:     "secret",
		DB:           2,
		PoolSize:     20,
		MinIdleConns: 4,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 4 * time.Second,
	}

	options := database.NewRedisOptions(cfg)

	if options.Addr != "redis:6379" {
		t.Fatalf("expected addr redis:6379, got %s", options.Addr)
	}

	if options.Password != "secret" {
		t.Fatalf("expected password secret, got %s", options.Password)
	}

	if options.DB != 2 {
		t.Fatalf("expected DB 2, got %d", options.DB)
	}

	if options.PoolSize != 20 {
		t.Fatalf("expected pool size 20, got %d", options.PoolSize)
	}
}

func TestRunMigrationsFromPathRequiresPath(t *testing.T) {
	err := database.RunMigrationsFromPath(testDatabaseConfig(), "   ")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "migrations path is required") {
		t.Fatalf("expected migrations path error, got %v", err)
	}
}

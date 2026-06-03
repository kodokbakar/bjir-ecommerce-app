package config_test

import (
	"strings"
	"testing"
	"time"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func clearConfigEnv(t *testing.T) {
	t.Helper()

	keys := []string{
		"APP_ENV",
		"APP_PORT",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"DB_POOL_MAX_CONNS",
		"DB_POOL_MIN_CONNS",
		"DB_POOL_MAX_CONN_LIFETIME",
		"DB_POOL_MAX_CONN_IDLE_TIME",
		"DB_POOL_HEALTH_CHECK_PERIOD",
		"REDIS_HOST",
		"REDIS_PORT",
		"REDIS_PASSWORD",
		"REDIS_DB",
		"REDIS_POOL_SIZE",
		"REDIS_MIN_IDLE_CONNS",
		"REDIS_DIAL_TIMEOUT",
		"REDIS_READ_TIMEOUT",
		"REDIS_WRITE_TIMEOUT",
		"JWT_SECRET",
		"JWT_EXPIRES_IN",
		"JWT_ISSUER",
	}

	for _, key := range keys {
		t.Setenv(key, "")
	}
}

func TestLoadConfigWithDefaults(t *testing.T) {
	clearConfigEnv(t)

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.App.Env != "development" {
		t.Fatalf("expected APP_ENV development, got %s", cfg.App.Env)
	}

	if cfg.App.Port != "8080" {
		t.Fatalf("expected APP_PORT 8080, got %s", cfg.App.Port)
	}

	if cfg.Database.Host != "localhost" {
		t.Fatalf("expected DB_HOST localhost, got %s", cfg.Database.Host)
	}

	if cfg.Database.Port != 5432 {
		t.Fatalf("expected DB_PORT 5432, got %d", cfg.Database.Port)
	}

	if cfg.Redis.Addr() != "localhost:6379" {
		t.Fatalf("expected redis addr localhost:6379, got %s", cfg.Redis.Addr())
	}

	if cfg.JWT.ExpiresIn != 24*time.Hour {
		t.Fatalf("expected JWT_EXPIRES_IN 24h, got %s", cfg.JWT.ExpiresIn)
	}
}

func TestLoadConfigWithEnvOverrides(t *testing.T) {
	clearConfigEnv(t)

	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_PORT", "9090")

	t.Setenv("DB_HOST", "postgres")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "app_user")
	t.Setenv("DB_PASSWORD", "secret")
	t.Setenv("DB_NAME", "app_db")
	t.Setenv("DB_SSLMODE", "disable")
	t.Setenv("DB_POOL_MAX_CONNS", "20")
	t.Setenv("DB_POOL_MIN_CONNS", "2")
	t.Setenv("DB_POOL_MAX_CONN_LIFETIME", "2h")
	t.Setenv("DB_POOL_MAX_CONN_IDLE_TIME", "10m")
	t.Setenv("DB_POOL_HEALTH_CHECK_PERIOD", "30s")

	t.Setenv("REDIS_HOST", "redis")
	t.Setenv("REDIS_PORT", "6380")
	t.Setenv("REDIS_PASSWORD", "redis_secret")
	t.Setenv("REDIS_DB", "1")
	t.Setenv("REDIS_POOL_SIZE", "20")
	t.Setenv("REDIS_MIN_IDLE_CONNS", "4")
	t.Setenv("REDIS_DIAL_TIMEOUT", "2s")
	t.Setenv("REDIS_READ_TIMEOUT", "4s")
	t.Setenv("REDIS_WRITE_TIMEOUT", "5s")

	t.Setenv("JWT_SECRET", "test_secret")
	t.Setenv("JWT_EXPIRES_IN", "2h")
	t.Setenv("JWT_ISSUER", "test-issuer")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.App.Env != "test" {
		t.Fatalf("expected APP_ENV test, got %s", cfg.App.Env)
	}

	if cfg.App.Port != "9090" {
		t.Fatalf("expected APP_PORT 9090, got %s", cfg.App.Port)
	}

	dsn := cfg.Database.DSN()
	if !strings.Contains(dsn, "postgres:5433") {
		t.Fatalf("expected DSN to contain postgres:5433, got %s", dsn)
	}

	if cfg.Redis.Addr() != "redis:6380" {
		t.Fatalf("expected redis addr redis:6380, got %s", cfg.Redis.Addr())
	}

	if cfg.JWT.Issuer != "test-issuer" {
		t.Fatalf("expected issuer test-issuer, got %s", cfg.JWT.Issuer)
	}

	if cfg.JWT.ExpiresIn != 2*time.Hour {
		t.Fatalf("expected expires in 2h, got %s", cfg.JWT.ExpiresIn)
	}
}

func TestLoadConfigInvalidDatabasePort(t *testing.T) {
	clearConfigEnv(t)

	t.Setenv("DB_PORT", "0")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DB_PORT") {
		t.Fatalf("expected DB_PORT error, got %v", err)
	}
}

func TestLoadConfigInvalidRedisDB(t *testing.T) {
	clearConfigEnv(t)

	t.Setenv("REDIS_DB", "-1")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "REDIS_DB") {
		t.Fatalf("expected REDIS_DB error, got %v", err)
	}
}

func TestLoadConfigInvalidDatabasePoolConfig(t *testing.T) {
	clearConfigEnv(t)

	t.Setenv("DB_POOL_MAX_CONNS", "2")
	t.Setenv("DB_POOL_MIN_CONNS", "3")

	_, err := config.LoadConfig()
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "DB_POOL_MIN_CONNS") {
		t.Fatalf("expected DB_POOL_MIN_CONNS error, got %v", err)
	}
}

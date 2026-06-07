package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type AppConfig struct {
	Env             string
	Port            string
	ShutdownTimeout time.Duration
}

type DatabaseConfig struct {
	URL               string
	Host              string
	Port              int
	User              string
	Password          string
	Name              string
	SSLMode           string
	MaxConns          int32
	MinConns          int32
	MaxConnLifetime   time.Duration
	MaxConnIdleTime   time.Duration
	HealthCheckPeriod time.Duration
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type JWTConfig struct {
	Secret    string
	ExpiresIn time.Duration
	Issuer    string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		App: AppConfig{
			Env:             getEnv("APP_ENV", "development"),
			Port:            getEnv("APP_PORT", getEnv("PORT", "8080")),
			ShutdownTimeout: getEnvAsDuration("SHUTDOWN_TIMEOUT", 30*time.Second),
		},
		Database: DatabaseConfig{
			URL:               getEnv("DATABASE_URL", ""),
			Host:              getEnv("DB_HOST", "localhost"),
			Port:              getEnvAsInt("DB_PORT", 5432),
			User:              getEnv("DB_USER", "postgres"),
			Password:          getEnv("DB_PASSWORD", "postgres"),
			Name:              getEnv("DB_NAME", "go_ecommerce_api"),
			SSLMode:           getEnv("DB_SSLMODE", "disable"),
			MaxConns:          int32(getEnvAsInt("DB_POOL_MAX_CONNS", 10)),
			MinConns:          int32(getEnvAsInt("DB_POOL_MIN_CONNS", 0)),
			MaxConnLifetime:   getEnvAsDuration("DB_POOL_MAX_CONN_LIFETIME", time.Hour),
			MaxConnIdleTime:   getEnvAsDuration("DB_POOL_MAX_CONN_IDLE_TIME", 30*time.Minute),
			HealthCheckPeriod: getEnvAsDuration("DB_POOL_HEALTH_CHECK_PERIOD", time.Minute),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnvAsInt("REDIS_PORT", 6379),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 2),
			DialTimeout:  getEnvAsDuration("REDIS_DIAL_TIMEOUT", 5*time.Second),
			ReadTimeout:  getEnvAsDuration("REDIS_READ_TIMEOUT", 3*time.Second),
			WriteTimeout: getEnvAsDuration("REDIS_WRITE_TIMEOUT", 3*time.Second),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "change_me_for_local_dev"),
			ExpiresIn: getEnvAsDuration("JWT_EXPIRES_IN", 24*time.Hour),
			Issuer:    getEnv("JWT_ISSUER", "go-ecommerce-api"),
		},
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.App.Port == "" {
		return fmt.Errorf("APP_PORT is required")
	}

	if cfg.App.ShutdownTimeout <= 0 {
		return fmt.Errorf("SHUTDOWN_TIMEOUT must be greater than 0")
	}

	if cfg.Database.URL == "" {
		if cfg.Database.Host == "" {
			return fmt.Errorf("DB_HOST is required")
		}

		if cfg.Database.Port <= 0 {
			return fmt.Errorf("DB_PORT must be greater than 0")
		}

		if cfg.Database.User == "" {
			return fmt.Errorf("DB_USER is required")
		}

		if cfg.Database.Name == "" {
			return fmt.Errorf("DB_NAME is required")
		}
	}

	if cfg.Database.MaxConns <= 0 {
		return fmt.Errorf("DB_POOL_MAX_CONNS must be greater than 0")
	}

	if cfg.Database.MinConns < 0 {
		return fmt.Errorf("DB_POOL_MIN_CONNS cannot be negative")
	}

	if cfg.Database.MinConns > cfg.Database.MaxConns {
		return fmt.Errorf("DB_POOL_MIN_CONNS cannot be greater than DB_POOL_MAX_CONNS")
	}

	if cfg.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}

	if cfg.JWT.Issuer == "" {
		return fmt.Errorf("JWT_ISSUER is required")
	}

	if cfg.App.Env == "production" && cfg.JWT.Secret == "change_me_for_local_dev" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	if cfg.Redis.Host == "" {
		return fmt.Errorf("REDIS_HOST is required")
	}

	if cfg.Redis.Port <= 0 {
		return fmt.Errorf("REDIS_PORT must be greater than 0")
	}

	if cfg.Redis.DB < 0 {
		return fmt.Errorf("REDIS_DB cannot be negative")
	}

	if cfg.Redis.PoolSize <= 0 {
		return fmt.Errorf("REDIS_POOL_SIZE must be greater than 0")
	}

	if cfg.Redis.MinIdleConns < 0 {
		return fmt.Errorf("REDIS_MIN_IDLE_CONNS cannot be negative")
	}

	return nil
}

func (db DatabaseConfig) DSN() string {
	if db.URL != "" {
		return db.URL
	}

	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(db.User, db.Password),
		Host:   fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:   "/" + db.Name,
	}

	query := dsn.Query()
	query.Set("sslmode", db.SSLMode)

	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func (r RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvAsInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return intValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return duration
}

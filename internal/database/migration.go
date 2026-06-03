package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func RunMigrations(cfg config.DatabaseConfig) error {
	m, err := migrate.New(
		"file://migrations",
		cfg.DSN(),
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	defer func() {
		sourceErr, databaseErr := m.Close()
		if sourceErr != nil {
			log.Printf("failed to close migration source: %v", sourceErr)
		}
		if databaseErr != nil {
			log.Printf("failed to close migration database: %v", databaseErr)
		}
	}()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database migrations already up to date")
			return nil
		}

		return fmt.Errorf("failed to run database migrations: %w", err)
	}

	log.Println("Database migrations applied successfully")

	return nil
}

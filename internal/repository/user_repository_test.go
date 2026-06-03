package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
	"github.com/kodokbakar/go-ecommerce-api/internal/repository"
)

func TestUserRepositoryCreateSuccess(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer mock.Close()

	now := time.Now()

	user := &models.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
		Role:         models.RoleCustomer,
		IsActive:     true,
	}

	rows := pgxmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow("8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21", now, now)

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive).
		WillReturnRows(rows)

	repo := repository.NewUserRepository(mock)

	if err := repo.Create(context.Background(), user); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != "8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21" {
		t.Fatalf("expected user ID to be set, got %s", user.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepositoryCreateDuplicateEmail(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer mock.Close()

	user := &models.User{
		Name:         "Test User",
		Email:        "test@example.com",
		PasswordHash: "hashed-password",
		Role:         models.RoleCustomer,
		IsActive:     true,
	}

	mock.ExpectQuery("INSERT INTO users").
		WithArgs(user.Name, user.Email, user.PasswordHash, user.Role, user.IsActive).
		WillReturnError(&pgconn.PgError{
			Code: "23505",
		})

	repo := repository.NewUserRepository(mock)

	err = repo.Create(context.Background(), user)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, models.ErrDuplicateEmail) {
		t.Fatalf("expected ErrDuplicateEmail, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepositoryGetByEmailSuccess(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer mock.Close()

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id",
		"name",
		"email",
		"password_hash",
		"role",
		"is_active",
		"created_at",
		"updated_at",
	}).AddRow(
		"8b5d1d9a-2b0c-4fd7-9c27-25e05e79ad21",
		"Test User",
		"test@example.com",
		"hashed-password",
		models.RoleCustomer,
		true,
		now,
		now,
	)

	mock.ExpectQuery("SELECT").
		WithArgs("test@example.com").
		WillReturnRows(rows)

	repo := repository.NewUserRepository(mock)

	user, err := repo.GetByEmail(context.Background(), "test@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.Email != "test@example.com" {
		t.Fatalf("expected email test@example.com, got %s", user.Email)
	}

	if user.Role != models.RoleCustomer {
		t.Fatalf("expected role customer, got %s", user.Role)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

func TestUserRepositoryGetByEmailNotFound(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer mock.Close()

	mock.ExpectQuery("SELECT").
		WithArgs("missing@example.com").
		WillReturnError(pgx.ErrNoRows)

	repo := repository.NewUserRepository(mock)

	_, err = repo.GetByEmail(context.Background(), "missing@example.com")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, models.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unmet expectations: %v", err)
	}
}

package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type PgxQuerier interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type UserRepository struct {
	db PgxQuerier
}

func NewUserRepository(db PgxQuerier) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (
			name,
			email,
			password_hash,
			role,
			is_active
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.IsActive,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if isUniqueViolation(err) {
			return models.ErrDuplicateEmail
		}

		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT
			id,
			name,
			email,
			password_hash,
			role,
			is_active,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	user := &models.User{}

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrUserNotFound
		}

		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return user, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}

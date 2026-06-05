package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *models.Category) error
	FindAll(ctx context.Context) ([]models.Category, error)
	FindAllPaginated(ctx context.Context, filter CategoryListFilter) ([]models.Category, int, error)
	FindByID(ctx context.Context, id string) (*models.Category, error)
	FindBySlug(ctx context.Context, slug string) (*models.Category, error)
	ExistsByName(ctx context.Context, name string, excludeID string) (bool, error)
	ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error)
	HasProducts(ctx context.Context, categoryID string) (bool, error)
	HasChildren(ctx context.Context, categoryID string) (bool, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id string) error
}

type CategoryListFilter struct {
	Limit  int
	Offset int
}

type categoryRepository struct {
	db PgxQuerier
}

func NewCategoryRepository(db PgxQuerier) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *models.Category) error {
	query := `
		INSERT INTO categories (parent_id, name, slug, description, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id::text,
			COALESCE(parent_id::text, ''),
			name,
			slug,
			COALESCE(description, ''),
			COALESCE(image_url, ''),
			created_at,
			updated_at
	`

	var parentID string

	err := r.db.QueryRow(
		ctx,
		query,
		nullableString(category.ParentID),
		category.Name,
		category.Slug,
		category.Description,
		category.ImageURL,
	).Scan(
		&category.ID,
		&parentID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ImageURL,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return err
	}

	category.ParentID = stringPtrOrNil(parentID)

	return nil
}

func (r *categoryRepository) FindAll(ctx context.Context) ([]models.Category, error) {
	query := `
		SELECT
			id::text,
			COALESCE(parent_id::text, ''),
			name,
			slug,
			COALESCE(description, ''),
			COALESCE(image_url, ''),
			created_at,
			updated_at
		FROM categories
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {
		category, err := scanCategoryRows(rows)
		if err != nil {
			return nil, err
		}

		categories = append(categories, *category)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *categoryRepository) FindAllPaginated(ctx context.Context, filter CategoryListFilter) ([]models.Category, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM categories
	`

	var total int

	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT
		id::text,
		COALESCE(parent_id::text, ''),
		name,
		slug,
		COALESCE(description, ''),
		COALESCE(image_url, ''),
		created_at,
		updated_at
	FROM categories
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
`

	rows, err := r.db.Query(ctx, query, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	categories := make([]models.Category, 0)

	for rows.Next() {
		category := models.Category{}
		var parentID string

		if err := rows.Scan(
			&category.ID,
			&parentID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.ImageURL,
			&category.CreatedAt,
			&category.UpdatedAt,
		); err != nil {
			return nil, 0, err
		}

		if parentID != "" {
			category.ParentID = &parentID
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

func (r *categoryRepository) FindByID(ctx context.Context, id string) (*models.Category, error) {
	query := `
		SELECT
			id::text,
			COALESCE(parent_id::text, ''),
			name,
			slug,
			COALESCE(description, ''),
			COALESCE(image_url, ''),
			created_at,
			updated_at
		FROM categories
		WHERE id = $1
	`

	category, err := r.scanCategoryRow(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrCategoryNotFound
		}

		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) FindBySlug(ctx context.Context, slug string) (*models.Category, error) {
	query := `
		SELECT
			id::text,
			COALESCE(parent_id::text, ''),
			name,
			slug,
			COALESCE(description, ''),
			COALESCE(image_url, ''),
			created_at,
			updated_at
		FROM categories
		WHERE slug = $1
	`

	category, err := r.scanCategoryRow(ctx, query, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrCategoryNotFound
		}

		return nil, err
	}

	return category, nil
}

func (r *categoryRepository) ExistsByName(ctx context.Context, name string, excludeID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM categories
			WHERE LOWER(name) = LOWER($1)
			AND ($2 = '' OR id::text <> $2)
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, name, excludeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM categories
			WHERE slug = $1
			AND ($2 = '' OR id::text <> $2)
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, slug, excludeID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) HasProducts(ctx context.Context, categoryID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM products
			WHERE category_id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) HasChildren(ctx context.Context, categoryID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM categories
			WHERE parent_id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, categoryID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *categoryRepository) Update(ctx context.Context, category *models.Category) error {
	query := `
		UPDATE categories
		SET
			parent_id = $2,
			name = $3,
			slug = $4,
			description = $5,
			image_url = $6
		WHERE id = $1
		RETURNING
			id::text,
			COALESCE(parent_id::text, ''),
			name,
			slug,
			COALESCE(description, ''),
			COALESCE(image_url, ''),
			created_at,
			updated_at
	`

	var parentID string

	err := r.db.QueryRow(
		ctx,
		query,
		category.ID,
		nullableString(category.ParentID),
		category.Name,
		category.Slug,
		category.Description,
		category.ImageURL,
	).Scan(
		&category.ID,
		&parentID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ImageURL,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrCategoryNotFound
		}

		return err
	}

	category.ParentID = stringPtrOrNil(parentID)

	return nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM categories
		WHERE id = $1
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrCategoryNotFound
	}

	return nil
}

func (r *categoryRepository) scanCategoryRow(ctx context.Context, query string, args ...any) (*models.Category, error) {
	var category models.Category
	var parentID string

	err := r.db.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&parentID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ImageURL,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	category.ParentID = stringPtrOrNil(parentID)

	return &category, nil
}

type categoryRowScanner interface {
	Scan(dest ...any) error
}

func scanCategoryRows(row categoryRowScanner) (*models.Category, error) {
	var category models.Category
	var parentID string

	err := row.Scan(
		&category.ID,
		&parentID,
		&category.Name,
		&category.Slug,
		&category.Description,
		&category.ImageURL,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	category.ParentID = stringPtrOrNil(parentID)

	return &category, nil
}

func nullableString(value *string) any {
	if value == nil {
		return nil
	}

	trimmed := *value
	if trimmed == "" {
		return nil
	}

	return trimmed
}

func stringPtrOrNil(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

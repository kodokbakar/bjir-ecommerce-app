package repository

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type ProductListFilter struct {
	CategoryID   string
	CategorySlug string
	Search       string
	Limit        int
	Offset       int
	SortBy       string
	SortOrder    string
}

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	FindAll(ctx context.Context, filter ProductListFilter) ([]models.Product, int, error)
	FindByID(ctx context.Context, id string) (*models.Product, error)
	FindBySlug(ctx context.Context, slug string) (*models.Product, error)
	ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error)
	Update(ctx context.Context, product *models.Product) error
	UpdateImageURL(ctx context.Context, id string, imageURL string) (*models.Product, error)
	Delete(ctx context.Context, id string) error
}

type productRepository struct {
	db PgxQuerier
}

func NewProductRepository(db PgxQuerier) ProductRepository {
	return &productRepository{db: db}
}

func buildProductOrderBy(sortBy string, sortOrder string) string {
	sortBy = strings.ToLower(strings.TrimSpace(sortBy))
	sortOrder = strings.ToUpper(strings.TrimSpace(sortOrder))

	if sortOrder != "ASC" && sortOrder != "DESC" {
		sortOrder = "DESC"
	}

	switch sortBy {
	case "price":
		return "p.price " + sortOrder + ", p.created_at DESC"
	case "name":
		return "LOWER(p.name) " + sortOrder + ", p.created_at DESC"
	case "created_at":
		return "p.created_at " + sortOrder
	default:
		return "p.created_at DESC"
	}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	query := `
		WITH inserted AS (
			INSERT INTO products (
				category_id,
				name,
				slug,
				description,
				price,
				stock,
				image_url,
				is_active
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, true)
			RETURNING *
		)
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM inserted p
		LEFT JOIN categories c ON c.id = p.category_id
	`

	err := r.scanProductRow(
		ctx,
		query,
		product,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) FindAll(ctx context.Context, filter ProductListFilter) ([]models.Product, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.is_active = true
		AND ($1 = '' OR p.category_id::text = $1)
		AND ($2 = '' OR p.name ILIKE '%' || $2 || '%')
		AND ($3 = '' OR c.slug = $3)
	`

	var total int

	if err := r.db.QueryRow(
		ctx,
		countQuery,
		filter.CategoryID,
		filter.Search,
		filter.CategorySlug,
	).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := buildProductOrderBy(filter.SortBy, filter.SortOrder)

	query := `
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.is_active = true
		AND ($1 = '' OR p.category_id::text = $1)
		AND ($2 = '' OR p.name ILIKE '%' || $2 || '%')
		AND ($3 = '' OR c.slug = $3)
		ORDER BY ` + orderBy + `
		LIMIT $4 OFFSET $5
	`

	rows, err := r.db.Query(
		ctx,
		query,
		filter.CategoryID,
		filter.Search,
		filter.CategorySlug,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)

	for rows.Next() {
		product := models.Product{}

		if err := scanProduct(&product, rows); err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*models.Product, error) {
	query := `
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.id = $1
		AND p.is_active = true
	`

	product := &models.Product{}

	err := r.scanProductRow(ctx, query, product, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	query := `
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM products p
		LEFT JOIN categories c ON c.id = p.category_id
		WHERE p.slug = $1
		AND p.is_active = true
	`

	product := &models.Product{}

	err := r.scanProductRow(ctx, query, product, slug)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepository) ExistsBySlug(ctx context.Context, slug string, excludeID string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM products
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

func (r *productRepository) Update(ctx context.Context, product *models.Product) error {
	query := `
		WITH updated AS (
			UPDATE products
			SET
				category_id = $2,
				name = $3,
				slug = $4,
				description = $5,
				price = $6,
				stock = $7,
				image_url = $8
			WHERE id = $1
			AND is_active = true
			RETURNING *
		)
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM updated p
		LEFT JOIN categories c ON c.id = p.category_id
	`

	err := r.scanProductRow(
		ctx,
		query,
		product,
		product.ID,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrProductNotFound
		}

		return err
	}

	return nil
}

func (r *productRepository) UpdateImageURL(ctx context.Context, id string, imageURL string) (*models.Product, error) {
	query := `
		WITH updated AS (
			UPDATE products
			SET image_url = $2
			WHERE id = $1
			AND is_active = true
			RETURNING *
		)
		SELECT
			p.id::text,
			COALESCE(p.category_id::text, ''),
			p.name,
			p.slug,
			COALESCE(p.description, ''),
			p.price::float8,
			p.stock,
			COALESCE(p.image_url, ''),
			p.is_active,
			p.created_at,
			p.updated_at,
			COALESCE(c.id::text, ''),
			COALESCE(c.parent_id::text, ''),
			COALESCE(c.name, ''),
			COALESCE(c.slug, ''),
			COALESCE(c.description, ''),
			COALESCE(c.image_url, ''),
			COALESCE(c.created_at, p.created_at),
			COALESCE(c.updated_at, p.updated_at)
		FROM updated p
		LEFT JOIN categories c ON c.id = p.category_id
	`

	product := &models.Product{}

	err := r.scanProductRow(ctx, query, product, id, imageURL)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrProductNotFound
		}

		return nil, err
	}

	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE products
		SET is_active = false
		WHERE id = $1
		AND is_active = true
	`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrProductNotFound
	}

	return nil
}

func (r *productRepository) scanProductRow(ctx context.Context, query string, product *models.Product, args ...any) error {
	return scanProduct(product, r.db.QueryRow(ctx, query, args...))
}

type productRowScanner interface {
	Scan(dest ...any) error
}

func scanProduct(product *models.Product, row productRowScanner) error {
	var categoryID string
	var categoryParentID string
	var categoryName string
	var categorySlug string
	var categoryDescription string
	var categoryImageURL string
	var categoryCreatedAt = product.CreatedAt
	var categoryUpdatedAt = product.UpdatedAt

	err := row.Scan(
		&product.ID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.ImageURL,
		&product.IsActive,
		&product.CreatedAt,
		&product.UpdatedAt,
		&categoryID,
		&categoryParentID,
		&categoryName,
		&categorySlug,
		&categoryDescription,
		&categoryImageURL,
		&categoryCreatedAt,
		&categoryUpdatedAt,
	)
	if err != nil {
		return err
	}

	if categoryID != "" {
		product.Category = &models.Category{
			ID:          categoryID,
			ParentID:    stringPtrOrNil(categoryParentID),
			Name:        categoryName,
			Slug:        categorySlug,
			Description: categoryDescription,
			ImageURL:    categoryImageURL,
			CreatedAt:   categoryCreatedAt,
			UpdatedAt:   categoryUpdatedAt,
		}
	}

	return nil
}

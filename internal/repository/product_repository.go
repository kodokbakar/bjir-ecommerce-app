package repository

import (
	"context"
	"errors"
	"fmt"
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

type ProductImageSortOrder struct {
	ID        string
	SortOrder int
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

	FindImagesByProductID(ctx context.Context, productID string) ([]models.ProductImage, error)
	CountImagesByProductID(ctx context.Context, productID string) (int, error)
	CreateProductImage(ctx context.Context, image *models.ProductImage) error
	DeleteProductImage(ctx context.Context, productID string, imageID string) error
	UpdateProductImageSortOrder(ctx context.Context, productID string, imageID string, sortOrder int) error
	BulkUpdateProductImageSortOrders(ctx context.Context, productID string, images []ProductImageSortOrder) error
	SetPrimaryProductImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error)
	SyncProductPrimaryImageURL(ctx context.Context, productID string) error
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

func (r *productRepository) FindImagesByProductID(ctx context.Context, productID string) ([]models.ProductImage, error) {
	query := `
		SELECT
			id::text,
			product_id::text,
			image_url,
			sort_order,
			is_primary,
			created_at,
			updated_at
		FROM product_images
		WHERE product_id = $1
		ORDER BY sort_order ASC, created_at ASC
	`

	rows, err := r.db.Query(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProductImageRows(rows)
}

func (r *productRepository) CountImagesByProductID(ctx context.Context, productID string) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM product_images
		WHERE product_id = $1
	`

	var count int
	if err := r.db.QueryRow(ctx, query, productID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (r *productRepository) CreateProductImage(ctx context.Context, image *models.ProductImage) error {
	query := `
		INSERT INTO product_images (
			product_id,
			image_url,
			sort_order,
			is_primary
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id::text,
			product_id::text,
			image_url,
			sort_order,
			is_primary,
			created_at,
			updated_at
	`

	err := scanProductImage(
		r.db.QueryRow(
			ctx,
			query,
			image.ProductID,
			image.ImageURL,
			image.SortOrder,
			image.IsPrimary,
		),
		image,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepository) DeleteProductImage(ctx context.Context, productID string, imageID string) error {
	query := `
		DELETE FROM product_images
		WHERE product_id = $1
		AND id = $2
		RETURNING id
	`

	var deletedID string
	if err := r.db.QueryRow(ctx, query, productID, imageID).Scan(&deletedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrProductImageNotFound
		}

		return err
	}

	return nil
}

func (r *productRepository) UpdateProductImageSortOrder(ctx context.Context, productID string, imageID string, sortOrder int) error {
	query := `
		UPDATE product_images
		SET sort_order = $3,
			updated_at = NOW()
		WHERE product_id = $1
		AND id = $2
		RETURNING id
	`

	var updatedID string
	if err := r.db.QueryRow(ctx, query, productID, imageID, sortOrder).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrProductImageNotFound
		}

		return err
	}

	return nil
}

func (r *productRepository) BulkUpdateProductImageSortOrders(ctx context.Context, productID string, images []ProductImageSortOrder) error {
	if len(images) == 0 {
		return nil
	}

	args := make([]any, 0, 1+(len(images)*2))
	args = append(args, productID)

	values := make([]string, 0, len(images))
	for _, image := range images {
		idPlaceholder := len(args) + 1
		sortOrderPlaceholder := len(args) + 2

		values = append(values, fmt.Sprintf("($%d::uuid, $%d::int)", idPlaceholder, sortOrderPlaceholder))
		args = append(args, image.ID, image.SortOrder)
	}

	query := `
		WITH input(id, sort_order) AS (
			VALUES ` + strings.Join(values, ", ") + `
		),
		updated AS (
			UPDATE product_images pi
			SET sort_order = input.sort_order,
				updated_at = NOW()
			FROM input
			WHERE pi.product_id = $1
			AND pi.id = input.id
			RETURNING pi.id
		)
		SELECT COUNT(*)
		FROM updated
	`

	var updatedCount int
	if err := r.db.QueryRow(ctx, query, args...).Scan(&updatedCount); err != nil {
		return err
	}

	if updatedCount != len(images) {
		return models.ErrProductImageNotFound
	}

	return nil
}

func (r *productRepository) SetPrimaryProductImage(ctx context.Context, productID string, imageID string) (*models.ProductImage, error) {
	query := `
		WITH target AS (
			SELECT
				id,
				product_id,
				image_url,
				sort_order,
				created_at,
				updated_at
			FROM product_images
			WHERE product_id = $1
			AND id = $2
		),
		unset_old AS (
			UPDATE product_images
			SET is_primary = FALSE,
				updated_at = NOW()
			WHERE product_id = $1
			AND EXISTS (SELECT 1 FROM target)
		),
		set_new AS (
			UPDATE product_images
			SET is_primary = TRUE,
				updated_at = NOW()
			WHERE product_id = $1
			AND id = $2
			RETURNING
				id::text,
				product_id::text,
				image_url,
				sort_order,
				is_primary,
				created_at,
				updated_at
		),
		sync_product AS (
			UPDATE products
			SET image_url = (SELECT image_url FROM set_new),
				updated_at = NOW()
			WHERE id = $1
		)
		SELECT
			id,
			product_id,
			image_url,
			sort_order,
			is_primary,
			created_at,
			updated_at
		FROM set_new
	`

	image := &models.ProductImage{}
	if err := scanProductImage(r.db.QueryRow(ctx, query, productID, imageID), image); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrProductImageNotFound
		}

		return nil, err
	}

	return image, nil
}

func (r *productRepository) SyncProductPrimaryImageURL(ctx context.Context, productID string) error {
	query := `
		WITH current_primary AS (
			SELECT id, image_url
			FROM product_images
			WHERE product_id = $1
			AND is_primary = TRUE
			ORDER BY sort_order ASC, created_at ASC
			LIMIT 1
		),
		first_image AS (
			SELECT id, image_url
			FROM product_images
			WHERE product_id = $1
			ORDER BY sort_order ASC, created_at ASC
			LIMIT 1
		),
		chosen AS (
			SELECT id, image_url FROM current_primary
			UNION ALL
			SELECT id, image_url FROM first_image
			WHERE NOT EXISTS (SELECT 1 FROM current_primary)
			LIMIT 1
		),
		ensure_primary AS (
			UPDATE product_images
			SET is_primary = TRUE,
				updated_at = NOW()
			WHERE id = (SELECT id FROM chosen)
			AND NOT EXISTS (SELECT 1 FROM current_primary)
		)
		UPDATE products
		SET image_url = COALESCE((SELECT image_url FROM chosen), ''),
			updated_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	var updatedID string
	if err := r.db.QueryRow(ctx, query, productID).Scan(&updatedID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.ErrProductNotFound
		}

		return err
	}

	return nil
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

func scanProductImage(row pgx.Row, image *models.ProductImage) error {
	return row.Scan(
		&image.ID,
		&image.ProductID,
		&image.ImageURL,
		&image.SortOrder,
		&image.IsPrimary,
		&image.CreatedAt,
		&image.UpdatedAt,
	)
}

func scanProductImageRows(rows pgx.Rows) ([]models.ProductImage, error) {
	images := make([]models.ProductImage, 0)

	for rows.Next() {
		image := models.ProductImage{}
		if err := rows.Scan(
			&image.ID,
			&image.ProductID,
			&image.ImageURL,
			&image.SortOrder,
			&image.IsPrimary,
			&image.CreatedAt,
			&image.UpdatedAt,
		); err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return images, nil
}

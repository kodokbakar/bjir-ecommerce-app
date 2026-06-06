package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type CartRepository interface {
	Create(ctx context.Context, item *models.CartItem) error
	FindAllByUserID(ctx context.Context, userID string) ([]models.CartItem, error)
	FindByID(ctx context.Context, id string, userID string) (*models.CartItem, error)
	FindByUserAndProduct(ctx context.Context, userID string, productID string) (*models.CartItem, error)
	UpdateQuantity(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error)
	Delete(ctx context.Context, id string, userID string) error
}

type cartRepository struct {
	db PgxQuerier
}

func NewCartRepository(db PgxQuerier) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Create(ctx context.Context, item *models.CartItem) error {
	query := `
		WITH inserted AS (
			INSERT INTO carts (user_id, product_id, quantity)
			VALUES ($1, $2, $3)
			RETURNING *
		)
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
			c.created_at,
			c.updated_at,
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
			p.updated_at
		FROM inserted c
		JOIN products p ON p.id = c.product_id
	`

	createdItem, err := r.scanCartItemRow(ctx, query, item.UserID, item.ProductID, item.Quantity)
	if err != nil {
		return err
	}

	*item = *createdItem

	return nil
}

func (r *cartRepository) FindAllByUserID(ctx context.Context, userID string) ([]models.CartItem, error) {
	query := `
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
			c.created_at,
			c.updated_at,
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
			p.updated_at
		FROM carts c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
		ORDER BY c.updated_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.CartItem, 0)

	for rows.Next() {
		item, err := scanCartItem(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *cartRepository) FindByID(ctx context.Context, id string, userID string) (*models.CartItem, error) {
	query := `
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
			c.created_at,
			c.updated_at,
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
			p.updated_at
		FROM carts c
		JOIN products p ON p.id = c.product_id
		WHERE c.id = $1
		AND c.user_id = $2
	`

	item, err := r.scanCartItemRow(ctx, query, id, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrCartItemNotFound
		}

		return nil, err
	}

	return item, nil
}

func (r *cartRepository) FindByUserAndProduct(ctx context.Context, userID string, productID string) (*models.CartItem, error) {
	query := `
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
			c.created_at,
			c.updated_at,
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
			p.updated_at
		FROM carts c
		JOIN products p ON p.id = c.product_id
		WHERE c.user_id = $1
		AND c.product_id = $2
	`

	item, err := r.scanCartItemRow(ctx, query, userID, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrCartItemNotFound
		}

		return nil, err
	}

	return item, nil
}

func (r *cartRepository) UpdateQuantity(ctx context.Context, id string, userID string, quantity int) (*models.CartItem, error) {
	query := `
		WITH updated AS (
			UPDATE carts
			SET quantity = $3
			WHERE id = $1
			AND user_id = $2
			RETURNING *
		)
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
			c.created_at,
			c.updated_at,
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
			p.updated_at
		FROM updated c
		JOIN products p ON p.id = c.product_id
	`

	item, err := r.scanCartItemRow(ctx, query, id, userID, quantity)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrCartItemNotFound
		}

		return nil, err
	}

	return item, nil
}

func (r *cartRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `
		DELETE FROM carts
		WHERE id = $1
		AND user_id = $2
	`

	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrCartItemNotFound
	}

	return nil
}

func (r *cartRepository) scanCartItemRow(ctx context.Context, query string, args ...any) (*models.CartItem, error) {
	return scanCartItem(r.db.QueryRow(ctx, query, args...))
}

type cartItemRowScanner interface {
	Scan(dest ...any) error
}

func scanCartItem(row cartItemRowScanner) (*models.CartItem, error) {
	var item models.CartItem
	var product models.Product

	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.ProductID,
		&item.Quantity,
		&item.CreatedAt,
		&item.UpdatedAt,
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
	)
	if err != nil {
		return nil, err
	}

	item.Product = &product
	item.Subtotal = product.Price * float64(item.Quantity)

	return &item, nil
}

package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type OrderRepository interface {
	Checkout(ctx context.Context, userID string) (*models.Order, error)
}

type orderRepository struct {
	db PgxTxBeginner
}

type PgxTxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewOrderRepository(db PgxTxBeginner) OrderRepository {
	return &orderRepository{db: db}
}

type checkoutCartItem struct {
	CartItemID string
	UserID     string
	ProductID  string
	Quantity   int
	Product    models.Product
}

func (r *orderRepository) Checkout(ctx context.Context, userID string) (*models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	cartItems, err := r.getCheckoutCartItems(ctx, tx, userID)
	if err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, models.ErrCartEmpty
	}

	totalAmount := calculateOrderTotal(cartItems)

	for _, cartItem := range cartItems {
		if err := r.decreaseProductStock(ctx, tx, cartItem.ProductID, cartItem.Quantity); err != nil {
			return nil, err
		}
	}

	order := &models.Order{
		UserID:          userID,
		OrderNumber:     generateOrderNumber(),
		Status:          models.OrderStatusPending,
		TotalAmount:     totalAmount,
		ShippingAddress: "",
		Notes:           "",
		Items:           make([]models.OrderItem, 0, len(cartItems)),
	}

	if err := r.createOrder(ctx, tx, order); err != nil {
		return nil, err
	}

	for _, cartItem := range cartItems {
		orderItem := models.OrderItem{
			OrderID:     order.ID,
			ProductID:   cartItem.ProductID,
			ProductName: cartItem.Product.Name,
			Quantity:    cartItem.Quantity,
			Price:       cartItem.Product.Price,
			Subtotal:    cartItem.Product.Price * float64(cartItem.Quantity),
		}

		if err := r.createOrderItem(ctx, tx, &orderItem); err != nil {
			return nil, err
		}

		order.Items = append(order.Items, orderItem)
	}

	if err := r.clearCart(ctx, tx, userID); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	committed = true

	return order, nil
}

func (r *orderRepository) getCheckoutCartItems(ctx context.Context, tx pgx.Tx, userID string) ([]checkoutCartItem, error) {
	query := `
		SELECT
			c.id::text,
			c.user_id::text,
			c.product_id::text,
			c.quantity,
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
		AND p.is_active = true
		ORDER BY c.created_at ASC
		FOR UPDATE OF c, p
	`

	rows, err := tx.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]checkoutCartItem, 0)

	for rows.Next() {
		var item checkoutCartItem

		err := rows.Scan(
			&item.CartItemID,
			&item.UserID,
			&item.ProductID,
			&item.Quantity,
			&item.Product.ID,
			&item.Product.CategoryID,
			&item.Product.Name,
			&item.Product.Slug,
			&item.Product.Description,
			&item.Product.Price,
			&item.Product.Stock,
			&item.Product.ImageURL,
			&item.Product.IsActive,
			&item.Product.CreatedAt,
			&item.Product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *orderRepository) createOrder(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `
		INSERT INTO orders (
			user_id,
			order_number,
			status,
			total_amount,
			shipping_address,
			notes
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id::text,
			user_id::text,
			order_number,
			status,
			total_amount::float8,
			COALESCE(shipping_address, ''),
			COALESCE(notes, ''),
			created_at,
			updated_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		order.UserID,
		order.OrderNumber,
		order.Status,
		order.TotalAmount,
		order.ShippingAddress,
		order.Notes,
	).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) decreaseProductStock(ctx context.Context, tx pgx.Tx, productID string, quantity int) error {
	query := `
		UPDATE products
		SET stock = stock - $1
		WHERE id = $2
		AND is_active = true
		AND stock >= $1
	`

	result, err := tx.Exec(ctx, query, quantity, productID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrInsufficientStock
	}

	return nil
}

func (r *orderRepository) createOrderItem(ctx context.Context, tx pgx.Tx, item *models.OrderItem) error {
	query := `
		INSERT INTO order_items (
			order_id,
			product_id,
			product_name,
			quantity,
			price,
			subtotal
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING
			id::text,
			order_id::text,
			product_id::text,
			product_name,
			quantity,
			price::float8,
			subtotal::float8,
			created_at
	`

	err := tx.QueryRow(
		ctx,
		query,
		item.OrderID,
		item.ProductID,
		item.ProductName,
		item.Quantity,
		item.Price,
		item.Subtotal,
	).Scan(
		&item.ID,
		&item.OrderID,
		&item.ProductID,
		&item.ProductName,
		&item.Quantity,
		&item.Price,
		&item.Subtotal,
		&item.CreatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepository) clearCart(ctx context.Context, tx pgx.Tx, userID string) error {
	query := `
		DELETE FROM carts
		WHERE user_id = $1
	`

	_, err := tx.Exec(ctx, query, userID)
	return err
}

func calculateOrderTotal(items []checkoutCartItem) float64 {
	var total float64

	for _, item := range items {
		total += item.Product.Price * float64(item.Quantity)
	}

	return total
}

func generateOrderNumber() string {
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return fmt.Sprintf("ORD-%d", time.Now().UnixNano())
	}

	return fmt.Sprintf("ORD-%s-%s", time.Now().Format("20060102150405"), hex.EncodeToString(randomBytes))
}

package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type OrderRepository interface {
	Checkout(ctx context.Context, userID string) (*models.Order, error)
	FindAll(ctx context.Context, filter OrderListFilter) ([]models.Order, int, error)
	FindAllByUserID(ctx context.Context, userID string, filter OrderListFilter) ([]models.Order, int, error)
	FindByIDAndUserID(ctx context.Context, orderID string, userID string) (*models.Order, error)
	FindByID(ctx context.Context, orderID string) (*models.Order, error)
	UpdateStatus(ctx context.Context, orderID string, currentStatus string, nextStatus string) (*models.Order, error)
}

type OrderListFilter struct {
	Limit  int
	Offset int
	Status string
	Search string
}

type orderRepository struct {
	db PgxOrderDB
}

type PgxOrderDB interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func NewOrderRepository(db PgxOrderDB) OrderRepository {
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

func (r *orderRepository) FindAll(ctx context.Context, filter OrderListFilter) ([]models.Order, int, error) {
	whereClause, args := buildAdminOrderWhereClause(filter)

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE %s
	`, whereClause)

	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	queryArgs := append([]any{}, args...)
	queryArgs = append(queryArgs, filter.Limit, filter.Offset)

	limitPlaceholder := fmt.Sprintf("$%d", len(args)+1)
	offsetPlaceholder := fmt.Sprintf("$%d", len(args)+2)

	query := fmt.Sprintf(`
		SELECT
			o.id::text,
			o.user_id::text,
			COALESCE(u.name, ''),
			COALESCE(u.email, ''),
			o.order_number,
			o.status,
			o.total_amount::float8,
			COALESCE(o.shipping_address, ''),
			COALESCE(o.notes, ''),
			o.created_at,
			o.updated_at
		FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE %s
		ORDER BY o.created_at DESC
		LIMIT %s OFFSET %s
	`, whereClause, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.Query(ctx, query, queryArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		var order models.Order

		if err := scanOrderWithCustomer(&order, rows); err != nil {
			return nil, 0, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func buildAdminOrderWhereClause(filter OrderListFilter) (string, []any) {
	clauses := []string{"1 = 1"}
	args := make([]any, 0, 2)

	status := strings.TrimSpace(filter.Status)
	if status != "" {
		args = append(args, status)
		clauses = append(clauses, fmt.Sprintf("o.status = $%d", len(args)))
	}

	search := strings.TrimSpace(filter.Search)
	if search != "" {
		args = append(args, "%"+search+"%")
		clauses = append(clauses, fmt.Sprintf("o.order_number ILIKE $%d", len(args)))
	}

	return strings.Join(clauses, " AND "), args
}

func (r *orderRepository) FindAllByUserID(ctx context.Context, userID string, filter OrderListFilter) ([]models.Order, int, error) {
	countQuery := `
		SELECT COUNT(*)
		FROM orders
		WHERE user_id = $1
	`

	var total int

	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT
			id::text,
			user_id::text,
			order_number,
			status,
			total_amount::float8,
			COALESCE(shipping_address, ''),
			COALESCE(notes, ''),
			created_at,
			updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, userID, filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	orders := make([]models.Order, 0)

	for rows.Next() {
		var order models.Order

		if err := scanOrder(&order, rows); err != nil {
			return nil, 0, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) FindByIDAndUserID(ctx context.Context, orderID string, userID string) (*models.Order, error) {
	query := `
		SELECT
			id::text,
			user_id::text,
			order_number,
			status,
			total_amount::float8,
			COALESCE(shipping_address, ''),
			COALESCE(notes, ''),
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
		AND user_id = $2
	`

	order := &models.Order{}

	if err := scanOrder(order, r.db.QueryRow(ctx, query, orderID, userID)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}

		return nil, err
	}

	items, err := r.findOrderItemsByOrderID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	order.Items = items

	return order, nil
}

func (r *orderRepository) FindByID(ctx context.Context, orderID string) (*models.Order, error) {
	query := `
		SELECT
			id::text,
			user_id::text,
			order_number,
			status,
			total_amount::float8,
			COALESCE(shipping_address, ''),
			COALESCE(notes, ''),
			created_at,
			updated_at
		FROM orders
		WHERE id = $1
	`

	order := &models.Order{}

	if err := scanOrder(order, r.db.QueryRow(ctx, query, orderID)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}

		return nil, err
	}

	return order, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID string, currentStatus string, nextStatus string) (*models.Order, error) {
	query := `
		UPDATE orders
		SET status = $3
		WHERE id = $1
		AND status = $2
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

	order := &models.Order{}

	if err := scanOrder(order, r.db.QueryRow(ctx, query, orderID, currentStatus, nextStatus)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrInvalidOrderStatusTransition
		}

		return nil, err
	}

	return order, nil
}

func (r *orderRepository) findOrderItemsByOrderID(ctx context.Context, orderID string) ([]models.OrderItem, error) {
	query := `
		SELECT
			id::text,
			order_id::text,
			product_id::text,
			product_name,
			quantity,
			price::float8,
			subtotal::float8,
			created_at
		FROM order_items
		WHERE order_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]models.OrderItem, 0)

	for rows.Next() {
		var item models.OrderItem

		if err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.ProductName,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
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

type orderRowScanner interface {
	Scan(dest ...any) error
}

func scanOrderWithCustomer(order *models.Order, row orderRowScanner) error {
	return row.Scan(
		&order.ID,
		&order.UserID,
		&order.UserName,
		&order.UserEmail,
		&order.OrderNumber,
		&order.Status,
		&order.TotalAmount,
		&order.ShippingAddress,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
}

func scanOrder(order *models.Order, row orderRowScanner) error {
	return row.Scan(
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

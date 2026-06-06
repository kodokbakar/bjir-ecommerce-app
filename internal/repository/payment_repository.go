package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/kodokbakar/go-ecommerce-api/internal/models"
)

type CreatePaymentInput struct {
	UserID        string
	OrderID       string
	Provider      string
	PaymentMethod string
	TransactionID string
	Status        string
}

type PaymentRepository interface {
	CreateForOrder(ctx context.Context, input CreatePaymentInput) (*models.Payment, error)
}

type paymentRepository struct {
	db PgxBeginner
}

type PgxBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

func NewPaymentRepository(db PgxBeginner) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) CreateForOrder(ctx context.Context, input CreatePaymentInput) (*models.Payment, error) {
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

	order, err := lockPayableOrder(ctx, tx, input.OrderID, input.UserID)
	if err != nil {
		return nil, err
	}

	if order.Status != models.OrderStatusPending {
		return nil, models.ErrOrderNotPayable
	}

	if err := markOrderPaid(ctx, tx, order.ID); err != nil {
		return nil, err
	}

	payment, err := insertPayment(ctx, tx, input, order.TotalAmount)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	committed = true

	return payment, nil
}

func lockPayableOrder(ctx context.Context, tx pgx.Tx, orderID string, userID string) (*models.Order, error) {
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
		FOR UPDATE
	`

	order := &models.Order{}

	err := scanOrder(order, tx.QueryRow(ctx, query, orderID, userID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, models.ErrOrderNotFound
		}

		return nil, err
	}

	return order, nil
}

func markOrderPaid(ctx context.Context, tx pgx.Tx, orderID string) error {
	query := `
		UPDATE orders
		SET status = $2
		WHERE id = $1
		AND status = $3
	`

	result, err := tx.Exec(
		ctx,
		query,
		orderID,
		models.OrderStatusPaid,
		models.OrderStatusPending,
	)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return models.ErrOrderNotPayable
	}

	return nil
}

func insertPayment(ctx context.Context, tx pgx.Tx, input CreatePaymentInput, amount float64) (*models.Payment, error) {
	query := `
		INSERT INTO payments (
			order_id,
			provider,
			payment_method,
			transaction_id,
			amount,
			status,
			paid_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		RETURNING
			id::text,
			order_id::text,
			provider,
			payment_method,
			COALESCE(transaction_id, ''),
			amount::float8,
			status,
			paid_at,
			created_at,
			updated_at
	`

	payment := &models.Payment{}

	err := scanPayment(
		payment,
		tx.QueryRow(
			ctx,
			query,
			input.OrderID,
			input.Provider,
			input.PaymentMethod,
			input.TransactionID,
			amount,
			input.Status,
		),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, models.ErrPaymentAlreadyExists
		}

		return nil, err
	}

	return payment, nil
}

type paymentRowScanner interface {
	Scan(dest ...any) error
}

func scanPayment(payment *models.Payment, row paymentRowScanner) error {
	var paidAt time.Time

	err := row.Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Provider,
		&payment.PaymentMethod,
		&payment.TransactionID,
		&payment.Amount,
		&payment.Status,
		&paidAt,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	payment.PaidAt = &paidAt

	return nil
}
